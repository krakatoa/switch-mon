build: ./src/switch_mon/switch_mon.go
	export GOPATH=`pwd`
	cd src/switch_mon; \
	godep go build -o ../../release/switch_mon

deps:
	go get github.com/tools/godep
	go get github.com/amir/raidman
	go get github.com/0x19/goesl
