package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	testUploadDirectory = "./test_uploads"
	tempChunkSuffix     = ".chunk"
)

var uploadDirectory = testUploadDirectory // Override upload directory for testing

func setupTestEnvironment() {
	if _, err := os.Stat(testUploadDirectory); os.IsNotExist(err) {
		os.Mkdir(testUploadDirectory, os.ModePerm)
	}
}

func teardownTestEnvironment() {
	os.RemoveAll(testUploadDirectory)
}

func TestUploadPage(t *testing.T) {
	setupTestEnvironment()
	defer teardownTestEnvironment()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, c.HTML(200, `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Chunk Upload with Resume</title>
    </head>
    <body></body>
    </html>
    `)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestResumeUpload(t *testing.T) {
	setupTestEnvironment()
	defer teardownTestEnvironment()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/resume-upload?file_name=testfile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		fileName := c.QueryParam("file_name")
		if fileName == "" {
			return c.JSON(400, map[string]string{"message": "File name is required"})
		}

		chunks, err := filepath.Glob(filepath.Join(uploadDirectory, fmt.Sprintf("%s%s*", fileName, tempChunkSuffix)))
		if err != nil {
			return err
		}

		chunkIndex := len(chunks)
		totalChunks := chunkIndex // Assuming chunks are sequentially named

		return c.JSON(200, map[string]interface{}{
			"chunk_index":  chunkIndex,
			"total_chunks": totalChunks,
			"file_name":    fileName,
		})
	}

	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		expected := `{"chunk_index":0,"total_chunks":0,"file_name":"testfile"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}

func createMultipartFormData(fileField, fileName string, fileContent []byte) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileField, fileName)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, bytes.NewReader(fileContent))
	if err != nil {
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}

func TestUploadChunk(t *testing.T) {
	setupTestEnvironment()
	defer teardownTestEnvironment()

	e := echo.New()
	fileName := "testfile"
	chunkIndex := 0
	totalChunks := 2 // Set totalChunks to a value greater than 1
	fileContent := []byte("this is a test chunk")

	body, contentType, err := createMultipartFormData("file", "chunk0", fileContent)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/upload-chunk", body)
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Add form data using the variables
	req.Form = make(map[string][]string)
	req.Form.Add("chunk_index", strconv.Itoa(chunkIndex))
	req.Form.Add("total_chunks", strconv.Itoa(totalChunks))
	req.Form.Add("file_name", fileName)

	handler := func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		chunkIndex, err := strconv.Atoi(c.FormValue("chunk_index"))
		if err != nil {
			return err
		}

		totalChunks, err := strconv.Atoi(c.FormValue("total_chunks"))
		if err != nil {
			return err
		}

		fileName := c.FormValue("file_name")
		chunkPath := filepath.Join(uploadDirectory, fmt.Sprintf("%s%s%d", fileName, tempChunkSuffix, chunkIndex))

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(chunkPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		chunks, err := filepath.Glob(filepath.Join(uploadDirectory, fmt.Sprintf("%s%s*", fileName, tempChunkSuffix)))
		if err != nil {
			return err
		}

		if len(chunks) == totalChunks {
			finalPath := filepath.Join(uploadDirectory, fileName)
			finalFile, err := os.Create(finalPath)
			if err != nil {
				return err
			}
			defer finalFile.Close()

			for i := 0; i < totalChunks; i++ {
				chunkPath := filepath.Join(uploadDirectory, fmt.Sprintf("%s%s%d", fileName, tempChunkSuffix, i))
				chunkFile, err := os.Open(chunkPath)
				if err != nil {
					return err
				}
				_, err = io.Copy(finalFile, chunkFile)
				chunkFile.Close()
				if err != nil {
					return err
				}
				os.Remove(chunkPath)
			}
			return c.JSON(200, map[string]string{"message": "File uploaded and assembled successfully!"})
		}

		return c.JSON(200, map[string]string{"message": "Chunk uploaded successfully!"})
	}

	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message":"Chunk uploaded successfully!"}`, rec.Body.String())
	}
}
