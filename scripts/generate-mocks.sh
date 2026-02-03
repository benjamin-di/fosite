#!/bin/sh
# Regenerates the gomock mocks in internal/ from the go:generate directives in
# generate.go. mockgen writes files without a license header, so the Ory
# copyright that upstream ships on each mock is kept and the license header
# tool then normalizes it.
set -eu

cd "$(dirname "$0")/.."

headers=$(mktemp -d)
trap 'rm -rf "$headers"' EXIT

for file in internal/*.go; do
	if head -1 "$file" | grep -q 'Ory Corp$'; then
		head -1 "$file" >"$headers/$(basename "$file")"
	fi
done

go generate .

for file in internal/*.go; do
	header="$headers/$(basename "$file")"
	if [ -f "$header" ] && ! head -1 "$file" | grep -q 'Copyright'; then
		{
			cat "$header"
			echo '// SPDX-License-Identifier: Apache-2.0'
			echo
			cat "$file"
		} >"$file.tmp"
		mv "$file.tmp" "$file"
	fi
done

go run ./internal/cmd/licenseheaders --fix
