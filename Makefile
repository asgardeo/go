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
	@mkdir -p pkg/application/internal
	@mkdir -p pkg/api_resource/internal
	@mkdir -p pkg/identity_provider/internal
	@mkdir -p pkg/authenticator/internal
	@mkdir -p pkg/claim/internal
	@mkdir -p pkg/oidc_scope/internal
	@oapi-codegen \
		-package internal \
		-generate types,client \
		-o pkg/application/internal/client_gen.go \
		api-specs/applications.yaml
	@oapi-codegen \
		-package internal \
		-generate types,client \
		-o pkg/api_resource/internal/client_gen.go \
		api-specs/api_resources.yaml
	@oapi-codegen \
		-package internal \
		-generate types,client \
		-o pkg/identity_provider/internal/client_gen.go \
		api-specs/idp.yaml
	@oapi-codegen \
		-package internal \
		-generate types,client \
		-o pkg/authenticator/internal/client_gen.go \
		api-specs/authenticators.yaml
	@oapi-codegen \
		-package internal \
		-generate types,client \
		-o pkg/claim/internal/client_gen.go \
		api-specs/claim-management.yaml
	@oapi-codegen \
		-package internal \
		-generate types,client \
		-o pkg/oidc_scope/internal/client_gen.go \
		api-specs/oidc-scope-management.yaml
	@echo "SDK generation complete."

# Clean generated code
clean:
	@echo "Cleaning generated SDK code..."
	@rm -f pkg/application/internal/client_gen.go
	@rm -f pkg/api_resource/internal/client_gen.go
	@rm -f pkg/identity_provider/internal/client_gen.go
	@rm -f pkg/authenticator/internal/client_gen.go
	@rm -f pkg/claim/internal/client_gen.go
	@rm -f pkg/oidc_scope/internal/client_gen.go
	@echo "Cleaning complete."
