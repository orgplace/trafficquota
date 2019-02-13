# Traffic Quota

[![Build Status](https://travis-ci.com/orgplace/trafficquota.svg?branch=master)](https://travis-ci.com/orgplace/trafficquota)
[![Go Report Card](https://goreportcard.com/badge/github.com/orgplace/trafficquota)](https://goreportcard.com/report/github.com/orgplace/trafficquota)
[![codecov](https://codecov.io/gh/orgplace/trafficquota/branch/master/graph/badge.svg)](https://codecov.io/gh/orgplace/trafficquota)
[![GoDoc](https://godoc.org/github.com/orgplace/trafficquota?status.svg)](https://godoc.org/github.com/orgplace/trafficquota)
[![Docker Pulls](https://img.shields.io/docker/pulls/orgplace/trafficquota.svg?style=flat)](https://hub.docker.com/r/orgplace/trafficquota)

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

Please see `examples` directory and godoc for more detail.

## Health Checking

`trafficquotad` supports [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
You can query health of the server using [grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe).

```sh
grpc-health-probe -addr=localhost:3895
```
