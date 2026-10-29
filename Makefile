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
	@mkdir -p pkg/consent/internal
	@mkdir -p pkg/self_service/internal/application
	@mkdir -p pkg/self_service/internal/association
	@mkdir -p pkg/self_service/internal/authorized_apps
	@mkdir -p pkg/self_service/internal/backup_code
	@mkdir -p pkg/self_service/internal/consent
	@mkdir -p pkg/self_service/internal/fido
	@mkdir -p pkg/self_service/internal/mfa
	@mkdir -p pkg/self_service/internal/push
	@mkdir -p pkg/self_service/internal/session
	@mkdir -p pkg/self_service/internal/totp
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
	@oapi-codegen \
		-package internal \
		-generate types,client \
		-o pkg/consent/internal/client_gen.go \
		api-specs/carbon-consent-management.yaml
	@oapi-codegen \
		-package application \
		-generate types,client \
		-o pkg/self_service/internal/application/client_gen.go \
		api-specs/self_service/application.yaml
	@oapi-codegen \
		-package association \
		-generate types,client \
		-o pkg/self_service/internal/association/client_gen.go \
		api-specs/self_service/association.yaml
	@oapi-codegen \
		-package authorized_apps \
		-generate types,client \
		-o pkg/self_service/internal/authorized_apps/client_gen.go \
		api-specs/self_service/authorizedApps.yaml
	@oapi-codegen \
		-package backup_code \
		-generate types,client \
		-o pkg/self_service/internal/backup_code/client_gen.go \
		api-specs/self_service/backupcode.yaml
	@oapi-codegen \
		-package consent \
		-generate types,client \
		-o pkg/self_service/internal/consent/client_gen.go \
		api-specs/self_service/consent.yaml
	@oapi-codegen \
		-package fido \
		-generate types,client \
		-o pkg/self_service/internal/fido/client_gen.go \
		api-specs/self_service/fido.yaml
	@oapi-codegen \
		-package mfa \
		-generate types,client \
		-o pkg/self_service/internal/mfa/client_gen.go \
		api-specs/self_service/mfa.yaml
	@oapi-codegen \
		-package push \
		-generate types,client \
		-o pkg/self_service/internal/push/client_gen.go \
		api-specs/self_service/push.yaml
	@oapi-codegen \
		-package session \
		-generate types,client \
		-o pkg/self_service/internal/session/client_gen.go \
		api-specs/self_service/session.yaml
	@oapi-codegen \
		-package totp \
		-generate types,client \
		-o pkg/self_service/internal/totp/client_gen.go \
		api-specs/self_service/totp.yaml
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
	@rm -f pkg/consent/internal/client_gen.go
	@rm -f pkg/self_service/internal/application/client_gen.go
	@rm -f pkg/self_service/internal/association/client_gen.go
	@rm -f pkg/self_service/internal/authorized_apps/client_gen.go
	@rm -f pkg/self_service/internal/backup_code/client_gen.go
	@rm -f pkg/self_service/internal/consent/client_gen.go
	@rm -f pkg/self_service/internal/fido/client_gen.go
	@rm -f pkg/self_service/internal/mfa/client_gen.go
	@rm -f pkg/self_service/internal/push/client_gen.go
	@rm -f pkg/self_service/internal/session/client_gen.go
	@rm -f pkg/self_service/internal/totp/client_gen.go
	@echo "Cleaning complete."
