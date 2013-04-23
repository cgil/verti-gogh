all: server game capture

pi: server game
pi: GOENV = GOARCH=arm GOOS=linux GOARM=5

server:
	cd server && $(GOENV) go build

game:
	cd game && $(GOENV) go build

capture:
	make -C capture

transfer: pi
	scp ./game/game pi:code/verti-gogh/game
	scp ./server/server ./server/*.html pi:code/verti-gogh/server


.PHONY: server game capture
