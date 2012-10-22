/*
Copyright (c) 2012, Erno Hopearuoho
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
    notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in the
    documentation and/or other materials provided with the distribution.
    * Neither the name of the <organization> nor the
    names of its contributors may be used to endorse or promote products
    derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
wde.buffer:

High speed buffering for go.wde
( https://github.com/skelterjohn/go.wde )

Mimics triple buffering.
There are two buffers the size of current screen.
One to hold the current screen and other to draw onto.
Flip() switches buffers and renders the screen.
*/
package buffer

import (
	"github.com/skelterjohn/go.wde"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

var (
	w          wde.Window
	b          *buffer
	empty      *image.Uniform
	flushmutex = new(sync.Mutex)
)

type buffer struct {
	draw   func(draw.Image, image.Rectangle, image.Image, image.Point, draw.Op)
	buffer [](*[]uint8)
	back   Image
}

type Image interface {
	draw.Image
	PixOffset(x, y int) int
	SetRGBA(x, y int, c color.RGBA)
}

// Create new buffer.
// params:
// window       - wde compatible window
// background   - optional background color
func Create(window wde.Window, background color.Color) {
	w = window
	setScreen(window.Screen())
	if background != nil {
		SetEmptyColor(background)
		b.draw(b.back, b.back.Bounds(), empty, image.ZP, draw.Src)
	}
}

// Switch buffers and render screen with the new one.
func Flip() {
	flushmutex.Lock()
	*b.buffer[0], *b.buffer[1] = *b.buffer[1], *b.buffer[0]
	go func() {
		if w != nil {
			w.FlushImage()
		}
		flushmutex.Unlock()
	}()
	if empty != nil {
		b.draw(b.back, b.back.Bounds(), empty, image.ZP, draw.Src)
	}
}

// Set background color to something else.
//
// NOTE: Change is applied after next Flip()
//       You need to Draw() the color to next buffer
//       if you want instant change.
func SetEmptyColor(c color.Color) {
	empty = image.NewUniform(c)
}

// Similar to image/draw Draw() function but without the first argument.
// Always draws to the next back buffer.
func Draw(r image.Rectangle, src image.Image, sp image.Point, op draw.Op) {
	b.draw(b.back, r, src, sp, op)
}

// Changes the color of single pixel (color.RGBA)
func Set(x, y int, c color.RGBA) {
	b.back.SetRGBA(x, y, c)
}
