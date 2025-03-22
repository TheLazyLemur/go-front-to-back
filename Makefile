build-wasm:
	templ generate
	GOOS=js GOARCH=wasm go build -o ./static/main.wasm ./client/...

run-server: build-wasm
	go run ./server
