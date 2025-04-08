package captcha

import (
	"errors"
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"os"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/f64"
	"golang.org/x/image/math/fixed"
)

type CharSet = string

type Captcha struct {
	// width of captcha image
	width int

	// height of captcha image
	height int

	// set of characters to use in captcha
	charSet CharSet

	// minimum length of captcha
	minLength int

	// maximum length of captcha
	maxLength int

	// path to font file, ttf or otf
	fontPath string

	// size of font
	fontSize float64

	// foreground color
	foreground color.Color

	// background color
	background color.Color

	// minimum spacing between characters
	minSpacing float64

	// maximum spacing between characters
	maxSpacing float64

	// minimum rotation of each character
	minRotation float64

	// maximum rotation of each character
	maxRotation float64

	// minimum scaling of each character
	minScale float64

	// maximum scaling of each character
	maxScale float64

	// minimum distortion of each character
	minDistortion float64

	// maximum distortion of each character
	maxDistortion float64

	// minimum number of lines to draw
	minLines int

	// maximum number of lines to draw
	maxLines int

	// level of noise to add to image
	noiseLevel float64

	// font face
	fontFace font.Face
}

func New(opts ...Option) (*Captcha, error) {
	c := &Captcha{
		width:         120,
		height:        50,
		charSet:       AlphaNumericWithoutConfusion,
		minLength:     4,
		maxLength:     4,
		fontPath:      "",
		fontSize:      36,
		foreground:    color.RGBA{0, 0, 0, 255},
		background:    color.RGBA{255, 255, 255, 255},
		minSpacing:    1.0,
		maxSpacing:    1.0,
		minRotation:   0.0,
		maxRotation:   0.0,
		minScale:      1.0,
		maxScale:      1.0,
		minDistortion: 0.0,
		maxDistortion: 0.0,
		minLines:      3,
		maxLines:      7,
		noiseLevel:    0.1,
	}

	for _, opt := range opts {
		opt(c)
	}

	var err error
	switch {
	case c.width <= 0 || c.height <= 0:
		err = errors.New("width and height must be greater than 0")

	case c.charSet == "":
		err = errors.New("char set is required")

	case c.fontPath == "":
		err = errors.New("font path is required")

	case c.fontSize <= 0:
		err = errors.New("font size must be greater than 0")

	case c.minLength < 0 || c.maxLength < 0 || c.minLength > c.maxLength:
		err = errors.New("min length must be greater than 0 and max length must be greater than min length")

	case c.minSpacing < 0 || c.maxSpacing < 0 || c.minSpacing > c.maxSpacing:
		err = errors.New("min spacing must be greater than 0 and max spacing must be greater than min spacing")

	case c.minRotation < -180 || c.maxRotation > 180 || c.minRotation > c.maxRotation:
		err = errors.New("min rotation must be between -180 and 180 and max rotation must be greater than min rotation")

	case c.minScale < 0 || c.maxScale < 0 || c.minScale > c.maxScale:
		err = errors.New("min scale must be greater than 0 and max scale must be greater than min scale")

	case c.minDistortion < 0 || c.maxDistortion < 0 || c.minDistortion > c.maxDistortion:
		err = errors.New("min distortionmust be greater than 0 and max distortion must be greater than min distortion")

	case c.minLines < 0 || c.maxLines < 0 || c.minLines > c.maxLines:
		err = errors.New("min lines must be greater than 0 and max lines must be greater than min lines")

	case c.noiseLevel < 0 || c.noiseLevel > 1:
		err = errors.New("noise level must be between 0 and 1")
	}

	if err != nil {
		return nil, err
	}

	// load font
	fontBytes, err := os.ReadFile(c.fontPath)
	if err != nil {
		return nil, err
	}

	// parse font
	_font, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	// create font face
	face, err := opentype.NewFace(
		_font,
		&opentype.FaceOptions{
			Size:    c.fontSize,
			DPI:     72,
			Hinting: font.HintingFull,
		},
	)
	if err != nil {
		return nil, err
	}

	c.fontFace = face

	return c, nil
}

// Generate generates a captcha image and returns the image and the code
func (c *Captcha) Generate() (image.Image, string) {
	code := c.randomString()

	// create canvas
	canvas := image.NewRGBA(image.Rect(0, 0, c.width, c.height))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(c.background), image.Point{}, draw.Over)

	// draw string
	stringImg := c.drawString(code)

	// caculate start position. make sure the string is centered
	sx, sy := (c.width-stringImg.Bounds().Dx())/2, (c.height-stringImg.Bounds().Dy())/2
	startRect := image.Rect(sx, sy, canvas.Bounds().Dx(), canvas.Bounds().Dy())

	// copy string image to canvas
	draw.Draw(canvas, startRect, stringImg, image.Pt(0, 0), draw.Over)

	canvas = c.drawNoise(canvas)
	canvas = c.drawLines(canvas)

	return canvas, code
}

