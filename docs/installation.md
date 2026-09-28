# Install and update

Purpose: build, install, update, and uninstall the supported Pathframe binary.

Pathframe 1.x supports Linux amd64. Install or update the latest release:

```sh
curl --proto '=https' --tlsv1.2 -fsSL https://raw.githubusercontent.com/0xkhdr/pathframe/main/scripts/install.sh | sh
```

This installs to `$HOME/.local/bin/pathframe`; ensure that directory is on `PATH`. To install a specific release:

```sh
curl --proto '=https' --tlsv1.2 -fsSL https://raw.githubusercontent.com/0xkhdr/pathframe/main/scripts/install.sh | PATHFRAME_VERSION=v1.0.0 sh
```

The installer downloads the release archive and its SHA-256 checksum over HTTPS, verifies it, and atomically replaces only the executable. An error or interruption before replacement leaves the prior binary installed. It never reads or changes `.pathframe/` project data or host configuration.

Uninstall only the binary:

```sh
curl --proto '=https' --tlsv1.2 -fsSL https://raw.githubusercontent.com/0xkhdr/pathframe/main/scripts/install.sh | sh -s -- uninstall
```

You can inspect the script before running it instead of piping it directly to `sh`.

To build and install from a checkout instead:

```sh
go build -trimpath -ldflags "-X github.com/0xkhdr/pathframe/internal/app.Version=v1.0.0" -o pathframe ./cmd/pathframe
./scripts/install.sh install ./pathframe "$HOME/.local/bin/pathframe"
pathframe --version
```

Set `PATHFRAME_INSTALL_DIR=/usr/local/bin` for a system-wide install; elevated privileges may be required. The script refuses symlink destinations. Back up project data separately; uninstall does not delete it.
