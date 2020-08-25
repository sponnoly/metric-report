# metric-report

This repository holds the source code for the metric-report.

## Prerequisites

### Download Go

You must have Go version **1.13.x** installed on your machine before compiling the application.
You can download Go from <https://golang.org/dl/>.

### Go Modules

see the [go-modules instructions](doc/go-modules.md) for information on how to setup your system.

### GetMetricSum

curl --header "Content-Type: application/json" --request GET http://localhost:8080/metric/{key}/sum

### InsertMetric

curl --header "Content-Type: application/json" --request POST --data '{"value":2}' http://localhost:8080/metric/{key}
