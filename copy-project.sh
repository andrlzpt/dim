#!/usr/bin/env bash
set -euo pipefail

usage() {
    printf 'Usage: %s TARGET_DIRECTORY\n' "${0##*/}"
}

if [[ $# -eq 1 && ( $1 == --help || $1 == -h ) ]]; then
    usage
    exit 0
fi
if [[ $# -ne 1 || -z $1 ]]; then
    usage >&2
    exit 1
fi

# Resolve physical paths so symlinks cannot point the copy back into itself.
source_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
# Relative destinations live beside this project, independent of the caller's cwd.
case "$1" in
    /*) target_dir=$(realpath -m -- "$1") ;;
    *) target_dir=$(realpath -m -- "$source_dir/../$1") ;;
esac

case "$target_dir/" in
    "$source_dir/"*)
        printf 'No copy inside own cave. Pick directory outside %s\n' "$source_dir" >&2
        exit 1
        ;;
esac
case "$source_dir/" in
    "$target_dir/"*)
        printf 'Target is parent cave. Pick new or empty directory instead.\n' >&2
        exit 1
        ;;
esac

if [[ -e $target_dir && ! -d $target_dir ]]; then
    printf 'Target is not directory: %s\n' "$target_dir" >&2
    exit 1
fi
if [[ -d $target_dir ]]; then
    shopt -s nullglob dotglob
    entries=("$target_dir"/*)
    if (( ${#entries[@]} > 0 )); then
        printf 'Target cave has stuff. Use empty directory. No squash stuff.\n' >&2
        exit 1
    fi
fi

mkdir -p -- "$target_dir"
cp -a -- "$source_dir/." "$target_dir/"
printf 'Copied cave to %s\n' "$target_dir"

if command -v code >/dev/null 2>&1; then
    if ! code --new-window "$target_dir"; then
        printf 'Copy ready, but VS Code could not open. Open directory yourself: %s\n' "$target_dir" >&2
    fi
else
    printf 'Copy ready. VS Code command missing. Add code to PATH to open automatically.\n' >&2
fi
