default: install

# install
install:
	@go install

# unit tests
test:
	@test/unit/testAllPackages.sh

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
	@docker build --rm=true --file=test/integration/Dockerfile --tag=golook:latest .

integration: fmt vet
	@test/unit/testAllPackages.sh

.PHONY: \
	install \
	test \
	fmt \
	vet \
	docker