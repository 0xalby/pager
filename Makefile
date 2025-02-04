VERSION=0.1

build:
	go build -o bin/pager -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .

run: build
	./bin/pager

clean:
	rm bin/*

release:
	@env GOOS="windows" GOARCH="amd64" go build -o bin/pager_windows_amd64.exe -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="windows" GOARCH="386" go build -o bin/pager_windows_x86.exe -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="darwin" GOARCH="amd64" go build -o bin/pager_macos_amd64 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="darwin" GOARCH="arm64" go build -o bin/pager_macos_arm64 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="linux" GOARCH="amd64" go build -o bin/pager_linux_amd64 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="linux" GOARCH="386" go build -o bin/pager_linux_x86 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="linux" GOARCH="arm64" go build -o bin/pager_linux_arm64 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="linux" GOARCH="arm" go build -o bin/pager_linux_arm32 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="freebsd" GOARCH="amd64" go build -o bin/pager_freebsd_amd64 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .
	@env GOOS="openbsd" GOARCH="amd64" go build -o bin/pager_openbsd_amd64 -ldflags="-s -w -X main.Release=$(VERSION) -extldflags='-static'" .