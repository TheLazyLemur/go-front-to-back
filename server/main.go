package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"wasmapp"
	"wasmapp/types"
)

var (
	counter  int
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections (modify as needed for security)
		},
	}
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/main.wasm", serveWasm)
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/wasm_exec.js", serveWasmExec)
	mux.HandleFunc("/message", handleMessage)
	mux.HandleFunc("/ws", handleWebSocket)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	var msg types.MyMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(msg.Message)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(
		`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Go WebAssembly</title>
			<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
			<script src="wasm_exec.js"></script>
			<script defer>
				const go = new Go();
				WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
					go.run(result.instance);
				});
			</script>
		</head>
		<body>
			<div id="root"></div>
		</body>
		</html>
		`,
	))
}

func serveWasmExec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(wasmapp.WasmExec)
}

func serveWasm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	w.Write(wasmapp.Wasm)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s\n", msg)

		err = conn.WriteMessage(messageType, fmt.Appendf(nil, "Server Echo: %s", msg))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
