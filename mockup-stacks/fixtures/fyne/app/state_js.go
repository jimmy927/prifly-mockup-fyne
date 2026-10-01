//go:build js

package main

import (
	"syscall/js"
	"time"
)

// initialState reads window.__state, which the packed page sets from location.hash.
func initialState() string {
	v := js.Global().Get("__state")
	if v.IsUndefined() || v.IsNull() {
		return ""
	}
	return v.String()
}

// announceFirstFrame tells the packed page that Fyne has drawn: it resolves window.__mockupReady.
func announceFirstFrame() {
	go func() {
		time.Sleep(1500 * time.Millisecond)
		if f := js.Global().Get("__goReady"); f.Type() == js.TypeFunction {
			f.Invoke()
		}
	}()
}
