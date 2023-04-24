.PHONY: build clean run

BUILD_DIR=build
JAVA_BUILD_DIR=$(BUILD_DIR)/classes
JAVA_INFO=$(JAVA_BUILD_DIR)/JvmInfo.class
GO_BUILD_DIR=$(BUILD_DIR)/go
MAIN_PROGRAM=$(GO_BUILD_DIR)/main

build: $(JAVA_INFO) $(MAIN_PROGRAM)

$(JAVA_INFO): JvmInfo.java
	@mkdir -p "$(JAVA_BUILD_DIR)"
	javac --release 8 -d "$(JAVA_BUILD_DIR)" JvmInfo.java

$(MAIN_PROGRAM): main.go
	@mkdir -p "$(GO_BUILD_DIR)"
	go build -ldflags "-w" -o "$(MAIN_PROGRAM)" main.go

clean:
	rm  -rf "$(BUILD_DIR)"

run: $(JAVA_INFO)
	@go run main.go
