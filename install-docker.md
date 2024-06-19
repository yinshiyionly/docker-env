# install docker

1. download the script

```shell
$ curl -fsSL https://get.docker.com -o install-docker.sh
```

2. verify the script's content

```shell
$ cat install-docker.sh
```

3. run the script with --dry-run to verify the steps it executes

```shell
$ sh install-docker.sh --dry-run
```

4. run the script either as root, or using sudo to perform the installation

```shell
$ sudo sh install-docker.sh
```

# Command-line options

1. `--version <VERSION>`, Use the `--version` option to install a specific version, for example:

```shell
$ sudo sh install-docker.sh --version 23.0
```

2. `--channel` <stable|test>, Use the --channel option to install from an alternative installation channel

```shell
$ sudo sh install-docker.sh --channel stable
```
