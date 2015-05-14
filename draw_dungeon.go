package dungeon

import (
	"bufio"
	"fmt"
	"os"
	"image"
	"image/png"
	"code.google.com/p/draw2d/draw2d"
	"math"
	"bitbucket.org/ariqchowdhury/bowyer_watson"
)	

func DrawToPngFile(filePath string, m image.Image) {
	f, err := os.Create(filePath)
	if err != nil {
		os.Exit(1)
	}

	defer f.Close()

	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("Wrote %s OK. \n", filePath)
}

func initGc(w, h int) (image.Image, draw2d.GraphicContext) {
	i := image.NewRGBA(image.Rect(0, 0, w, h))
	gc := draw2d.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)

	return i, gc
}

func DrawCircle(gc draw2d.GraphicContext, x, y, radius float64) {
	xc, yc := x, y

	radiusX := radius
	radiusY := radius

	startAngle := 0.0
	endAngle := 2.0 * math.Pi

	gc.SetLineWidth(1)
	gc.SetLineCap(draw2d.ButtCap)
	gc.SetStrokeColor(image.Black)
	gc.ArcTo(xc, yc, radiusX, radiusY, startAngle, endAngle)
	gc.Stroke()

	gc.Fill()
}

func DrawLine(gc draw2d.GraphicContext, a, b bowyer_watson.Point) {	
	gc.SetStrokeColor(image.Black)
	gc.SetLineWidth(3.0)
	gc.MoveTo(a.X, a.Y)
	gc.LineTo(b.X, b.Y)
	gc.Stroke()
}

func DrawTriangle(gc draw2d.GraphicContext, t bowyer_watson.Triangle) {
	DrawLine(gc, t.A, t.B)
	DrawLine(gc, t.B, t.C)
	DrawLine(gc, t.C, t.A)
}

func (d *Dungeon) DrawDungeon(fileName string) {
	i, gc := initGc(1000, 1000)

	for itr := d.cells.Front(); itr != nil; itr = itr.Next() {
		DrawCircle(gc, itr.Value.(*Cell).x, itr.Value.(*Cell).y, itr.Value.(*Cell).radius)
	}

	DrawToPngFile(fileName, i)
}