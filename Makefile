GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOFMT := "goimports"

fmt: ## Run gofmt for all .go files
	$(GOFMT) -w $(GOFMT_FILES)

test: ## Run go test for whole project
	go test -v ./...

gendoc: ## Generate document in docs
	@gomarkdoc -o docs/datetime.md ./datetime
	@gomarkdoc -o docs/db.md ./db
	@gomarkdoc -o docs/generator.md ./generator
	@gomarkdoc -o docs/notifications.md ./notifications
	@gomarkdoc -o docs/number.md ./number
	@gomarkdoc -o docs/pagination.md ./pagination

goupdate: ## Remove go.sum and Go get dependencies
	@GOSUMDB=off rm -rf go.sum
	@GOSUMDB=off go mod tidy

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
