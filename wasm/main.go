//go:build js && wasm

package main

import "syscall/js"

func hello(this js.Value, args []js.Value) any {
	return js.ValueOf("Hello from Go WASM!")
}

func main() {
	js.Global().Set("goHello", js.FuncOf(hello))
	select {} // keep Go runtime alive
}
