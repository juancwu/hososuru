package views

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

func ToBuffer(c templ.Component) []byte {
    buffer := &bytes.Buffer{}
    c.Render(context.Background(), buffer)
    return buffer.Bytes()
}
