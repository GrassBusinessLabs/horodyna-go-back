package filesystem

import (
	"bufio"
	"bytes"
	stdimg "image"
	"image/jpeg"
	"io/ioutil"
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
	name, err := FileExists(s.loc, filename, content)
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

func (s imageStorageService) RemoveImage(filename string) error {
	location := path.Join(s.loc, filename)
	err := os.Remove(location)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func FileExists(loc string, name string, file []byte) (string, error) {
	location := path.Join(loc, name)
	file_cont, err := ioutil.ReadFile(location)
	if err != nil {
		return name, nil
	}

	rand.Seed(time.Now().UnixNano())
	if !AreBytesEqual(file_cont, file) {
		num := strconv.FormatUint(rand.Uint64(), 10)
		splited := strings.Split(name, ".")
		return FileExists(loc, splited[0]+"_"+num+"."+splited[1], file)
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