// randomString returns a random string of length between minLength and maxLength
func (c *Captcha) randomString() string {
	var length int

	if c.minLength == c.maxLength {
		length = c.minLength
	} else {
		length = rand.IntN(c.maxLength-c.minLength+1) + c.minLength
	}

	return randomString(length, c.charSet)
}

// drawString draws the string on the canvas
func (c *Captcha) drawString(code string) *image.RGBA {
	// calculate width and height.
	width := font.MeasureString(c.fontFace, code).Ceil() * int(c.maxSpacing*c.maxScale) * 2
	height := c.fontFace.Metrics().Height.Ceil() * int(c.maxScale) * 2

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	x, y := 0, height/3

	for _, char := range code {
		charImg := c.drawChar(char)

		charBounds := charImg.Bounds()

		// random y offset
		yOfs := rand.IntN(charBounds.Dy()/4) - charBounds.Dy()/8

		// copy char image to canvas
		startRect := image.Rect(x, y+yOfs, canvas.Bounds().Dx(), canvas.Bounds().Dy())
		draw.Draw(canvas, startRect, charImg, image.Pt(0, 0), draw.Over)

		spacing := 1.0
		if c.minSpacing == c.maxSpacing {
			spacing = c.minSpacing
		} else {
			spacing = rand.Float64()*(c.maxSpacing-c.minSpacing) + c.minSpacing
		}

		// next character x position
		x += int(float64(charBounds.Dx()) * spacing)
	}

	// cut canvas to remove empty space
	bounds := pixelBounds(canvas)

	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), canvas, bounds.Min, draw.Over)

	return dst
}

// drawChar draws the character on the canvas
func (c *Captcha) drawChar(char rune) *image.RGBA {
	// calculate char width and height
	charWith := font.MeasureString(c.fontFace, string(char)).Ceil()
	metrics := c.fontFace.Metrics()
	ascent := metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	charHeight := ascent + descent

	// create canvas
	size := int(math.Max(float64(charWith), float64(charHeight)))
	canvas := image.NewRGBA(image.Rect(0, 0, size, size))

	// draw char
	drawer := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(c.foreground),
		Face: c.fontFace,
		Dot: fixed.Point26_6{
			X: fixed.I(size/2 - charWith/2),
			Y: fixed.I(size/2 + descent),
		},
	}
	drawer.DrawString(string(char))

	// apply effectors
	canvas = c.scaleChar(canvas)
	canvas = c.distortChar(canvas)
	canvas = c.rotateChar(canvas)

	// cut canvas to remove empty space
	bounds := pixelBounds(canvas)

	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), canvas, bounds.Min, draw.Over)

	return dst
}

func (c *Captcha) rotateChar(src *image.RGBA) *image.RGBA {
	var rotation float64

	if c.minRotation == c.maxRotation {
		// no rotation
		if c.minRotation == 0 {
			return src
		}

		rotation = c.minRotation
	} else {
		rotation = rand.Float64()*(c.maxRotation-c.minRotation) + c.minRotation
	}

	srcBounds := src.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()

	// calculate dst size
	sin, cos := math.Sincos(math.Pi * rotation / 180)
	dstW := int(math.Abs(float64(srcW)*cos + math.Abs(float64(srcW)*sin)))
	dstH := int(math.Abs(float64(srcH)*cos + math.Abs(float64(srcH)*sin)))

	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	// calculate center
	srcCx, srcCy := float64(srcW)/2, float64(dstH)/2
	dstCx, dstCy := float64(dstW)/2, float64(dstH)/2

	// make affine transformation matrix
	m := f64.Aff3{
		cos, -sin, dstCx - srcCx*cos + srcCy*sin,
		sin, cos, dstCy - srcCx*sin - srcCy*cos,
	}

	// apply transformation with Catmull-Rom
	draw.CatmullRom.Transform(dst, m, src, srcBounds, draw.Over, nil)

	return dst
}

// scaleChar scales the character
func (c *Captcha) scaleChar(src *image.RGBA) *image.RGBA {
	var scaleX, scaleY float64

	if c.minScale == c.maxScale {
		// no scale
		if c.minScale == 1 {
			return src
		}

		scaleX = c.minScale
		scaleY = c.minScale
	} else {
		scaleX = rand.Float64()*(c.maxScale-c.minScale) + c.minScale
		scaleY = rand.Float64()*(c.maxScale-c.minScale) + c.minScale
	}

	// calculate new size
	newW := int(float64(src.Bounds().Dx()) * scaleX)
	newH := int(float64(src.Bounds().Dy()) * scaleY)
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// scale the image with Catmull-Rom
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return dst
}

