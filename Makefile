.PHONY: build run-http run-socks5 run test clean

ifeq ($(OS),Windows_NT)
    RM = del /Q
    RMDIR = rmdir /Q /S
    MKDIR = mkdir
    ECHO = echo
    SLASH = \\
    EXE_EXT = .exe
else
    RM = rm -f
    RMDIR = rm -rf
    MKDIR = mkdir -p
    ECHO = echo
    SLASH = /
    EXE_EXT =
endif

BIN_DIR = bin
TEST_BIN_DIR = test_bin
TARGET = $(BIN_DIR)$(SLASH)proxy-checker$(EXE_EXT)
TEST_TARGET = $(TEST_BIN_DIR)$(SLASH)proxy-checker$(EXE_EXT)

build:
	go build -o $(TARGET) ./cmd/app/

run-http:
	@$(TARGET) --input $(BIN_DIR)$(SLASH)proxies-http.txt --type http --target https://www.google.com --timeout 5 --save $(BIN_DIR)$(SLASH)good-http.csv

run-socks5:
	@$(TARGET) --input $(BIN_DIR)$(SLASH)proxies-socks5.txt --type socks5 --target https://www.google.com --timeout 5 --save $(BIN_DIR)$(SLASH)good-socks5.csv

run: run-http run-socks5

test:
	@cd test && docker-compose up --build -d
	@$(MKDIR) $(TEST_BIN_DIR)
	
	@go build -o $(TEST_TARGET) ./cmd/app/
	
	@$(ECHO) 127.0.0.1:8888 > $(TEST_BIN_DIR)$(SLASH)proxies-http.txt
	@$(ECHO) 127.0.0.1:8889:user:pass >> $(TEST_BIN_DIR)$(SLASH)proxies-http.txt
	@$(ECHO) 127.0.0.1:1080 > $(TEST_BIN_DIR)$(SLASH)proxies-socks5.txt
	@$(ECHO) 127.0.0.1:1081:test:secret >> $(TEST_BIN_DIR)$(SLASH)proxies-socks5.txt
	
	@$(TEST_TARGET) --input $(TEST_BIN_DIR)$(SLASH)proxies-http.txt --type http --target https://www.google.com --timeout 5 --save $(TEST_BIN_DIR)$(SLASH)good-http.csv
	@$(TEST_TARGET) --input $(TEST_BIN_DIR)$(SLASH)proxies-socks5.txt --type socks5 --target https://www.google.com --timeout 5 --save $(TEST_BIN_DIR)$(SLASH)good-socks5.csv
	
	@cd test && docker-compose down
	@$(RMDIR) $(TEST_BIN_DIR)

clean:
	@$(RM) $(TARGET) $(BIN_DIR)$(SLASH)*.csv
	@$(RMDIR) $(TEST_BIN_DIR)

DEFAULT_GOAL := build