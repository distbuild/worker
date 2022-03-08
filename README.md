# worker

[![Build Status](https://github.com/distbuild/worker/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/distbuild/worker/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/distbuild/worker/branch/main/graph/badge.svg?token=FM4NOMPT7Q)](https://codecov.io/gh/distbuild/worker)
[![License](https://img.shields.io/github/license/distbuild/worker.svg)](https://github.com/distbuild/worker/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/distbuild/worker.svg)](https://github.com/distbuild/worker/tags)
[![Gitter chat](https://badges.gitter.im/craftslab/distbuild.png)](https://gitter.im/craftslab/distbuild)



## Introduction

*worker* is the build worker of [distbuild](https://github.com/distbuild) written in Rust.



## Prerequisites

- Rust >= 1.59.0



## Run

```bash
git clone https://github.com/distbuild/worker.git

cd worker
make build
./target/release/worker --config-file="$PWD/src/config/config.yml"
```



## Docker

```bash
git clone https://github.com/distbuild/worker.git

cd worker
make docker
docker run -v "$PWD"/src/config:/tmp ghcr.io/distbuild/worker:latest --config-file="/tmp/config.yml"
```



## Usage

```
USAGE:
    worker --config-file <NAME>

OPTIONS:
    -c, --config-file <NAME>    Config file (.yml)
    -h, --help                  Print help information
    -V, --version               Print version information
```



## Settings

*worker* parameters can be set in the directory [config](https://github.com/distbuild/worker/blob/main/src/config).

An example of configuration in [config.yml](https://github.com/distbuild/worker/blob/main/src/config/config.yml):

```yaml
apiVersion: v1
kind: worker
metadata:
  name: worker
spec:
  foo: foo
```



## License

Project License can be found [here](LICENSE).



## Reference
