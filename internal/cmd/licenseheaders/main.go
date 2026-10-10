// Copyright © 2026 Pocket ID
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// baseFile pins the last upstream commit before the fork's changes. Update it
	// whenever the fork is rebased onto a newer upstream revision.
	baseFile          = ".license-header-base"
	modifiedNotice    = "// Modified by contributors to the Pocket ID fork."
	spdxNotice        = "// SPDX-License-Identifier: Apache-2.0"
	pocketIDCopyright = "// Copyright © %d Pocket ID"
)

var pocketIDCopyrightPattern = regexp.MustCompile(`^// Copyright © ([0-9]{4}) Pocket ID$`)
var oryCopyrightPattern = regexp.MustCompile(`^// Copyright © [0-9]{4}(?:-[0-9]{4})? Ory Corp$`)

type fileKind int

const (
	inheritedUnmodified fileKind = iota
	inheritedModified
	newFile
)

func main() {
	fix := flag.Bool("fix", false, "rewrite incorrect license headers")
	check := flag.Bool("check", false, "check license headers without changing files")
	flag.Parse()

	if *fix && *check {
		fatal(errors.New("--fix and --check are mutually exclusive"))
	}
	if !*fix && !*check {
		*check = true
	}

	root, err := gitOutput("rev-parse", "--show-toplevel")
	if err != nil {
		fatal(fmt.Errorf("find repository root: %w", err))
	}
	repositoryRoot := strings.TrimSpace(string(root))

	base, err := os.ReadFile(filepath.Join(repositoryRoot, baseFile))
	if err != nil {
		fatal(fmt.Errorf("read %s: %w", baseFile, err))
	}
	baseCommit := strings.TrimSpace(string(base))
	if baseCommit == "" {
		fatal(fmt.Errorf("%s is empty", baseFile))
	}
	if _, err := gitOutputAt(repositoryRoot, "cat-file", "-e", baseCommit+"^{commit}"); err != nil {
		fatal(fmt.Errorf("upstream base commit %s is unavailable; fetch full Git history: %w", baseCommit, err))
	}

	baseFiles, err := listFilesAt(repositoryRoot, baseCommit)
	if err != nil {
		fatal(err)
	}
	currentFiles, err := listCurrentFiles(repositoryRoot)
	if err != nil {
		fatal(err)
	}

	var incorrect []string
	updated := 0
	for _, path := range currentFiles {
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		current, err := os.ReadFile(absolutePath)
		if err != nil {
			fatal(fmt.Errorf("read %s: %w", path, err))
		}

		kind := newFile
		upstreamCopyright := ""
		if _, ok := baseFiles[path]; ok {
			original, err := gitOutputAt(repositoryRoot, "show", baseCommit+":"+path)
			if err != nil {
				fatal(fmt.Errorf("read %s from upstream base: %w", path, err))
			}
			upstreamCopyright, err = findOryCopyright(original)
			if err != nil {
				fatal(fmt.Errorf("read upstream copyright for %s: %w", path, err))
			}
			kind = inheritedUnmodified
			if !bytes.Equal(stripManagedHeader(current), stripManagedHeader(original)) {
				kind = inheritedModified
			}
		} else if oryCopyright, err := findOryCopyright(current); err == nil {
			// Files ported from upstream after the base commit keep Ory's
			// copyright. The fork adapts them (for example their import paths),
			// so they are marked as modified.
			kind = inheritedModified
			upstreamCopyright = oryCopyright
		}

		year := newFileYear(repositoryRoot, path, current)
		expected := withHeader(current, desiredHeader(kind, year, upstreamCopyright))
		if bytes.Equal(current, expected) {
			continue
		}

		incorrect = append(incorrect, path)
		if *fix {
			info, err := os.Stat(absolutePath)
			if err != nil {
				fatal(fmt.Errorf("stat %s: %w", path, err))
			}
			if err := os.WriteFile(absolutePath, expected, info.Mode().Perm()); err != nil {
				fatal(fmt.Errorf("write %s: %w", path, err))
			}
			updated++
		}
	}

	if *fix {
		fmt.Printf("updated %d license header(s)\n", updated)
		return
	}
	if len(incorrect) > 0 {
		for _, path := range incorrect {
			fmt.Fprintf(os.Stderr, "incorrect license header: %s\n", path)
		}
		os.Exit(1)
	}
	fmt.Printf("checked %d license header(s)\n", len(currentFiles))
}

