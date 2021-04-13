ARTIFACT_ID=tempdel
VERSION=0.1.0

MAKEFILES_VERSION=4.5.0

.DEFAULT_GOAL:=compile

include build/make/variables.mk

include build/make/self-update.mk
include build/make/info.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
