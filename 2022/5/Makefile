GOC = go build
GOV = go vet
GOL = golint

BENCH = ../bench.sh

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

run:
	cat $(IN) | go run ./$(SRC)

.PHONY: bench build check clean exemple run