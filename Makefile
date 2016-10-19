build: ./src/switch_mon/switch_mon.go
	export GOPATH=`pwd`
	cd src/switch_mon; \
	godep go build -o ../../release/switch_mon

deps:
	cd src/switch_mon && \
	godep get github.com/amir/raidman && \
	godep get github.com/0x19/goesl && \
	godep get golang.org/x/text/encoding && \
	godep get golang.org/x/sys/unix && \
	godep get golang.org/x/tools/go/buildutil && \
	godep get golang.org/x/crypto/ssh/terminal && \
	godep save ./...
