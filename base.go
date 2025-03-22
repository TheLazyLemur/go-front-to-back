package wasmapp

import (
	_ "embed"
)

//go:embed static/main.wasm
var Wasm []byte

//go:embed static/wasm_exec.js
var WasmExec []byte
