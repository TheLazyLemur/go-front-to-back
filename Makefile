build-wasm:
	templ generate
	cp /opt/homebrew/Cellar/go/1.24.0/libexec/lib/wasm/wasm_exec.js ./static
	GOOS=js GOARCH=wasm go build -o ./static/main.wasm ./client/...

run-server: build-wasm
	go run ./server
