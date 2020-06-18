package interactor

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/ezio1119/fishapp-image/conf"
	"github.com/ezio1119/fishapp-image/models"
)

func resizeImage(i *models.Image) error {

	img, t, err := image.Decode(i.Buf)
	if err != nil {
		return err
	}

	nrgba := imaging.Fit(img, conf.C.Sv.ImageWidth, conf.C.Sv.ImageHeight, imaging.Lanczos)
	buf := &bytes.Buffer{}

	switch t {
	case "jpeg":
		if err := jpeg.Encode(buf, nrgba, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return err
		}
		i.Name = i.Name + ".jpg"
	case "png":
		if err := png.Encode(buf, nrgba); err != nil {
			return err
		}
		i.Name = i.Name + ".png"
	case "gif":
		if err := gif.Encode(buf, nrgba, nil); err != nil {
			return err
		}
		i.Name = i.Name + ".gif"
	}

	return nil
}
