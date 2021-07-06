package spaghetti

import (
	"errors"
	"syscall/js"
)

// PromiseResult is the result of a promise
type PromiseResult struct {
	This   js.Value
	Values []js.Value
	Error  error
}

// ResolvePromise resolves a JS promise
func ResolvePromise(promise js.Value) <-chan PromiseResult {
	channel := make(chan PromiseResult)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer close(channel)
		channel <- PromiseResult{
			Values: args,
			Error:  nil,
		}
		return nil
	}))
	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer close(channel)
		channel <- PromiseResult{
			This:   this,
			Values: args,
			Error:  errors.New(args[0].String()),
		}
		return nil
	}))

	return channel
}
