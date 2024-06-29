# Stream Deck plugin library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/tardisx/streamdeck-plugin.svg)](https://pkg.go.dev/github.com/tardisx/streamdeck-plugin)

You can find fully-formed examples using this library in
[streamdeck-plugin-examples](https://github.com/tardisx/streamdeck-plugin-examples)

## Basic usage

```go
package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/tardisx/streamdeck-plugin"
	"github.com/tardisx/streamdeck-plugin/events"
)

// keep track of instances we've seen
var contexts = map[string]bool{}

func main() {
	slog.Info("Starting up")
	c := streamdeck.New()

	slog.Info("Registering handlers")
	c.RegisterHandler(func(e events.ERWillAppear) {
		slog.Info(fmt.Sprintf("action %s appeared, context %s", e.Action, e.Context))
		contexts[e.Context] = true
	})
	c.RegisterHandler(func(e events.ERWillDisappear) {
		slog.Info(fmt.Sprintf("action %s disappeared, context %s", e.Action, e.Context))
		delete(contexts, e.Context)
	})
	c.RegisterHandler(func(e events.ERKeyDown) {
		slog.Info(fmt.Sprintf("action %s appeared, context %s", e.Action, e.Context))
	})

	slog.Info("Connecting web socket")
	err := c.Connect()
	if err != nil {
		panic(err)
	}

	// update the title once a second, for all "seen" contexts
	go func() {
		for {
			for context := range contexts {
				c.Send(events.NewESSetTitle(
					context,
					time.Now().Format(time.Kitchen),
					events.EventTargetBoth,
					0))
			}
			time.Sleep(time.Second)
		}
	}()

	slog.Info("waiting for the end")
	c.WaitForPluginExit()
}
```