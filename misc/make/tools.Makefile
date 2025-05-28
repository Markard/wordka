# This makefile should be used to hold functions/variables
ifeq ($(ARCH),x86_64)
	ARCH := amd64
else ifeq ($(ARCH),aarch64)
	ARCH := arm64
endif

define github_url
    https://github.com/$(GITHUB)/releases/download/v$(VERSION)/$(ARCHIVE)
endef

# creates a directory bin.
bin:
	@ mkdir -p $@

# ~~~ Tools ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# ~~ [migrate] ~~~ https://github.com/golang-migrate/migrate ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

MIGRATE := $(shell command -v migrate || echo "bin/migrate")
migrate: bin/migrate ## Install migrate (database migration)

bin/migrate: VERSION := 4.18.3
bin/migrate: GITHUB  := golang-migrate/migrate
bin/migrate: ARCHIVE := migrate.$(OSTYPE)-$(ARCH).tar.gz
bin/migrate: bin
	@ printf "Install migrate from $(call github_url)... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - migrate > $@ && chmod +x $@
	@ echo "done."

# ~~ [ air ] ~~~ https://github.com/cosmtrek/air ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

AIR := $(shell command -v air || echo "bin/air")
air: bin/air ## Installs air (go file watcher)

bin/air: VERSION := 1.61.7
bin/air: GITHUB  := cosmtrek/air
bin/air: ARCHIVE := air_$(VERSION)_$(OSTYPE)_$(ARCH).tar.gz
bin/air: bin
	@ printf "Install air from $(call github_url)... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - air > $@ && chmod +x $@
	@ echo "done."


# ~~ [ gotestsum ] ~~~ https://github.com/gotestyourself/gotestsum ~~~~~~~~~~~~~~~~~~~~~~~

GOTESTSUM := $(shell command -v gotestsum || echo "bin/gotestsum")
gotestsum: bin/gotestsum ## Installs gotestsum (testing go code)

bin/gotestsum: VERSION := 1.12.2
bin/gotestsum: GITHUB  := gotestyourself/gotestsum
bin/gotestsum: ARCHIVE := gotestsum_$(VERSION)_$(OSTYPE)_$(ARCH).tar.gz
bin/gotestsum: bin
	@ printf "Install gotestsum from $(call github_url)... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - gotestsum > $@ && chmod +x $@
	@ echo "done."

# ~~ [ tparse ] ~~~ https://github.com/mfridman/tparse ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

TPARSE := $(shell command -v tparse || echo "bin/tparse")
tparse: bin/tparse ## Installs tparse (testing go code)

# eg https://github.com/mfridman/tparse/releases/download/v0.13.2/tparse_darwin_arm64
bin/tparse: VERSION := 0.17.0
bin/tparse: GITHUB  := mfridman/tparse
bin/tparse: ARCHIVE := tparse_$(OSTYPE)_$(ARCH)
bin/tparse: bin
	@ printf "Install tparse from $(call github_url)... "
	@ curl -Ls $(call github_url) > $@ && chmod +x $@
	@ echo "done."

# ~~ [ mockery ] ~~~ https://github.com/vektra/mockery ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

MOCKERY := $(shell command -v mockery || echo "bin/mockery")
mockery: bin/mockery ## Installs mockery (mocks generation)

bin/mockery: VERSION := 3.3.0
bin/mockery: GITHUB  := vektra/mockery
bin/mockery: ARCHIVE := mockery_$(VERSION)_$(OSTYPE)_$(ARCH).tar.gz
bin/mockery: bin
	@ printf "Install mockery from $(call github_url)... "
	@ curl -Ls $(call github_url) | tar -zOxf - mockery > $@ && chmod +x $@
	@ echo "done."

# ~~ [ golangci-lint ] ~~~ https://github.com/golangci/golangci-lint ~~~~~~~~~~~~~~~~~~~~~

GOLANGCI := $(shell command -v golangci-lint || echo "bin/golangci-lint")
golangci-lint: bin/golangci-lint ## Installs golangci-lint (linter)

bin/golangci-lint: VERSION := 2.1.6
bin/golangci-lint: GITHUB  := golangci/golangci-lint
bin/golangci-lint: ARCHIVE := golangci-lint-$(VERSION)-$(OSTYPE)-$(ARCH).tar.gz
bin/golangci-lint: bin
	@ printf "Install golangci-linter from $(call github_url)... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - golangci-lint > $@ && chmod +x $@
	@ echo "done."

# ~~ [ testfixtures ] ~~~ https://github.com/go-testfixtures/testfixtures ~~~~~~~~~~~~~~~~~~~~~

GOLANGCI := $(shell command -v testfixtures || echo "bin/testfixtures")
testfixtures: bin/testfixtures ## Installs testfixtures

bin/testfixtures: VERSION := 3.16.0
bin/testfixtures: GITHUB  := go-testfixtures/testfixtures
bin/testfixtures: ARCHIVE := testfixtures_$(OSTYPE)_$(ARCH).tar.gz
bin/testfixtures: bin
	@ printf "Install testfixtures from $(call github_url)... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - testfixtures > $@ && chmod +x $@
	@ echo "done."