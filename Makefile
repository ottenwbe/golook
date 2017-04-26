default: install

# release
release:
	@go install -ldflags "-s -w"

# install
install:
	@go install

# unit tests
test:
	@test/unit/run_unit_tests.sh

fmt:
	@for d in $$(go list ./... | grep -v vendor); do \
		go fmt $${d}; \
	done

vet:
	@for d in $$(go list ./... | grep -v vendor); do \
		go vet $${d};  \
	done

lint:
	@for d in $$(go list ./... | grep -v vendor); do \
		golint $${d};  \
	done

# create a docker image, i.e., for integration tests
docker:
	@docker pull golang:1.8.1
	@docker build --rm=true --file=test/integration/Dockerfile --tag=golook:latest .

integration: fmt vet docker
	@test/integration/run_integration_tests.sh

.PHONY: \
	release \
	install \
	test \
	fmt \
	vet \
	lint \
	docker \
	integration