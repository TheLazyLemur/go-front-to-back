package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"

	"wasmapp"
	"wasmapp/server/internal"
	"wasmapp/types"
)

var repo = internal.NewContactRepo()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /main.wasm", serveWasm)
	mux.HandleFunc("GET /wasm_exec.js", serveWasmExec)
	mux.HandleFunc("POST /contact", handleCreateContact)
	mux.HandleFunc("GET /contacts", handleGetContacts)

	log.Fatal(http.ListenAndServe(":8080", mux))
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

func handleCreateContact(w http.ResponseWriter, r *http.Request) {
	var req types.CreateContact

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := repo.CreateContact(req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := types.CreateContactResponse{
		ID:    c.ID,
		Name:  c.Name,
		Email: c.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	contacts := repo.GetContacts()

	cs := make([]*types.Contact, len(contacts))
	for i, c := range contacts {
		cs[i] = &types.Contact{
			ID:    c.ID,
			Name:  c.Name,
			Email: c.Email,
		}
	}

	resp := types.GetContactsReponse{
		Data: cs,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
