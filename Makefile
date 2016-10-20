
REPO?=github.com/rydrman/spotr

ldflags=\
	-X $(REPO).clientID=$(SPOTIFY_CLIENT_ID) \
	-X $(REPO).clientSecret=$(SPOTIFY_CLIENT_SECRET)

build:
	go build -ldflags "$(ldflags)"

install: build
	go install

run: build install
	spotr



