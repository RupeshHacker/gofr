package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/file"
)

func main() {
	app := gofr.New()

	app.POST("/upload", UploadHandler)

	app.Run()
}

// Data is the struct that we are trying to bind files to
type Data struct {
	// The Compressed field is of type zip,
	// the tag `upload` signifies the key for the form where the file is uploaded
	// if the tag is not present, the field name would be taken as a key.
	Compressed file.Zip `file:"upload"`

	// The FileHeader determines the generic file format that we can get
	// from the multipart form that gets parsed by the incoming HTTP request
	FileHeader *multipart.FileHeader `file:"a"`
}

func UploadHandler(c *gofr.Context) (interface{}, error) {
	var d Data

	// bind the multipart data into the variable d
	err := c.Bind(&d)
	if err != nil {
		return nil, err
	}

	// Defer removing the temporary folder after the function exits
	defer func() {
		if err := os.RemoveAll("tmp"); err != nil {
			c.Logger("Error removing tmp folder:", err)
		}
	}()

	f, err := d.FileHeader.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// read the file content
	content, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}

	// return the number of compressed files recieved
	return fmt.Sprintf("Zipped files: %d, Length of generic file `a`: %d", len(d.Compressed.Files), len(content)), nil
}
