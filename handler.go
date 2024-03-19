package main

import (
	"image"
	"io"
	"log"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gocv.io/x/gocv"
)

var supportedTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
}

func converter(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "You have to provide `image` field!")
	}

	if file.Header.Get("Content-Type") != "image/png" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide an image with PNG extension!")
	}

	f, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to stream image file")
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read image file")
	}

	img, err := gocv.IMDecode(b, gocv.IMReadColor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to decode image")
	}

	buf, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to convert image to JPG")
	}

	c.Set("Content-Type", "image/jpeg")
	if err = c.Send(buf.GetBytes()); err != nil {
		return err
	}

	log.Println("=== CONVERTER ===")
	log.Println("From PNG to JPG")

	return nil
}

func resizer(c *fiber.Ctx) error {
	widthValue, heightValue := c.FormValue("width", "0"), c.FormValue("height", "0")

	width, err := strconv.Atoi(widthValue)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide only an integer number for `width`!")
	}

	height, err := strconv.Atoi(heightValue)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide only an integer number for `height`!")
	}

	if width <= 0 && height <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "You must specify the dimensions using the `width` and `height` fields!")
	}

	file, err := c.FormFile("image")
	if err != nil {

		return fiber.NewError(fiber.StatusBadRequest, "You have to provide `image` field!")
	}

	extension := filepath.Ext(file.Filename)
	mime := file.Header.Get("Content-Type")

	if mime != "image/jpeg" && mime != "image/png" {
		return fiber.NewError(fiber.StatusBadRequest, "The image must be in JPEG or PNG format!")
	}

	f, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to stream image file")
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read image file")
	}

	img, err := gocv.IMDecode(b, gocv.IMReadColor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to decode image")
	}

	resizedImg := gocv.NewMat()
	defer resizedImg.Close()

	gocv.Resize(img, &resizedImg, image.Pt(width, height), 0, 0, gocv.InterpolationDefault)

	buf, err := gocv.IMEncode(gocv.FileExt(extension), resizedImg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to resize image dimension")
	}

	c.Set("Content-Type", supportedTypes[extension])
	if err = c.Send(buf.GetBytes()); err != nil {
		return err
	}

	log.Println("=== RESIZER ===")
	log.Println("Original\t:", img.Cols(), " x ", img.Rows())
	log.Println("Resized\t:", resizedImg.Cols(), " x ", resizedImg.Rows())

	return nil
}

func compressor(c *fiber.Ctx) error {

	file, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "You have to provide `image` field!")
	}

	extension := filepath.Ext(file.Filename)
	mime := file.Header.Get("Content-Type")

	if mime != "image/jpeg" && mime != "image/png" {
		return fiber.NewError(fiber.StatusBadRequest, "The image must be in JPEG or PNG format!")
	}

	f, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to stream image file")
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read image file")
	}

	img, err := gocv.IMDecode(b, gocv.IMReadColor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to decode image")
	}

	compressionParams := []int{}

	switch mime {
	case "image/jpeg":
		compressionParams = append(compressionParams, gocv.IMWriteJpegQuality, 50)
	case "image/png":
		compressionParams = append(compressionParams, gocv.IMWritePngCompression, 5)
	}

	buf, err := gocv.IMEncodeWithParams(gocv.FileExt(extension), img, compressionParams)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to compress image size")
	}

	c.Set("Content-Type", supportedTypes[extension])
	if err = c.Send(buf.GetBytes()); err != nil {
		return err
	}

	log.Println("=== COMPRESSOR ===")
	log.Println("Before\t:", file.Size)
	log.Println("After\t:", len(buf.GetBytes()))

	return nil
}
