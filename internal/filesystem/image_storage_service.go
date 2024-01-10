package filesystem

import (
	"bufio"
	"bytes"
	stdimg "image"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

type ImageStorageService interface {
	SaveImage(filename string, content []byte) (string, error)
	RemoveImage(filename string) error
	UpdateImage(oldfilename string, filename string, content []byte) (string, error)
}

type imageStorageService struct {
	loc string
}

func NewImageStorageService(location string) ImageStorageService {
	return imageStorageService{
		loc: location,
	}
}

func (s imageStorageService) SaveImage(filename string, content []byte) (string, error) {
	name, err := GenerateFileName(s.loc, filename)
	if err != nil {
		return "", err
	}

	location := path.Join(s.loc, name)
	err = writeFileToStorage(location, content)
	if err != nil {
		log.Print(err)
		return "", err
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

	return name, nil
}

func (s imageStorageService) UpdateImage(oldfilename string, filename string, content []byte) (string, error) {
	location := path.Join(s.loc, oldfilename)
	filer, err := os.Open(location)
	if err == nil {
		file_content, err := io.ReadAll(filer)
		if err != nil {
			return "", err
		}
		filer.Close()
		if AreBytesEqual(content, file_content) {
			return oldfilename, nil
		}
		err = s.RemoveImage(oldfilename)
		if err != nil {
			return "", err
		}
	}

	name, err := GenerateFileName(s.loc, filename)
	if err != nil {
		return "", err
	}

	location = path.Join(s.loc, name)
	err = writeFileToStorage(location, content)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return name, nil
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

func GenerateFileName(loc string, name string) (string, error) {
	location := path.Join(loc, name)
	_, err := os.Stat(location)
	if !os.IsNotExist(err) {
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		num := strconv.FormatUint(rng.Uint64(), 10)
		splited := strings.Split(name, ".")
		return GenerateFileName(loc, splited[0]+"_"+num+"."+splited[1])
	}

	return name, nil
}

func AreBytesEqual(b1, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}
	for i := range b1 {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
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
