.PHONY: clean all format test build run

BUILD_DIR=build
JAVA_BUILD_DIR=$(BUILD_DIR)/classes
JAVA_INFO=$(JAVA_BUILD_DIR)/JvmInfo.class
GO_BUILD_DIR=$(BUILD_DIR)/go
MAIN_PROGRAM=$(GO_BUILD_DIR)/jvm-finder

all: format test build

build: $(JAVA_INFO) $(MAIN_PROGRAM)

test: $(JAVA_INFO)
	@go test

$(JAVA_INFO): JvmInfo.java
	@mkdir -p "$(JAVA_BUILD_DIR)"
	@javac --release 8 -d "$(JAVA_BUILD_DIR)" JvmInfo.java

$(MAIN_PROGRAM): *.go
	@mkdir -p "$(GO_BUILD_DIR)"
	@go build -ldflags "-s -w" -o "$(MAIN_PROGRAM)" jvm-finder

format:
	@go fmt

clean:
	@rm  -rf "$(BUILD_DIR)"

run: $(JAVA_INFO)
	@go run jvm-finder
