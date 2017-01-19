
REPO?=github.com/rydrman/itunes-to-spotify

ldflags=\
	-X main.clientID=$(SPOTIFY_CLIENT_ID) \
	-X main.clientSecret=$(SPOTIFY_CLIENT_SECRET)

build:
	go build -a -x -ldflags "$(ldflags)"

install:
	go install -ldflags "$(ldflags)"

re-install:
	go install -x -ldflags "$(ldflags)"

run: install
	itunes-to-spotify



