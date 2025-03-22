//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"syscall/js"

	"wasmapp/types"
)

func createContact(this js.Value, p []js.Value) any {
	email := js.Global().Get("document").Call("getElementById", "email").Get("value").String()
	name := js.Global().Get("document").Call("getElementById", "name").Get("value").String()

	reqPL := types.CreateContact{
		Name:  name,
		Email: email,
	}

	pl, err := json.Marshal(reqPL)
	if err != nil {
		panic(err)
	}

	go func() {
		defer func() {
			js.Global().Get("document").Call("getElementById", "name").Set("value", "")
			js.Global().Get("document").Call("getElementById", "email").Set("value", "")

			if r := recover(); r != nil {
				js.Global().Get("console").Call("error", "Recovered from panic:", r)
			}
		}()

		res, err := http.Post(
			"http://localhost:8080/contact",
			"application/json",
			bytes.NewBuffer(pl),
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		var c types.CreateContactResponse
		if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
			panic(err)
		}

		list := js.Global().
			Get("document").
			Call("getElementById", "contact-list")

		newDiv := js.Global().Get("document").Call("createElement", "div")
		newDiv.Set("innerHTML", fmt.Sprintf("<div>%s %s %s</div>", c.ID, c.Name, c.Email))
		list.Call("appendChild", newDiv)
	}()

	return nil
}

func navigateToIndex(this js.Value, args []js.Value) any {
	var contactListHTML strings.Builder
	if err := ContactBook().Render(context.TODO(), &contactListHTML); err != nil {
		panic(err)
	}

	js.Global().
		Get("document").
		Call("getElementById", "root").
		Set("innerHTML", contactListHTML.String())

	go func() {
		resp, err := http.Get("http://localhost:8080/contacts")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var respPL types.GetContactsReponse
		if err := json.NewDecoder(resp.Body).Decode(&respPL); err != nil {
			panic(err)
		}

		list := js.Global().
			Get("document").
			Call("getElementById", "contact-list")

		for _, c := range respPL.Data {
			newDiv := js.Global().Get("document").Call("createElement", "div")
			newDiv.Set("innerHTML", fmt.Sprintf("<div>%s %s %s</div>", c.ID, c.Name, c.Email))
			list.Call("appendChild", newDiv)
		}
	}()

	return nil
}

func main() {
	js.Global().Set("navigateToIndex", js.FuncOf(navigateToIndex))
	js.Global().Set("createContact", js.FuncOf(createContact))
	navigateToIndex(js.ValueOf(nil), []js.Value{})

	select {}
}
