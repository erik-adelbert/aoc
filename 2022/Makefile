GO = go build 

TARGETS := bench build check clean run

SUBDIRS := $(wildcard */.)
SUBDIRS := $(filter-out download/. runtime/., $(SUBDIRS))

$(TARGETS): $(SUBDIRS)
$(SUBDIRS):
		@$(MAKE) -C $@ $(MAKECMDGOALS)

cyclo:
	@gocyclo -top 10 -avg .

lines:
	@find . -name '*go' \( -not -iname "main.go" \) | grep -v 'v1' | xargs wc -l | sort

runtime:
	@./runtime.sh | sort -n -k 4 > runtime.md

.PHONY: $(TARGETS) $(SUBDIRS) cyclo lines runtime 
