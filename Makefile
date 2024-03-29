CURDIR := $(shell pwd)
TEST?=$$(go list ./... |grep -v 'vendor')
TERRAFORM_VERSION_v011="0.11.13"
TERRAFORM_VERSION_v012="0.12.0"
ARCH=$(shell uname -s | tr A-Z a-z)

test: test_v011 test_v012

test_v011:
	rm -rf work || true
	mkdir work ; \
	cd work ; \
	wget https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION_v011)/terraform_$(TERRAFORM_VERSION_v011)_$(ARCH)_amd64.zip ; \
	unzip terraform_$(TERRAFORM_VERSION_v011)_$(ARCH)_amd64.zip
	PATH=$(CURDIR)/work:$(PATH) go test -v -run="V011" ./... -count=1

test_v012:
	rm -rf work || true
	mkdir work ; \
	cd work ; \
	wget https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION_v012)/terraform_$(TERRAFORM_VERSION_v012)_$(ARCH)_amd64.zip ; \
	unzip terraform_$(TERRAFORM_VERSION_v012)_$(ARCH)_amd64.zip
	PATH=$(CURDIR)/work:$(PATH) go test -v -run "V012" ./... -count=1

build:
	go install

fmtcheck:
	echo "==> Checking that code complies with gofmt requirements..."
	files=$$(find . -name '*.go' | grep -v 'vendor' ) ; \
	gofmt_files=`gofmt -l $$files`; \
	if [ -n "$$gofmt_files" ]; then \
		echo 'gofmt needs running on the following files:'; \
		echo "$$gofmt_files"; \
		echo "You can use the command: \`make fmt\` to reformat code."; \
		exit 1; \
	fi

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
