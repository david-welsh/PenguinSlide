package fonts

import (
	"PenguinSlide/assets"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

var (
	fontFaceSource *text.GoTextFaceSource
)

func fontInit() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(assets.MenuFontTtf))
	if err != nil {
		log.Fatal(err)
	}

	fontFaceSource = s
}

func FontFaceSource() *text.GoTextFaceSource {
	if fontFaceSource == nil {
		fontInit()
	}
	return fontFaceSource
}
