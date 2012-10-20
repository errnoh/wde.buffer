package buffer

import (
	"fmt"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/skelterjohn/go.wde"
	"github.com/skelterjohn/go.wde/xgb"
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"
)

func Screen(width, height int) wde.Image {
	r := image.Rect(0, 0, width, height)
	return &xgb.Image{
		&xgraphics.Image{
			X:      nil,
			Pixmap: 0,
			Pix:    make([]uint8, 4*r.Dx()*r.Dy()),
			Stride: 4 * r.Dx(),
			Rect:   r,
			Subimg: false,
		},
	}
}

func debug() {
	fmt.Println(
		"Front:", b.buffer[0], "\n",
		"Middle:", b.buffer[1], "\n",
		"Back:", b.buffer[2], "\n",
	)
}

func printbuffer() {
	for y := 0; y < b.back.Bounds().Dy(); y++ {
		for x := 0; x < b.back.Bounds().Dx(); x++ {
			r, _, _, _ := b.back.At(x, y).RGBA()
			if r>>8 == 0 {
				fmt.Print(0)
			} else {
				fmt.Print(1)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func colorcheck(t *testing.T, current, target color.Color) {
	r1, g1, b1, a1 := current.RGBA()
	r2, g2, b2, a2 := target.RGBA()
	if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
		t.Errorf("Got [%d %d %d %d], expected [%d %d %d %d]", r1>>8, g1>>8, b1>>8, a1>>8, r2>>8, g2>>8, b2>>8, a2>>8)
	}
}

func TestStuff(t *testing.T) {
	screen := Screen(10, 10)
	setScreen(screen)
	SetEmptyColor(image.White)
	colorcheck(t, screen.At(4, 4), color.RGBA{0, 0, 0, 0})
	Draw(image.Rect(3, 3, 3+2, 3+2), image.Black, image.ZP, draw.Src)
	colorcheck(t, screen.At(4, 4), color.RGBA{0, 0, 0, 0})
	Flip()
	colorcheck(t, screen.At(4, 4), color.Black)
	Flip()
	colorcheck(t, screen.At(4, 4), color.White)
}

func TestBench(t *testing.T) {
	fmt.Println("with empty")
	start := time.Now()
	screen := Screen(2000, 2000)
	setScreen(screen)
	SetEmptyColor(image.White)
	var count int
	go func(i *int) {
		start := time.Now()
		for {
			Flip()
			count++
			if time.Since(start) > time.Second {
				return
			}
		}
	}(&count)
	<-time.After(time.Second)
	fmt.Println(count, time.Since(start))
}

func TestBenchWithNoEmpty(t *testing.T) {
	fmt.Println("no empty")
	start := time.Now()
	screen := Screen(2000, 2000)
	setScreen(screen)
	var count int
	go func(i *int) {
		start := time.Now()
		for {
			Flip()
			count++
			if time.Since(start) > time.Second {
				return
			}
		}
	}(&count)
	<-time.After(time.Second)
	fmt.Println(count, time.Since(start))
}

func TestBenchWithWdeImage(t *testing.T) {
	fmt.Println("wde")
	r := image.Rect(0, 0, 2000, 2000)
	source := image.NewRGBA(r)
	draw.Draw(source, source.Bounds(), image.White, image.ZP, draw.Src)
	target := xgb.Image{&xgraphics.Image{
		X:      nil,
		Pixmap: 0,
		Pix:    make([]uint8, 4*r.Dx()*r.Dy()),
		Stride: 4 * r.Dx(),
		Rect:   r,
		Subimg: false,
	}}
	start := time.Now()
	var count int
	go func(i *int) {
		start := time.Now()
		for {
			target.CopyRGBA(source, r)
			count++
			if time.Since(start) > time.Second {
				return
			}
		}
	}(&count)
	<-time.After(time.Second)
	fmt.Println(count, time.Since(start))
}
