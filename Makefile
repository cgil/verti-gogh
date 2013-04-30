all: server capture

pi: server
pi: GOENV = GOARCH=arm GOOS=linux GOARM=5

server:
	$(GOENV) go build

capture:
	make -C capture

transfer: pi
	scp ./verti-gogh pi:code/verti-gogh
	scp ./server/*.html pi:code/verti-gogh/server
	scp -r ./server/static pi:code/verti-gogh/server
	scp ./capture/*.c ./capture/*.h pi:code/verti-gogh/capture


.PHONY: server game capture
