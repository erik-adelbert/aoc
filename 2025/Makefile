GO = go build

TARGETS := bench binrun build check clean header run

SUBDIRS := $(wildcard */.)
SUBDIRS := $(filter-out download/. header/. runtime/., images/., $(SUBDIRS))

$(TARGETS): $(SUBDIRS)
$(SUBDIRS):
		@$(MAKE) -C $@ $(MAKECMDGOALS)

cyclo:
	@gocyclo -ignore "download|images|header|runtime" -top 10 -avg .

lines:
	@find . -name '*go' \( -not -iname "main.go" \) | grep -v 'v1' | xargs wc -l | sort

runtime:
	@./runtime.sh | sort -n -k 4 > runtime.md

day:
	@day="$(filter-out day,$(MAKECMDGOALS))"; \
	if [ -z "$$day" ]; then \
		echo "Usage: make day <day_number>"; \
		echo "Example: make day 3"; \
		exit 1; \
	fi; \
	if [ -d "$$day" ]; then \
		echo "Directory $$day already exists"; \
		exit 1; \
	fi; \
	mkdir -p $$day && \
	echo "include ../day.mk" > $$day/Makefile && \
	echo "" >> $$day/aoc$$day.go && \
	echo "package main" >> $$day/aoc$$day.go && \
	echo "" >> $$day/aoc$$day.go && \
	echo "import (" >> $$day/aoc$$day.go && \
	echo "	\"fmt\"" >> $$day/aoc$$day.go && \
	echo ")" >> $$day/aoc$$day.go && \
	echo "" >> $$day/aoc$$day.go && \
	echo "func main() {" >> $$day/aoc$$day.go && \
	echo "	// TODO: Implement solution" >> $$day/aoc$$day.go && \
	echo "	fmt.Println(\"Day $$day - Not implemented yet\")" >> $$day/aoc$$day.go && \
	echo "}" >> $$day/aoc$$day.go && \
	$(MAKE) -C $$day header && \
	echo "Created day $$day structure:" && \
	echo "  - Directory: $$day/" && \
	echo "  - File: $$day/aoc$$day.go (with header)" && \
	echo "  - File: $$day/Makefile"

# Prevent make from treating day numbers as targets
%:
	@:

.PHONY: $(TARGETS) $(SUBDIRS) cyclo lines runtime day
