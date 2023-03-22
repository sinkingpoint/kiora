#!/bin/sh

docker run --name kiora-prometheus --net host -d -v $(pwd)/testdata/prometheus:/etc/prometheus prom/prometheus
