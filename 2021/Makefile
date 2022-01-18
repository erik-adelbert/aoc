GO = go build 

TARGETS := all bench build check clean run 

SUBDIRS := $(wildcard */.)

$(TARGETS): $(SUBDIRS)
$(SUBDIRS):
		@$(MAKE) -C $@ $(MAKECMDGOALS)

runtime:
	$(GO) runtime.go
	./runtime.sh | sort -n -k 4 > runtime.md
	rm -f runtime

cyclo:
	gocyclo -top 10 .

.PHONY: $(TOPTARGETS) $(SUBDIRS) cyclo runtime