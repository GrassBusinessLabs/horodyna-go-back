package filesystem

import (
	"bufio"
	"bytes"
	stdimg "image"
	"image/jpeg"
	"log"
	"os"
	"path"

	"golang.org/x/image/draw"
)

type ImageStorageService interface {
	SaveImage(filename string, content []byte) error
	RemoveImage(filename string) error
}

type imageStorageService struct {
	loc string
}

func NewImageStorageService(location string) ImageStorageService {
	return imageStorageService{
		loc: location,
	}
}

func (s imageStorageService) SaveImage(filename string, content []byte) error {
	location := path.Join(s.loc, filename)
	err := writeFileToStorage(location, content)
	if err != nil {
		log.Print(err)
		return err
	}

	/*
		scaledImage, err := downScaleImage(image.Image)
		if err != nil {
			log.Print(err)
			return err
		}

		err = writeFileToStorage(s.loc, image.Link+"_thumbnail", scaledImage)
		if err != nil {
			log.Print(err)
			return err
		}
	*/

	return nil
}

func (s imageStorageService) RemoveImage(filename string) error {
	location := path.Join(s.loc, filename)
	err := os.Remove(location)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// nolint
func downScaleImage(input []byte) ([]byte, error) {
	src, err := jpeg.Decode(bytes.NewReader(input))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if src.Bounds().Max.X <= 500 {
		return input, nil
	} else {
		scaleFactor := src.Bounds().Max.X / 500
		dst := stdimg.NewRGBA(stdimg.Rect(0, 0, src.Bounds().Max.X/scaleFactor, src.Bounds().Max.Y/scaleFactor))
		draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		err = jpeg.Encode(writer, dst, nil)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		result := b.Bytes()
		return result, nil
	}
}

func writeFileToStorage(location string, file []byte) error {
	dirLocation := path.Dir(location)
	err := os.MkdirAll(dirLocation, os.ModePerm)
	if err != nil {
		log.Print(err)
		return err
	}

	err = os.WriteFile(location, file, os.ModePerm)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
