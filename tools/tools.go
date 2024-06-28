// Package tools provides some helper functions to assist with
// creating Stream Deck plugins
package tools

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
)

// Turns an image.Image into a string suitable for delivering
// via an events.ESSetImage struct
func ImageToPayload(i image.Image) string {

	out := bytes.Buffer{}
	b64 := base64.NewEncoder(base64.RawStdEncoding, &out)
	err := png.Encode(b64, i)
	if err != nil {
		panic(err)
	}
	return "data:image/png;base64," + out.String()
}

// SVGToPayload create the string necessary to send an SVG
// via a ESSetImage struct
func SVGToPayload(svg string) string {
	return "data:image/svg+xml;charset=utf8," + svg
}
