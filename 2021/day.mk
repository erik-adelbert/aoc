GOC = go build
GOV = go vet

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

clean:
	go clean
	rm -f $(BIN)

input.txt:
	go run $(DOWNLOAD)

go.mod: 
	@go mod init 2>/dev/null

gobench: go.mod input.txt
	go test -bench=. -benchmem

cpuprof: build
	./$(BIN) -cpuprofile=$(BIN).cpu.prof < $(IN)

memprof: build
	./$(BIN) -memprofile=$(BIN).mem.prof < $(IN)

run: input.txt
	go run ./$(SRC) < $(IN)

sample:
	go run ./$(SRC) < $(EX)


.PHONY: bench build check clean cpuprof exemple gobench memprof run sample