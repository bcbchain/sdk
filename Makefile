
all: dist
.PHONY: all

dist:
	@sh -c "'$(CURDIR)/dist.sh'"

.PHONY: dist