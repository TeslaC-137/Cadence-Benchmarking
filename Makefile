.PHONY: test bins clean
PROJECT_ROOT = github.com/samarabbas/cadence-samples

export PATH := $(GOPATH)/bin:$(PATH)

# default target
default: test

PROGS = sequence-workflow \

TEST_ARG ?= -race -v -timeout 5m
BUILD := ./build
SAMPLES_DIR=./workflows

export PATH := $(GOPATH)/bin:$(PATH)

# Automatically gather all srcs
ALL_SRC := $(shell find ./samples/common -name "*.go")

# all directories with *_test.go files in them
# TEST_DIRS=./sequence-workflow

dep-ensured:
	dep ensure


sequence-workflow: dep-ensured $(ALL_SRC)
	go build -i -o bin/sequence samples/workflows/sequence-workflow/*.go


bins: sequence-workflow \


test: bins
	@rm -f test
	@rm -f test.log
	@echo $(TEST_DIRS)
	@for dir in $(TEST_DIRS); do \
		go test -coverprofile=$@ "$$dir" | tee -a test.log; \
	done;

clean:
	rm -rf bin
	rm -Rf $(BUILD)
