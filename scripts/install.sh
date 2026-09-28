#!/bin/sh
set -eu

repository=https://github.com/0xkhdr/pathframe
asset=pathframe_linux_amd64.tar.gz
target=${PATHFRAME_INSTALL_DIR:-"$HOME/.local/bin"}/pathframe
download_tmp=
install_tmp=

cleanup() {
	[ -z "$download_tmp" ] || rm -rf "$download_tmp"
	[ -z "$install_tmp" ] || rm -f "$install_tmp"
}
trap cleanup EXIT HUP INT TERM

usage() {
	echo "usage: install.sh [install [SOURCE [TARGET]] | uninstall [TARGET]]" >&2
	exit 2
}

install_file() {
	source=$1
	target=$2
	[ -f "$source" ] && [ -x "$source" ] || { echo "install.sh: source must be an executable regular file" >&2; exit 1; }
	[ ! -L "$target" ] || { echo "install.sh: refusing symlink target" >&2; exit 1; }
	dir=$(dirname "$target")
	mkdir -p "$dir"
	install_tmp="$dir/.pathframe-install.$$"
	cp "$source" "$install_tmp"
	chmod 0755 "$install_tmp"
	mv -f "$install_tmp" "$target"
	install_tmp=
}

download() {
	[ "$(uname -s)" = Linux ] && [ "$(uname -m)" = x86_64 ] || {
		echo "install.sh: only Linux amd64 is supported" >&2
		exit 1
	}
	for command in curl sha256sum tar mktemp; do
		command -v "$command" >/dev/null 2>&1 || { echo "install.sh: $command is required" >&2; exit 1; }
	done
	version=${PATHFRAME_VERSION:-latest}
	case "$version" in
	latest) base=$repository/releases/latest/download ;;
	v[0-9]*) base=$repository/releases/download/$version ;;
	*) echo "install.sh: PATHFRAME_VERSION must be 'latest' or start with 'v'" >&2; exit 1 ;;
	esac
	download_tmp=$(mktemp -d)
	curl -fsSL --proto '=https' --tlsv1.2 -o "$download_tmp/$asset" "$base/$asset"
	curl -fsSL --proto '=https' --tlsv1.2 -o "$download_tmp/$asset.sha256" "$base/$asset.sha256"
	(cd "$download_tmp" && sha256sum -c "$asset.sha256")
	tar -xzf "$download_tmp/$asset" -C "$download_tmp" pathframe
	install_file "$download_tmp/pathframe" "$target"
	echo "installed pathframe to $target"
}

action=${1-}
case "$action" in
"") download ;;
install)
	[ "$#" -le 3 ] || usage
	if [ "$#" -eq 1 ]; then download; else install_file "$2" "${3:-$target}"; fi
	;;
uninstall)
	[ "$#" -le 2 ] || usage
	target=${2:-$target}
	[ ! -L "$target" ] || { echo "install.sh: refusing symlink target" >&2; exit 1; }
	rm -f "$target"
	echo "uninstalled pathframe from $target"
	;;
*) usage ;;
esac
