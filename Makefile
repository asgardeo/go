# Makefile for SDK generation

.PHONY: generate clean all install-tools

# Generate all
all: install-tools generate

# Install required tools
install-tools:
	@echo "Installing required tools..."
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Generate SDK code from OpenAPI specs
generate:
	@echo "Generating SDK code from OpenAPI specs..."
	@mkdir -p pkg/application
	@mkdir -p pkg/api_resource
	@oapi-codegen \
		-package application \
		-generate types,client \
		-o pkg/application/client_gen.go \
		api-specs/applications.yaml
	@oapi-codegen \
		-package api_resource \
		-generate types,client \
		-o pkg/api_resource/client_gen.go \
		api-specs/api_resources.yaml
	@echo "SDK generation complete."

# Clean generated code
clean:
	@echo "Cleaning generated SDK code..."
	@rm -f pkg/application/client_gen.go
	@rm -f pkg/api_resource/client_gen.go
	@echo "Cleaning complete."
