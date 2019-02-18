# Traffic Quota Example HTTP Server

Running `trafficquotad` on `localhost:3895`, start this server:

```sh
go run .
```

And open another terminal:

```sh
# 20 (or little more) requests are accepted at default.
# And other requests are rejected by the burst limit.
ab -c100 -n100 http://localhost:8080/
# Each request is accepted because these are in the rate limit.
watch -n0.1 ab -c10 -n10 http://localhost:8080/
# 9 (or little less) reqests will exceed the rate limit.
watch -n0.1 ab -c20 -n20 http://localhost:8080/
```
