PACKAGES = $(shell ./scripts/packages.sh)

EXTERNAL_TOOLS = \
	github.com/golang/dep/cmd/dep \
	github.com/motemen/gobump \
	github.com/Songmu/ghch/cmd/ghch

setup:
	@for tool in $(EXTERNAL_TOOLS) ; do \
      echo "Installing $$tool" ; \
      go get $$tool; \
    done

test-all: vet lint test

test:
	./scripts/test.sh

vet:
	go vet $(PACKAGES)

lint:
	@if [ -z `which errcheck 2> /dev/null` ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	echo $(PACKAGES) | xargs -n 1 golint -set_exit_status

errcheck:
	@if [ -z `which errcheck 2> /dev/null` ]; then \
		go get -u github.com/kisielk/errcheck; \
	fi
	echo $(PACKAGES) | xargs errcheck -ignoretests

bump: setup
	./scripts/bumpup.sh

deploy: bump
	./scripts/deploy.sh

upload: bump
	./scripts/upload.sh

.PHONY: test-all test vet lint setup bump
