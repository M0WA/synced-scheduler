default: test benchmark

.PHONY: benchmark
benchmark: syncedsched-benchmark storagesched-benchmark

.PHONY: test
test: syncedsched-test storagesched-test

.PHONY: syncedsched-benchmark
syncedsched-benchmark:
	go test -bench=. ./pkgs/syncedscheduler/...

.PHONY: storagesched-benchmark
storagesched-benchmark:
	go test -bench=. ./pkgs/storagescheduler/...

.PHONY: syncedsched-test
syncedsched-test:
	go test -v ./pkgs/syncedscheduler/...

.PHONY: storagesched-test
storagesched-test:
	go test -v ./pkgs/storagescheduler/...