func listCurrentFiles(root string) ([]string, error) {
	output, err := gitOutputAt(root, "ls-files", "--cached", "--others", "--exclude-standard", "-z", "--", "*.go")
	if err != nil {
		return nil, fmt.Errorf("list current Go files: %w", err)
	}
	return splitNUL(output), nil
}

func listFilesAt(root, revision string) (map[string]struct{}, error) {
	output, err := gitOutputAt(root, "ls-tree", "-r", "--name-only", "-z", revision)
	if err != nil {
		return nil, fmt.Errorf("list Go files at upstream base: %w", err)
	}
	files := make(map[string]struct{})
	for _, path := range splitNUL(output) {
		if strings.HasSuffix(path, ".go") {
			files[path] = struct{}{}
		}
	}
	return files, nil
}

func splitNUL(output []byte) []string {
	parts := bytes.Split(output, []byte{0})
	files := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) > 0 {
			files = append(files, string(part))
		}
	}
	sort.Strings(files)
	return files
}

func desiredHeader(kind fileKind, year int, upstreamCopyright string) []string {
	switch kind {
	case inheritedUnmodified:
		return []string{upstreamCopyright, spdxNotice}
	case inheritedModified:
		return []string{upstreamCopyright, modifiedNotice, spdxNotice}
	case newFile:
		return []string{fmt.Sprintf(pocketIDCopyright, year), spdxNotice}
	default:
		panic("unknown file kind")
	}
}

func withHeader(content []byte, header []string) []byte {
	body := stripManagedHeader(content)
	result := []byte(strings.Join(header, "\n") + "\n\n")
	return append(result, body...)
}

func stripManagedHeader(content []byte) []byte {
	lines := bytes.Split(content, []byte("\n"))
	if len(lines) == 0 {
		return content
	}

	consumed := 0
	for consumed < len(lines) && isManagedHeaderLine(string(lines[consumed])) {
		consumed++
	}
	if consumed == 0 {
		return content
	}
	if consumed < len(lines) && len(lines[consumed]) == 0 {
		consumed++
	}
	return bytes.Join(lines[consumed:], []byte("\n"))
}

func isManagedHeaderLine(line string) bool {
	if line == modifiedNotice || line == spdxNotice {
		return true
	}
	if oryCopyrightPattern.MatchString(line) {
		return true
	}
	return pocketIDCopyrightPattern.MatchString(line)
}

func findOryCopyright(content []byte) (string, error) {
	for _, line := range bytes.Split(content, []byte("\n")) {
		text := string(line)
		if oryCopyrightPattern.MatchString(text) {
			return text, nil
		}
		if !isManagedHeaderLine(text) {
			break
		}
	}
	return "", errors.New("missing Ory copyright header")
}

func newFileYear(root, path string, content []byte) int {
	for _, line := range bytes.Split(content, []byte("\n")) {
		match := pocketIDCopyrightPattern.FindSubmatch(line)
		if len(match) == 2 {
			year, err := strconv.Atoi(string(match[1]))
			if err == nil {
				return year
			}
		}
		if !isManagedHeaderLine(string(line)) {
			break
		}
	}

	output, err := gitOutputAt(root, "log", "--diff-filter=A", "--format=%ad", "--date=format:%Y", "--", path)
	if err == nil {
		years := strings.Fields(string(output))
		if len(years) > 0 {
			year, conversionErr := strconv.Atoi(years[len(years)-1])
			if conversionErr == nil {
				return year
			}
		}
	}

	return time.Now().Year()
}

func gitOutput(arguments ...string) ([]byte, error) {
	return gitOutputAt("", arguments...)
}

func gitOutputAt(directory string, arguments ...string) ([]byte, error) {
	command := exec.Command("git", arguments...)
	command.Dir = directory
	output, err := command.Output()
	if err == nil {
		return output, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return nil, fmt.Errorf("git %s: %s", strings.Join(arguments, " "), strings.TrimSpace(string(exitError.Stderr)))
	}
	return nil, fmt.Errorf("git %s: %w", strings.Join(arguments, " "), err)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
