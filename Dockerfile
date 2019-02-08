FROM golang:1.11 AS builder

ARG pkg_dir=/go/src/github.com/orgplace/trafficquota
ADD . ${pkg_dir}
WORKDIR ${pkg_dir}
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /trafficquotad ./cmd/trafficquotad

FROM scratch

COPY --from=builder /trafficquotad /trafficquotad

ENTRYPOINT [ "/trafficquotad" ]
