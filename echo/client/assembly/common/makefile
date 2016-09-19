.PHONY: all dep dir bin test clean show
.SUFFIXES: .go

SRC_DIR		:= src
APP_DIR     := app
OUTPUT_DIR	:= release

DEBUG		:= true
BUILD		:= GOOS=linux GOARCH=amd64 go build
CLEAN		:= GOOS=linux GOARCH=amd64 go clean
RUN			:= GOOS=linux GOARCH=amd64 go run
INSTALL 	:= GOOS=linux GOARCH=amd64 go install
# go test -test.bench=".*"表示测试全部的压力测试函数
TEST		:= GOOS=linux GOARCH=amd64 go test -v
TYPE		:= GOOS=linux GOARCH=amd64 gotype  -e -a -v dir
# 查看变量分配的位置，如实在stack还是在堆上
WATCH       := go build -gcflags -m
# 查看是否有同名变量冲突
SHADOW      := go tool vet -shadow

ifeq ($(DEBUG), true)
	# GFLAGS  := -gcflags "-N -l" -race -x -v # -x会把go build的详细过程输出
	# GFLAGS  := -gcflags "-N -l" -race -v
	GFLAGS  := -gcflags "-N -l" -v
	OUTPUT_DIR := debug
else
	GFLAGS  := -ldflags "-s -w"
	OUTPUT_DIR := release
endif

BIN			:= client
all: dir clean $(BIN)

dep:
	export GOPATH=${PWD}:${GOPATH}

dir:
	@-mkdir -p $(OUTPUT_DIR)

client:
	# export GOPATH=${PWD}:${GOPATH} && cd $(SRC_DIR) && $(BUILD) $(GFLAGS) -o ../$(OUTPUT_DIR)/$@ app/client.go
	cd $(APP_DIR) && $(BUILD) $(GFLAGS) -o ./$(OUTPUT_DIR)/$@

test:
	export GOPATH=${PWD}:${GOPATH} && cd $(SRC_DIR) && $(TEST)

bench:
	export GOPATH=${PWD}:${GOPATH} && cd $(SRC_DIR)/consumer && $(TEST) consumer -bench ^BenchmarkConsumer_Call$ -run ^$ -benchtime=2s -cpuprofile=prof.cpu

clean:
	-$(CLEAN)
	-rm -rf $(OUTPUT_DIR)

install:
	-$(INSTALL)

show:
	@-echo $(OUTPUT_DIR)

