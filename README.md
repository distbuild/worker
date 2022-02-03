# worker

[![Build Status](https://github.com/distbuild/worker/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/distbuild/worker/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/distbuild/worker/branch/main/graph/badge.svg?token=U2RWRDSJGK)](https://codecov.io/gh/distbuild/worker)
[![Go Report Card](https://goreportcard.com/badge/github.com/distbuild/worker)](https://goreportcard.com/report/github.com/distbuild/worker)
[![License](https://img.shields.io/github/license/distbuild/worker.svg)](https://github.com/distbuild/worker/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/distbuild/worker.svg)](https://github.com/distbuild/worker/releases/latest)
[![Gitter chat](https://badges.gitter.im/craftslab/distbuild.png)](https://gitter.im/craftslab/distbuild)



## Introduction

*worker* is the worker of [distbuild](https://github.com/distbuild) written in Go.



## Prerequisites

- Go >= 1.17.0



## Run

```bash
git clone https://github.com/distbuild/worker.git

cd worker
version=latest make build
./bin/worker --config-file="$PWD/config/config.yml"
```



## Docker

```bash
git clone https://github.com/distbuild/worker.git

cd worker
version=latest make docker
docker run -v "$PWD"/config:/tmp ghcr.io/distbuild/worker:latest --config-file="/tmp/config.yml"
```



## Usage

```
usage: worker --config-file=CONFIG-FILE [<flags>]

distbuild worker

Flags:
  --help                     Show context-sensitive help (also try --help-long
                             and --help-man).
  --version                  Show application version.
  --config-file=CONFIG-FILE  Config file (.yml)
```



## Settings

*worker* parameters can be set in the directory [config](https://github.com/distbuild/worker/blob/main/config).

An example of configuration in [config.yml](https://github.com/distbuild/worker/blob/main/config/config.yml):

```yaml
apiVersion: v1
kind: worker
metadata:
  name: worker
spec:
```



## License

Project License can be found [here](LICENSE).



## Reference
