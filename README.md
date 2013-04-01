wde.buffer
==========

High speed buffering for go.wde


Init like you usually would in wde, just create buffer in the end:

    win, err = wde.NewWindow(w, h)
    win.Show()
    buffer.Create(win, nil)


Render thread example:

    r := win.Screen().Bounds()
    for {
        buffer.Draw(r, image.Black, image.ZP, draw.Src)
        buffer.Flip()
    }

When you're done, close like usual:
    win.Close()
