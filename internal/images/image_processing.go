package images

import (
	"github.com/fogleman/gg"
	"image"
	"image/png"
	"net/http"
)

func DownloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func AddBackgroundToImage(src image.Image) image.Image {
	const (
		width  = 725 // ширина нового изображения с фоном
		height = 400 // высота нового изображения с фоном
	)

	dc := gg.NewContext(width, height)
	dc.SetRGB(0, 0, 0) // черный цвет
	dc.Clear()

	x := (width - src.Bounds().Dx()) / 2
	y := (height - src.Bounds().Dy()) / 2

	dc.DrawImage(src, x, y)

	return dc.Image()
}
