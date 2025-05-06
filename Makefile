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
	@mkdir -p pkg/identity_provider
	@mkdir -p pkg/authenticator
	@mkdir -p pkg/claim
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
	@oapi-codegen \
		-package identity_provider \
		-generate types,client \
		-o pkg/identity_provider/client_gen.go \
		api-specs/idp.yaml
	@oapi-codegen \
		-package authenticator \
		-generate types,client \
		-o pkg/authenticator/client_gen.go \
		api-specs/authenticators.yaml
	@oapi-codegen \
		-package claim \
		-generate types,client \
		-o pkg/claim/client_gen.go \
		api-specs/claim-management.yaml
	@echo "SDK generation complete."

# Clean generated code
clean:
	@echo "Cleaning generated SDK code..."
	@rm -f pkg/application/client_gen.go
	@rm -f pkg/api_resource/client_gen.go
	@rm -f pkg/identity_provider/client_gen.go
	@rm -f pkg/authenticator/client_gen.go
	@rm -f pkg/claim/client_gen.go
	@echo "Cleaning complete."
