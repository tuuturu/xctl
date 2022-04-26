# Install XCTL

## Prerequisites

* Go 1.17

## Install

```shell
# First clone the repository
git clone git@github.com:tuuturu/xctl.git && cd xctl

# Build xctl
make build

# Install xctl
make install
```

Running `make install` will install `xctl` into your `~/.local/bin` folder.

⚠ Remember to add `~/.local/bin` to your path by adding `PATH=$PATH:~/.local/bin` to your `~/.bashrc` equivalent.

To modify where the script installs `xctl`, use `INSTALL_DIR=<new directory>`. For example:

```shell
INSTALL_DIR=~/.local/binaries make install
```

That should be it. `xctl` should now be available when running `xctl --help`.