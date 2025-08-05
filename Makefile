
.PHONY: test
test: ## Run unit tests for each module
	@hack/do.sh test

.PHONY: dagger-version
dagger-version: ## Get the latest dagger version used
	@hack/do.sh dagger_version