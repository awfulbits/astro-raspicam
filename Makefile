GOCMD=go
GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build
OS=linux
ARCH=arm
ARM=7
PATH_DIST=dist
BINARY_NAME=astro-raspicam

clean:
	$(GOCLEAN)
	rm -rf $(PATH_DIST)
build:
	make clean
	mkdir $(PATH_DIST)
	cp config.toml $(PATH_DIST)
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) GOARM=$(ARM) $(GOBUILD) -o $(PATH_DIST)/$(BINARY_NAME) -v
run:
	make build
	$(PATH_DIST)/$(BINARY_NAME)