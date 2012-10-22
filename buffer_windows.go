package buffer

import (
	"github.com/errnoh/utils/bgra"
	"github.com/skelterjohn/go.wde"
	"github.com/skelterjohn/go.wde/win"
)

func setScreen(screen wde.Image) {
	b = &buffer{buffer: make([](*[]uint8), 3)}
	t := screen.(*win.DIB)
	b.draw = bgra.Draw
	back := bgra.New(t.Bounds())
	b.back = back
	b.buffer[0] = &t.Pix
	b.buffer[1] = &back.Pix
}
