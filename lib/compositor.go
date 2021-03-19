package lib

import (
	"image"
	"image/draw"
)

// Compositor builds a composite image from constituent tile
// images.
type Compositor struct {
	Bounds     image.Rectangle
	addedAreas []image.Rectangle
	Result     draw.Image
}

func NewCompositor(rect image.Rectangle) Compositor {
	return Compositor{
		rect,
		[]image.Rectangle{},
		image.NewRGBA(rect),
	}
}

// Add a new image.  Adjust its contrast range as necessary to match
// any overlapping image data that has already been composited.
func (comp *Compositor) AddImage(image image.Image, subframeRect image.Rectangle) {
	// First draft: don't worry about contrast matching.  Just create the composite image.
	srcPoint := image.Bounds().Min
	draw.Src.Draw(comp.Result, subframeRect, image, srcPoint)
	comp.addedAreas = append(comp.addedAreas, subframeRect)
}
