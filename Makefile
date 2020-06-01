
all: dist
.PHONY: all

dist:
	@sh -c "'$(CURDIR)/scripts/dist.sh'"

.PHONY: dist