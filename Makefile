ARTIFACT_ID=tempdel
VERSION=0.1.0

MAKEFILES_VERSION=9.0.4
GOTAG=1.22.4
LINT_VERSION=v1.57.2
# overwrite ADDITIONAL_LDFLAGS to disable static compilation
# this should fix https://github.com/golang/go/issues/13470
ADDITIONAL_LDFLAGS=""
GO_ENVIRONMENT=GO111MODULE=on
GO_ENV_VARS=GOPRIVATE=CGO_ENABLED=0

.DEFAULT_GOAL:=compile

include build/make/variables.mk

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
