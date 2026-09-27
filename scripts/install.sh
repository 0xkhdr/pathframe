#!/bin/sh
set -eu

usage() {
	echo "usage: install.sh install SOURCE [TARGET] | uninstall [TARGET]" >&2
	exit 2
}

action=${1-}
case "$action" in
install)
	[ "$#" -ge 2 ] && [ "$#" -le 3 ] || usage
	source=$2
	target=${3:-/usr/local/bin/pathframe}
	[ -f "$source" ] && [ -x "$source" ] || { echo "install.sh: source must be an executable regular file" >&2; exit 1; }
	[ ! -L "$target" ] || { echo "install.sh: refusing symlink target" >&2; exit 1; }
	dir=$(dirname "$target")
	mkdir -p "$dir"
	tmp="$dir/.pathframe-install.$$"
	trap 'rm -f "$tmp"' EXIT HUP INT TERM
	cp "$source" "$tmp"
	chmod 0755 "$tmp"
	mv -f "$tmp" "$target"
	trap - EXIT HUP INT TERM
	;;
uninstall)
	[ "$#" -le 2 ] || usage
	target=${2:-/usr/local/bin/pathframe}
	[ ! -L "$target" ] || { echo "install.sh: refusing symlink target" >&2; exit 1; }
	rm -f "$target"
	;;
*) usage ;;
esac
