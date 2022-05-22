# Install XCTL

## Prerequisites

* [Go 1.17](https://go.dev/)
* [Linode account](https://www.linode.com/)

## Installation

```shell
# First clone the repository
git clone git@github.com:tuuturu/xctl.git && cd xctl

# Build xctl
make build

# Install xctl
make install
```

This will install `xctl` into your `~/.local/bin` folder.

!!! note
    Remember to add `~/.local/bin` to your path by adding `export PATH=$PATH:~/.local/bin` to your `~/.bashrc` equivalent.

To modify where the script installs `xctl`, use `INSTALL_DIR=<new directory>`. For example:

```shell
INSTALL_DIR=~/.local/binaries make install
```

That should be it. `xctl` should now be available. Test by running a command. For example:

```shell
xctl --version
```
