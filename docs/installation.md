# Install and update

Purpose: build, install, update, and uninstall the supported Pathframe binary.

Pathframe 1.x supports Linux amd64. Build a versioned binary with Go 1.26, then install it atomically:

```sh
go build -trimpath -ldflags "-X github.com/0xkhdr/pathframe/internal/app.Version=v1.0.0" -o pathframe ./cmd/pathframe
./scripts/install.sh install ./pathframe "$HOME/.local/bin/pathframe"
pathframe --version
```

Update by passing a newly built, executable binary to the same command. The installer stages the complete binary beside the destination and performs one atomic rename, so an error or interruption before that boundary leaves the prior binary installed. It never reads or changes `.pathframe/` project data or host configuration.

Uninstall only the binary:

```sh
./scripts/install.sh uninstall "$HOME/.local/bin/pathframe"
```

The default destination is `/usr/local/bin/pathframe`. The script refuses symlink destinations. Back up project data separately; uninstall does not delete it.
