package main

import (
	"flag"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	overlayWidth := flag.Int("overlayWidth", 80, "width of image overlay")
	xOffset := flag.Int("xOffset", 16, "x offset of image overlay")
	yOffset := flag.Int("yOffset", 20, "y offset of image overlay")
	parrotPath := flag.String("parrotPath", "parrot.gif", "location of parrot file (.gif)")

	flag.Parse()

	parrot, err := loadParrotFile(*parrotPath)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(*addr, makeParrotHandler(*overlayWidth, *xOffset, *yOffset, parrot))
}

func loadParrotFile(path string) (*gif.GIF, error) {
	parrotFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer parrotFile.Close()

	return gif.DecodeAll(parrotFile)
}

func makeParrotHandler(overlayWidth, xOffset, yOffset int, parrot *gif.GIF) http.Handler {
	var coordinates = [][]int{
		[]int{64 - xOffset, 64 - yOffset},
		[]int{40 - xOffset, 50 - yOffset},
		[]int{26 - xOffset, 53 - yOffset},
		[]int{17 - xOffset, 58 - yOffset},
		[]int{12 - xOffset, 59 - yOffset},
		[]int{18 - xOffset, 65 - yOffset},
		[]int{35 - xOffset, 67 - yOffset},
		[]int{47 - xOffset, 70 - yOffset},
		[]int{56 - xOffset, 66 - yOffset},
		[]int{64 - xOffset, 62 - yOffset},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageFile, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var input image.Image
		var decodeErr error
		switch filepath.Ext(header.Filename) {
		case ".jpeg":
			fallthrough
		case ".jpg":
			input, decodeErr = jpeg.Decode(imageFile)
		case ".png":
			input, decodeErr = png.Decode(imageFile)
		default:
			http.Error(w, "invalid file extension (only .jpg, .jpeg, .png supported)", http.StatusBadRequest)
			return
		}

		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}

		ratio := float64(input.Bounds().Dx()) / float64(overlayWidth)
		input = resize.Resize(uint(overlayWidth), uint(float64(input.Bounds().Dy())/ratio), input, resize.Lanczos3)

		out := &gif.GIF{
			BackgroundIndex: parrot.BackgroundIndex,
			Disposal:        parrot.Disposal,
			Config:          parrot.Config,
			Delay:           parrot.Delay,
			LoopCount:       parrot.LoopCount,
		}

		for i, originalFrame := range parrot.Image {
			frame := image.NewPaletted(originalFrame.Bounds(), originalFrame.Palette)

			draw.Draw(frame, frame.Bounds(), originalFrame, image.ZP, draw.Src)
			draw.Draw(frame, frame.Bounds(), input, image.Point{X: -coordinates[i][0], Y: -coordinates[i][1]}, draw.Over)

			out.Image = append(out.Image, frame)
		}

		if err := gif.EncodeAll(w, out); err != nil {
			http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
			return
		}
	})
}
