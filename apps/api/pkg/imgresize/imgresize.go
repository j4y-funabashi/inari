package imgresize

import (
	"fmt"
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

func NewResizer(baseDir string) app.Resizer {
	err := os.MkdirAll(baseDir, 0700)
	if err != nil {
		panic("failed to create thumbnails dir: " + err.Error())
	}

	return func(inPath, outPath string) (app.MediaSrc, error) {

		thumbnails := app.MediaSrc{}
		src, err := imaging.Open(inPath, imaging.AutoOrientation(true))
		if err != nil {
			return app.MediaSrc{}, err
		}

		// figure out landscape or portrait
		orientation := orientation(src.Bounds().Dx(), src.Bounds().Dy())

		// -- create lg image
		if orientation == Landscape {
			src = imaging.Resize(src, ImgSizeLG, 0, imaging.Lanczos)
		} else {
			src = imaging.Resize(src, 0, ImgSizeLG, imaging.Lanczos)
		}
		err = imaging.Save(src, filepath.Join(baseDir, generateFilename(ImgSizeLGPrefix, outPath)))
		if err != nil {
			return app.MediaSrc{}, err
		}
		thumbnails.Large = generateFilename(ImgSizeLGPrefix, outPath)

		// -- create sqmd image
		if orientation == Landscape {
			src = imaging.Resize(src, 0, ImgSizeSQMD, imaging.Lanczos)
		} else {
			src = imaging.Resize(src, ImgSizeSQMD, 0, imaging.Lanczos)
		}
		src = imaging.CropAnchor(src, ImgSizeSQMD, ImgSizeSQMD, imaging.Center)
		err = imaging.Save(src, filepath.Join(baseDir, generateFilename("sqmd", outPath)))
		if err != nil {
			return app.MediaSrc{}, err
		}
		thumbnails.Medium = generateFilename("sqmd", outPath)

		// -- create sqsm image
		src = imaging.Resize(src, ImgSizeSQSM, 0, imaging.Lanczos)
		src = imaging.CropAnchor(src, ImgSizeSQSM, ImgSizeSQSM, imaging.Center)
		err = imaging.Save(src, filepath.Join(baseDir, generateFilename("sqsm", outPath)))
		if err != nil {
			return app.MediaSrc{}, err
		}
		thumbnails.Small = generateFilename("sqsm", outPath)

		return thumbnails, nil
	}
}

func generateFilename(prefix, originalImgFilename string) string {
	return fmt.Sprintf("%s_%s", prefix, filepath.Base(originalImgFilename))
}

func orientation(w, h int) string {
	if w > h {
		return "l"
	}
	return "p"
}
