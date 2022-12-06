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

sample:
	cat $(EX) | go run ./$(SRC)

run:
	cat $(IN) | go run ./$(SRC)

.PHONY: bench build check clean exemple run