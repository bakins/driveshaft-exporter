driveshaft-exporter
================

Export [driveshaft](https://github.com/keyurdg/driveshaft) metrics in [Prometheus](https://prometheus.io/) format.

See [Releases](https://github.com/bakins/driveshaft-exporter/releases) for pre-built binaries.

Build
=====

Requires [Go](https://golang.org/doc/install). Tested with Go 1.8+.

Clone this repo into your `GOPATH` (`$HOME/go` by default) and run build:

```
mkdir -p $HOME/go/src/github.com/bakins
cd $HOME/go/src/github.com/bakins
git clone https://github.com/bakins/driveshaft-exporter
cd driveshaft-exporter
./script/build
```

You should then have two executables: driveshaft-exporter.linux.amd64 and driveshaft-exporter.darwin.amd64

You may want to rename for your local OS, ie `mv driveshaft-exporter.darwin.amd64 driveshaft-exporter`

Running
=======

```
./driveshaft-exporter --help
driveshaft metrics exporter

Usage:
  driveshaft-exporter [flags]

Flags:
      --addr string       listen address for metrics handler (default "127.0.0.1:8080")
      --driveshaft string   address of driveshaft (default "127.0.0.1:4730")
```

When running, a simple healthcheck is availible on `/healthz`

Metrics
=======

Metrics will be exposes on `/metrics`

```
curl http://localhost:8080/metrics

# HELP driveshaft_threads_count count of threads
# TYPE driveshaft_threads_count gauge
driveshaft_threads_count{function="FakeSnaps"} 4
# HELP driveshaft_up is driveshaft up
# TYPE driveshaft_up gauge
driveshaft_up 1
```

LICENSE
========

See [LICENSE](./LICENSE)
