package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/simple-go-image-processing"
	"github.com/stretchr/testify/assert"
)

func TestConverter(t *testing.T) {
	app := main.SetupApp()

	t.Run("should respond with a JSON error message when the image field is missing", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/converter", nil)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "You have to provide `image` field!", data["message"])
	})

	t.Run("should respond with a JSON error message for wrong image extension; expected PNG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://placehold.co/600x400.jpg"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreateFormFile("image", "image.jpg")
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/converter", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "Please provide an image with PNG extension!", data["message"])
	})

	t.Run("should respond with success and provide a JPEG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://placehold.co/600x400.png"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "image.png")},
			"Content-Type":              []string{"image/png"},
			"Content-Transfer-Encoding": []string{"binary"},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/converter", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Greater(t, resp.ContentLength, int64(0))
		assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
	})
}

func TestResizer(t *testing.T) {
	app := main.SetupApp()

	t.Run("should respond with a JSON error message for missing width", func(t *testing.T) {
		body := bytes.Buffer{}

		req := httptest.NewRequest(http.MethodPost, "/resizer", &body)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "`width` field is missing!", data["message"])
	})

	t.Run("should respond with a JSON error message for missing height", func(t *testing.T) {
		body := bytes.Buffer{}

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=ABC", &body)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "`height` field is missing!", data["message"])
	})

	t.Run("should respond with a JSON validation error message for non-numeric width", func(t *testing.T) {
		body := bytes.Buffer{}

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=ABC&height=DEF", &body)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "Please provide only an integer number for `width`!", data["message"])
	})

	t.Run("should respond with a JSON validation error message for non-numeric height", func(t *testing.T) {
		body := bytes.Buffer{}

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=-600&height=DEF", &body)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "Please provide only an integer number for `height`!", data["message"])
	})

	t.Run("should respond with a JSON error message for invalid dimension values (negative or zero height or width)", func(t *testing.T) {
		body := bytes.Buffer{}

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=-600&height=-400", &body)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "Invalid dimension values for `width` or/and `height` fields!", data["message"])
	})

	t.Run("should respond with a JSON error message when the image field is missing", func(t *testing.T) {
		body := bytes.Buffer{}

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=600&height=400", &body)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "You have to provide `image` field!", data["message"])
	})

	t.Run("should respond with a JSON error message for wrong image extension; expected PNG or JPEG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://placehold.co/600x400.svg"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreateFormFile("image", "image.svg")
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=600&height=400", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "The image must be in JPEG or PNG format!", data["message"])
	})

	t.Run("should respond with success and provide the resized JPEG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://placehold.co/600x400.jpg"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "image.jpg")},
			"Content-Type":              []string{"image/jpeg"},
			"Content-Transfer-Encoding": []string{"binary"},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=1000&height=1000", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Greater(t, resp.ContentLength, int64(0))
		assert.Equal(t, imgResp.Header.Get("Content-Type"), resp.Header.Get("Content-Type"))
	})

	t.Run("should respond with success and provide the resized PNG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://placehold.co/600x400.png"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "image.png")},
			"Content-Type":              []string{"image/png"},
			"Content-Transfer-Encoding": []string{"binary"},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/resizer?width=1000&height=1000", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Greater(t, resp.ContentLength, int64(0))
		assert.Equal(t, imgResp.Header.Get("Content-Type"), resp.Header.Get("Content-Type"))
	})
}

func TestCompressor(t *testing.T) {
	app := main.SetupApp()

	t.Run("should respond with a JSON error message when the image field is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/compressor", nil)
		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "You have to provide `image` field!", data["message"])
	})

	t.Run("should respond with a JSON error message for wrong image extension; expected PNG or JPEG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://placehold.co/600x400.svg"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreateFormFile("image", "image.svg")
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/compressor", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "error", data["status"])
		assert.Equal(t, "The image must be in JPEG or PNG format!", data["message"])
	})

	t.Run("should respond with success and provide the compressed JPEG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://picsum.photos/200/300"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "image.jpg")},
			"Content-Type":              []string{"image/jpeg"},
			"Content-Transfer-Encoding": []string{"binary"},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/compressor", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Greater(t, resp.ContentLength, int64(0))
		assert.Less(t, resp.ContentLength, imgResp.ContentLength)
		assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
	})

	t.Run("should respond with success and provide the compressed PNG image", func(t *testing.T) {
		body := bytes.Buffer{}

		fileUrl := "https://placehold.co/600x400.png"
		imgResp, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer imgResp.Body.Close()

		writer := multipart.NewWriter(&body)

		fw, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       []string{fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "image.png")},
			"Content-Type":              []string{"image/png"},
			"Content-Transfer-Encoding": []string{"binary"},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(fw, imgResp.Body)
		if err != nil {
			log.Fatal(err)
		}
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/compressor", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, _ := app.Test(req)

		data := map[string]any{}

		_ = json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Greater(t, resp.ContentLength, int64(0))
		assert.Less(t, resp.ContentLength, imgResp.ContentLength)
		assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
	})
}
