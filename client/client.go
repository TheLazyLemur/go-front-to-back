//go:build js && wasm
// +build js,wasm

package main

import (
	"context"
	"fmt"
	"strings"
	"syscall/js"
)

var counter int

func inc(this js.Value, p []js.Value) any {
	counter++
	countEle := js.Global().Get("document").Call("getElementById", "count")
	countEle.Set("innerHTML", fmt.Sprintf("%d", counter))
	js.Global().Call("send", counter)
	return nil
}

func dec(this js.Value, p []js.Value) any {
	counter--
	countEle := js.Global().Get("document").Call("getElementById", "count")
	countEle.Set("innerHTML", fmt.Sprintf("%d", counter))
	js.Global().Call("send", counter)
	return nil
}

func onMessage(this js.Value, args []js.Value) any {
	msg := args[0].Get("data").String()
	fmt.Println("Received from server:", msg)
	return nil
}

func navigateToIndex(this js.Value, args []js.Value) any {
	counter = 0
	var counterHtml strings.Builder
	if err := Counter().Render(context.TODO(), &counterHtml); err != nil {
		panic(err)
	}

	js.Global().
		Get("document").
		Call("getElementById", "root").
		Set("innerHTML", counterHtml.String())
	return nil
}

func navigateToAbout(this js.Value, args []js.Value) any {
	var aboutHTML strings.Builder
	if err := About().Render(context.TODO(), &aboutHTML); err != nil {
		panic(err)
	}

	js.Global().
		Get("document").
		Call("getElementById", "root").
		Set("innerHTML", aboutHTML.String())

	return nil
}

func main() {
	js.Global().Set("inc", js.FuncOf(inc))
	js.Global().Set("dec", js.FuncOf(dec))
	js.Global().Set("navigateToIndex", js.FuncOf(navigateToIndex))
	js.Global().Set("navigateToAbout", js.FuncOf(navigateToAbout))

	ws := js.Global().
		Get("WebSocket").
		New("ws://localhost:8080/ws")

	ws.Set("onmessage", js.FuncOf(onMessage))

	js.Global().Set("send", js.FuncOf(func(this js.Value, args []js.Value) any {
		num := args[0].Int()
		ws.Call("send", num)
		return nil
	}))

	navigateToIndex(js.ValueOf(nil), []js.Value{})

	js.Global().
		Get("document").
		Call("getElementById", "count").
		Set("innerHTML", fmt.Sprintf("%d", counter))

	select {}
}
