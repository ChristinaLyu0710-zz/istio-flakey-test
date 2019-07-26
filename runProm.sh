#!/bin/bash

pushd "${GOPATH}/src/istio.io/tools/perf/benchmark/runner" > /dev/null
sleep 1m
python prom.py http://localhost:8060 60 --no-aggregate
popd > /dev/null
