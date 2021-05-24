# `deploy-keychain`

Command-line tool to permit the easy usage of multiple GitHub deploy keys simultaneously.

## Installation

### macOS

```shell
brew install nint8835/formulae/deploy-keychain
```

### Linux (with `apt`)

```shell
echo "deb [trusted=yes] https://packages.bootleg.technology/apt/ /" | sudo tee -a /etc/apt/sources.list.d/bootleg.technology.list
sudo apt-get update
sudo apt-get install deploy-keychain
```

### Linux (with `yum`)

```shell
echo "[bootleg.technology]\nname=bootleg.technology\nbaseurl=https://packages.bootleg.technology/yum/\nenabled=1\ngpgcheck=0" | sudo tee -a /etc/yum.repos.d/bootleg.technology.repo
sudo yum update
sudo yum install deploy-keychain
```

### Other OSes (and package managers)

If your desired target operating system or package manager aren't listed here, binaries and packages of various formats for most operating systems and CPU architectures are available via the GitHub releases in this repository.

## Usage

To use this tool, set the `GIT_SSH_COMMAND` environment variable to `deploy-keychain`, and then use Git as normal. For example:

```shell
# Use for all Git commands in this shell
export GIT_SSH_COMMAND=deploy-keychain

git fetch
git pull

# Use for a single Git command
GIT_SSH_COMMAND=deploy-keychain git clone ...
```

## Configuration

This tool looks for a file called `deploy-keychain.yml` in the either the folder `~/.deploy-keychain/`, or in the current directory.

This file follows the following format (all keys are optional)

| Key               | Description                                                                                                                                                  | Default                            |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------- |
| `key_path`        | Path to look for keys in by default.                                                                                                                         | `~/.ssh/deploy-keys/`              |
| `key_name_format` | Go text template string declaring the format of key file names to look for by default.                                                                       | `{{.account}}-{{.repository}}.pem` |
| `fallback_key`    | Path to a key to use if no other key matches for the current repo. Leave blank to error out instead.                                                         |                                    |
| `keys`            | A map of GitHub repository names (in the format `nint8835/deploy-keychain`, for example) to a path to a key on disk that should be used for that repository. |
| `ssh_command`     | Command used to connect over SSH.                                                                                                                            | `ssh`                              |

## Key Source Order
This tool will look for deploy keys matching the current repository in the following order:

1. `keys` attribute
2. Key in `key_path` matching `key_name_format`
3. `fallback_key` (if provided)

## Acknowledgements

The idea for this tool (and some of the trickier Git / SSH implementation details) is taken from [this Gist](https://gist.github.com/vhermecz/4e2ae9468f2ff7532bf3f8155ac95c74), for which the license can be found [here](https://gist.github.com/vhermecz/67b6e3491a8de566bf6d8577b9d431f1).
