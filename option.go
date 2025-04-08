package captcha

import "image/color"

const (
	Lowercase                    CharSet = "abcdefghijklmnopqrstuvwxyz"
	Uppercase                    CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alphabetic                   CharSet = Lowercase + Uppercase
	Numeric                      CharSet = "0123456789"
	Hex                          CharSet = Numeric + "abcdef"
	LowerNumeric                 CharSet = Lowercase + Numeric
	UpperNumeric                 CharSet = Uppercase + Numeric
	AlphaNumeric                 CharSet = Alphabetic + Numeric
	AlphaNumericWithoutConfusion CharSet = "ABCDEFGHKLMNPQRSTUVWXYZabcdefghkmnpqsuvwxyz23456789"
)

type Option func(*Captcha)

// WithSize sets the size of the captcha image.
//
// Default: 140x50
func WithSize(width, height int) Option {
	return func(c *Captcha) {
		c.width = width
		c.height = height
	}
}

// WithCharSet sets the character set used to generate the captcha.
//
// Default: AlphaNumeric
func WithCharSet(charSet CharSet) Option {
	return func(c *Captcha) {
		c.charSet = charSet
	}
}

// WithLength sets the minimum and maximum length of the captcha.
//
// Default: 4, 4
func WithLength(minLength, maxLength int) Option {
	return func(c *Captcha) {
		c.minLength = minLength
		c.maxLength = maxLength
	}
}

// WithFont sets the font used to generate the captcha.
//
// Default: nil, 40
func WithFont(path string, size float64) Option {
	return func(c *Captcha) {
		c.fontPath = path
		c.fontSize = size
	}
}

// WithBackground sets the background color of the captcha.
//
// Default: color.White
func WithBackground(background color.Color) Option {
	return func(c *Captcha) {
		c.background = background
	}
}

// WithForeground sets the foreground color of the captcha.
//
// Default: color.Black
func WithForeground(foreground color.Color) Option {
	return func(c *Captcha) {
		c.foreground = foreground
	}
}

// WithSpacing sets the spacing between characters.
//
// Default: 1.0, 1.0
func WithSpacing(minSpacing, maxSpacing float64) Option {
	return func(c *Captcha) {
		c.minSpacing = minSpacing
		c.maxSpacing = maxSpacing
	}
}

// WithRotation sets the rotation angle of characters.
//
// Default: 0.0, 0.0
func WithRotation(minRotation, maxRotation float64) Option {
	return func(c *Captcha) {
		c.minRotation = minRotation
		c.maxRotation = maxRotation
	}
}

// WithScale sets the scaling factor of characters.
//
// Default: 1.0, 1.0
func WithScale(minScale, maxScale float64) Option {
	return func(c *Captcha) {
		c.minScale = minScale
		c.maxScale = maxScale
	}
}

// WithDistortion sets the distortion factor of characters.
//
// Default: 0.0, 0.0
func WithDistortion(minDistortion, maxDistortion float64) Option {
	return func(c *Captcha) {
		c.minDistortion = minDistortion
		c.maxDistortion = maxDistortion
	}
}

// WithLines sets the number of lines to draw on the captcha.
//
// Default: 3, 7
func WithLines(minLines, maxLines int) Option {
	return func(c *Captcha) {
		c.minLines = minLines
		c.maxLines = maxLines
	}
}

// WithNoise sets the level of noise to add to the captcha.
//
// Default: 0.1
func WithNoise(level float64) Option {
	return func(c *Captcha) {
		c.noiseLevel = level
	}
}
