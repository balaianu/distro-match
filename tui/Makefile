.PHONY: build run serve clean distros

build: distros.json
	go build -o distro-match .

run: distros.json
	go run .

serve: distros.json
	go run . --serve :2323 --host-key .ssh/host_ed25519

distros.json: ../commons/data/distros.json
	cp ../commons/data/distros.json distros.json

clean:
	rm -f distros.json distro-match distro-match-*
	rm -rf .ssh

# Cross-compilation targets — builds binaries for other platforms.
# Usage: make dist  or  make dist-linux-amd64 etc.
DIST_DIR := dist
TARGETS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

dist: $(TARGETS)

$(TARGETS):
	@mkdir -p $(DIST_DIR)
	$(eval GOOS := $(word 1,$(subst -, ,$@)))
	$(eval GOARCH := $(word 2,$(subst -, ,$@)))
	$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(DIST_DIR)/distro-match-$(GOOS)-$(GOARCH)$(EXT) .

dist-clean:
	rm -rf $(DIST_DIR)