// distortChar distorts the character
func (c *Captcha) distortChar(src *image.RGBA) *image.RGBA {
	var amplitude float64

	if c.minDistortion == c.maxDistortion {
		// no distortion
		if c.minDistortion == 0 {
			return src
		}

		amplitude = c.minDistortion
	} else {
		amplitude = rand.Float64()*(c.maxDistortion-c.minDistortion) + c.minDistortion
	}

	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	dst := image.NewRGBA(src.Bounds())
	period := float64(h) / 2

	for y := range h {
		for x := range w {
			// calculate x and y offset
			xOfs := amplitude * math.Sin(2*math.Pi*float64(y)/period)
			yOfs := amplitude * math.Cos(2*math.Pi*float64(x)/period)

			nx := x + int(xOfs)
			ny := y + int(yOfs)

			// check if the offset is within the image bounds
			if nx >= 0 && nx < w && ny >= 0 && ny < h {
				dst.Set(x, y, src.At(nx, ny))
			}
		}
	}

	return dst
}

// drawNoise draws noise on the image
func (c *Captcha) drawNoise(src *image.RGBA) *image.RGBA {
	if c.noiseLevel == 0 {
		return src
	}

	if c.noiseLevel > 1 {
		c.noiseLevel = 1
	}

	// get random noise count
	noiseCnt := int(float64(src.Bounds().Dx()*src.Bounds().Dy()) * c.noiseLevel)

	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, image.Pt(0, 0), draw.Src)

	for range noiseCnt {
		x := rand.IntN(src.Bounds().Dx())
		y := rand.IntN(src.Bounds().Dy())
		dst.Set(x, y, c.foreground)
	}

	return dst
}

// drawLines draws lines on the image
func (c *Captcha) drawLines(src *image.RGBA) *image.RGBA {
	var lineCnt int

	if c.minLines == c.maxLines {
		// no lines
		if c.minLines == 0 {
			return src
		}

		lineCnt = c.minLines
	} else {
		lineCnt = rand.IntN(c.maxLines-c.minLines) + c.minLines
	}

	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, image.Pt(0, 0), draw.Src)

	for range lineCnt {
		// draw straight line or curve line randomly
		if rand.Float64() < 0.5 {
			c.drawStraightLine(dst)
		} else {
			c.drawCurveLine(dst)
		}
	}

	return dst
}

// drawStraightLine draws a straight line on the image
func (c *Captcha) drawStraightLine(img *image.RGBA) {
	// get random line start and end point
	x1, y1 := rand.IntN(img.Bounds().Dx()), rand.IntN(img.Bounds().Dy())
	x2, y2 := rand.IntN(img.Bounds().Dx()), rand.IntN(img.Bounds().Dy())

	// get random line color
	lineColor := randomNearColor(c.foreground)

	// Bresenham line algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	sx, sy := -1, -1
	if x1 < x2 {
		sx = 1
	}
	if y1 < y2 {
		sy = 1
	}

	// error term
	err := dx - dy

	for {
		img.Set(x1, y1, lineColor)

		if x1 == x2 && y1 == y2 {
			break
		}

		err2 := 2 * err
		if err2 > -dy {
			err -= dy
			x1 += sx
		}
		if err2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// drawCurveLine draws a curve line on the image
func (c *Captcha) drawCurveLine(img *image.RGBA) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()

	// get random curve line start and end point
	x1, y1 := rand.IntN(w), rand.IntN(h)
	x2, y2 := rand.IntN(w), rand.IntN(h)

	// get random curve line control point
	ctrlX, ctrlY := rand.IntN(w), rand.IntN(h)

	// get random line color
	lineColor := randomNearColor(c.foreground)

	// Bezier curve algorithm
	// change steps to smooth the curve
	steps := 120
	for step := range steps {
		t := float64(step) / float64(steps)
		x := int(lerp(float64(x1), float64(x2), t))
		y := int(lerp(float64(y1), float64(y2), t))
		ctrlXt := int(lerp(float64(x1), float64(ctrlX), t))
		ctrlYt := int(lerp(float64(y1), float64(ctrlY), t))
		finalX := int(lerp(float64(x), float64(ctrlXt), t))
		finalY := int(lerp(float64(y), float64(ctrlYt), t))

		if finalX >= 0 && finalX < w && finalY >= 0 && finalY < h {
			img.Set(finalX, finalY, lineColor)
		}
	}
}
