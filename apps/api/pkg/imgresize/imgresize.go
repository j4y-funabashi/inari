package imgresize

import (
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

const ImgSizeSQSMPrefix = "sqsm"
const ImgSizeSQMDPrefix = "sqmd"
const ImgSizeLGPrefix = "lg"
const ImgSizeSQSM = 92
const ImgSizeSQMD = 420
const ImgSizeLG = 1080
const Landscape = "l"

func NewResizer() app.Resizer {
	return func(originalImgFilename string) (app.MediaSrc, error) {
		resizedFiles := app.MediaSrc{}

		src, err := imaging.Open(originalImgFilename, imaging.AutoOrientation(true))
		if err != nil {
			return resizedFiles, err
		}

		// figure out landscape or portrait
		orientation := orientation(src.Bounds().Dx(), src.Bounds().Dy())

		// -- create lg image
		if orientation == Landscape {
			src = imaging.Resize(src, ImgSizeLG, 0, imaging.Lanczos)
		} else {
			src = imaging.Resize(src, 0, ImgSizeLG, imaging.Lanczos)
		}
		err = imaging.Save(src, generateFilename(ImgSizeLGPrefix, originalImgFilename))
		if err != nil {
			return resizedFiles, err
		}
		resizedFiles.Large = generateFilename(ImgSizeLGPrefix, originalImgFilename)

		// -- create sqmd image
		if orientation == Landscape {
			src = imaging.Resize(src, 0, ImgSizeSQMD, imaging.Lanczos)
		} else {
			src = imaging.Resize(src, ImgSizeSQMD, 0, imaging.Lanczos)
		}
		src = imaging.CropAnchor(src, ImgSizeSQMD, ImgSizeSQMD, imaging.Center)
		err = imaging.Save(src, generateFilename("sqmd", originalImgFilename))
		if err != nil {
			return resizedFiles, err
		}
		resizedFiles.Medium = generateFilename("sqmd", originalImgFilename)

		// -- create sqsm image
		src = imaging.Resize(src, ImgSizeSQSM, 0, imaging.Lanczos)
		src = imaging.CropAnchor(src, ImgSizeSQSM, ImgSizeSQSM, imaging.Center)
		err = imaging.Save(src, generateFilename("sqsm", originalImgFilename))
		if err != nil {
			return resizedFiles, err
		}
		resizedFiles.Small = generateFilename("sqsm", originalImgFilename)

		return resizedFiles, nil
	}
}

func generateFilename(prefix, originalImgFilename string) string {
	return filepath.Join(os.TempDir(), prefix+"_"+filepath.Base(originalImgFilename))
}

func orientation(w, h int) string {
	if w > h {
		return "l"
	}
	return "p"
}
