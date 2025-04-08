# Captcha

This module provides a simple captcha image generator.

## Installation

```bash
go get github.com/cyberxnomad/captcha
```

## Usage

### example

```go
package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/cyberxnomad/captcha"
)

func main() {
	c, err := captcha.New(
		// replace with your font file path
		captcha.WithFont("/path/to/font.ttf", 40),
		captcha.WithSize(150, 50),
		captcha.WithLength(4, 6),
		captcha.WithScale(0.8, 1.2),
		captcha.WithRotation(-30, 30),
		captcha.WithDistortion(1.2, 2),
		captcha.WithSpacing(0.7, 0.9),
		captcha.WithNoise(0.05),
		captcha.WithLines(4, 6),
	)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("captcha.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, code := c.Generate()
	fmt.Println(code)

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
```

### output

code: `wqufPK`  

image: ![captcha](/sample.png)

## License

This project is licensed under the terms of the **MIT** license. See [LICENSE](LICENSE) for more details.
