//go:build !js

package main

import "os"

// initialState reads MOCKUP_STATE on desktop builds.
func initialState() string { return os.Getenv("MOCKUP_STATE") }

// announceFirstFrame is only needed in the browser.
func announceFirstFrame() {}
