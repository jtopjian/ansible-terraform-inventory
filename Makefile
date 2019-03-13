CURDIR := $(shell pwd)
TEST?=$$(go list ./... |grep -v 'vendor')
TARGETS=darwin linux
TERRAFORM_VERSION="0.11.13"

test:
	rm -rf work || true
	mkdir work ; \
	cd work ; \
	wget https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_linux_amd64.zip ; \
	unzip terraform_$(TERRAFORM_VERSION)_linux_amd64.zip
	PATH=$(CURDIR)/work:$(PATH) go test -v $(TEST) ./... -count=1

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

targets: $(TARGETS)

$(TARGETS):
	GOOS=$@ GOARCH=amd64 go build -o "dist/$@/terraform-provider-sensu_${TRAVIS_TAG}_x4"
	zip -j dist/terraform-provider-sensu_${TRAVIS_TAG}_$@_amd64.zip dist/$@/terraform-provider-sensu_${TRAVIS_TAG}_x4
