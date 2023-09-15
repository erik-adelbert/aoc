GOC = go build
GOV = go vet
GOL = golint

BENCH = ../bench.sh

DOWNLOAD = ../download/main.go

EX = sample.txt
IN = input.txt

BIN = $(addprefix aoc,$(shell basename $(CURDIR)))
SRC = $(BIN).go

bench: build
	$(BENCH) $(BIN)
	@$(MAKE) clean

build: input.txt
	$(GOC) $(SRC)

check:
	$(GOV) $(SRC)
	$(GOL) $(SRC)

clean:
	go clean
	rm -f $(BIN)

input.txt:
	go run $(DOWNLOAD)

cpuprof: build
	./$(BIN) -cpuprofile=$(BIN).cpu.prof < $(IN)

memprof: build
	./$(BIN) -memprofile=$(BIN).mem.prof < $(IN)

run: input.txt
	go run ./$(SRC) < $(IN)

sample:
	go run ./$(SRC) < $(EX)


.PHONY: bench build check clean cpuprof exemple memprof run