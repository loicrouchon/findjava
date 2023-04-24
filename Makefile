.PHONY: build clean run

BUILD_DIR=build
JAVA_BUILD_DIR=$(BUILD_DIR)/classes
JAVA_INFO=$(JAVA_BUILD_DIR)/JvmInfo.class
GO_BUILD_DIR=$(BUILD_DIR)/go
MAIN_PROGRAM=$(GO_BUILD_DIR)/jvm-finder

build: $(JAVA_INFO) $(MAIN_PROGRAM)

$(JAVA_INFO): JvmInfo.java
	@mkdir -p "$(JAVA_BUILD_DIR)"
	@javac --release 8 -d "$(JAVA_BUILD_DIR)" JvmInfo.java

$(MAIN_PROGRAM): jvm-finder
	@mkdir -p "$(GO_BUILD_DIR)"
	@cd jvm-finder && go build -ldflags "-s -w" -o "../$(MAIN_PROGRAM)" jvm-finder

clean:
	@rm  -rf "$(BUILD_DIR)"

run: $(JAVA_INFO)
	@cd jvm-finder && go run jvm-finder
