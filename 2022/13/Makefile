GOC = go build
GOV = go vet
GOL = golint

BENCH = ../bench.sh

EX = sample.txt
IN = input.txt

BIN = $(addprefix aoc,$(shell basename $(CURDIR)))
SRC = $(BIN).go

bench:
	@$(MAKE) build
	$(BENCH) $(BIN)
	@$(MAKE) clean

build:
	$(GOC) $(SRC)

check:
	$(GOV) $(SRC)
	$(GOL) $(SRC)

clean:
	go clean
	rm -f $(BIN)

cpuprof: build
	cat $(IN) | ./$(BIN) -cpuprofile=$(BIN).cpu.prof

memprof: build
	cat $(IN) | ./$(BIN) -memprofile=$(BIN).mem.prof

run:
	cat $(IN) | go run ./$(SRC)

sample:
	cat $(EX) | go run ./$(SRC)



.PHONY: bench build check clean exemple run