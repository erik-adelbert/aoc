GO = go build 

TARGETS := all bench build check clean run 

SUBDIRS := $(wildcard */.)

$(TARGETS): $(SUBDIRS)
$(SUBDIRS):
		@$(MAKE) -C $@ $(MAKECMDGOALS)

cyclo:
	gocyclo -top 10 .

runtime:
	./runtime.sh | sort -n -k 4 > runtime.md

.PHONY: $(TARGETS) $(SUBDIRS) cyclo runtime