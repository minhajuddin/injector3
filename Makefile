.PHONY: all
all: gen build test bench

.PHONY: gen
gen:
	@echo "Generating/Updating bazel files..."
	@bazel run //:gazelle
	@bazel run @rules_go//go -- mod tidy
	@bazel mod tidy

.PHONY: build
build:
	@bazel build //...

.PHONY: test
test:
	@bazel test //... --test_output=errors --flaky_test_attempts=3

.PHONY: bench
bench:
	@bazel run //:injector3_test -- -test.bench=.
