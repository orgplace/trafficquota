# Traffic Quota

[![Build Status](https://travis-ci.com/orgplace/trafficquota.svg?branch=master)](https://travis-ci.com/orgplace/trafficquota)
[![Go Report Card](https://goreportcard.com/badge/github.com/orgplace/trafficquota)](https://goreportcard.com/report/github.com/orgplace/trafficquota)
[![codecov](https://codecov.io/gh/orgplace/trafficquota/branch/master/graph/badge.svg)](https://codecov.io/gh/orgplace/trafficquota)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Forgplace%2Ftrafficquota.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Forgplace%2Ftrafficquota?ref=badge_shield)
[![GoDoc](https://godoc.org/github.com/orgplace/trafficquota?status.svg)](https://godoc.org/github.com/orgplace/trafficquota)
[![Docker Pulls](https://img.shields.io/docker/pulls/orgplace/trafficquota.svg?style=flat)](https://hub.docker.com/r/orgplace/trafficquota)
[![packagecloud deb](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/orgplace/trafficquota?filter=debs)
[![packagecloud rpm](https://img.shields.io/badge/rpm-packagecloud.io-844fec.svg)](https://packagecloud.io/orgplace/trafficquota?filter=rpms)

[Token bucket](https://en.wikipedia.org/wiki/Token_bucket) server to control web API requests.
You can check whether the request conforms to defined limits on rate and burstiness by querying this service at the beginning of your API.

## Quick Start

### Starting `trafficquotad`

```sh
go run ./cmd/trafficquotad
```

### Using `trafficquotad` from Your Application

After `go get github.com/orgplace/trafficquota/client`
or `dep ensure -add github.com/orgplace/trafficquota/client`:

```go
import "github.com/orgplace/trafficquota/client"

c, _ := client.NewInsecureClient("localhost:3895")
allowed, _ := c.Take("key")
```

Please see [examples directory](examples) and [godoc](https://godoc.org/github.com/orgplace/trafficquota/client) for more detail.

## Installation

### Docker

Pull image from [Docker Hub](https://hub.docker.com/r/orgplace/trafficquota) and run:

```sh
docker run -it --rm -p3895:3895 orgplace/trafficquota:latest
```

In your `docker-compose.yml`:

```yml
services:
  trafficquota:
    image: orgplace/trafficquota
    ports:
    - "3895:3895"
```

### DEB/RPM Repository

You can use [DEB/RPM repository (packagecloud)](https://packagecloud.io/orgplace/trafficquota).
Currently, Ubuntu 18.04, Ubuntu 18.10, Debian 9, Fedora 29 and RHEL 7 are supprted.

To register repository, [follow packagecloud instruction](https://packagecloud.io/orgplace/trafficquota/install) ([deb](https://packagecloud.io/orgplace/trafficquota/install#bash-deb), [rpm](https://packagecloud.io/orgplace/trafficquota/install#bash-rpm)).
After registration:

```sh
# Ubuntu/Debian
apt-get install trafficquota
# Fedora
dnf install trafficquota
# RHEL
yum install trafficquota
```

When you use the package, systemd unit file for `trafficquotad` is also installed.

```sh
# Start
sudo systemctl start trafficquotad
# Enable
sudo systemctl enable trafficquotad
# Customize
sudo systemctl edit trafficquotad
# Log
journalctl -xeu trafficquotad
```

### From tarball/Build Source

Download tarball form [GitHub Releases](https://github.com/orgplace/trafficquota/releases) or build from source:

```sh
go build -o trafficquotad ./cmd/trafficquotad
```

## Monitoring

### Health Check

`trafficquotad` supports [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
You can query health of the server using [grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe).

```sh
grpc-health-probe -addr=localhost:3895
```
