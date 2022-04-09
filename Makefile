build: ## build
	go build -o main
run: ## Run with pprof
	./main -cpuprofile cpu.prof -memprofile mem.prof

profmem: ## pprof analyze memory
	go tool pprof mem.prof

profcpu: ## pprof analyze cpu
	go tool pprof cpu.prof

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
