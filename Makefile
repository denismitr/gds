TIMESTAMP := $$(date +%Y%m%d%H%M%S)

deps:
	go mod tidy
	go mod vendor

test:
	go test ./...

analyze/contains:
	go test ./daryheap -bench=BenchmarkDaryHeap_Contains -benchmem -run=xxx -cpuprofile ./pprof/contains_cpu.pprof -memprofile ./pprof/contains_mem.pprof -benchtime=20s > ./bench/contains_$(TIMESTAMP).bench