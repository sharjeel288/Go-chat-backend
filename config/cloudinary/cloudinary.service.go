package cloudinaryService

import (
	"ChaiLabs/config"
	"ChaiLabs/utils"
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/google/uuid"
)

type BulkFileUploadRequest struct {
	buffer   []byte
	filename string
}

func UploadBufferToCloudinary(ctx context.Context, buffer []byte, filename string, opts ...bool) (string, error) {

	randomizeFilename := true
	if len(opts) > 0 {
		randomizeFilename = opts[0]
	}

	cld, err := cloudinary.NewFromParams(config.CloudinaryCloudName(), config.CloudinaryApiKey(), config.CloudinaryApiSecret())
	if err != nil {
		fmt.Println("Error initializing Cloudinary:", err)
		return "", err
	}

	if filename == "" {
		filename = uuid.New().String()
	}

	// Sanitize filename against malicious characters.
	// Note: Implement the sanitizeFilename function based on your requirements.
	filename = utils.SanitizeFilename(filename)

	// Remove the extension from the filename.
	filename = filename[:len(filename)-len(filepath.Ext(filename))]

	// If randomizeFilename is true, prepend a uuid to the filename.
	if randomizeFilename {
		filename = fmt.Sprintf("%s-%s", uuid.New().String(), filename)
	}

	resp, err := cld.Upload.Upload(ctx, bytes.NewReader(buffer), uploader.UploadParams{
		PublicID: filename,
	})
	if err != nil {
		fmt.Println("Error uploading file to Cloudinary:", err)
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadBulkToCloudinary(
	ctx context.Context,
	files []BulkFileUploadRequest,
	opts ...bool,
) ([]string, error) {

	randomizeFilename := true
	if len(opts) > 0 {
		randomizeFilename = opts[0]
	}

	// WaitGroup is used to wait for all goroutines to finish.

	var wg sync.WaitGroup
	resultUrls := make([]string, len(files))
	errors := make(chan error, len(files))

	/*
		wg.Add(1) is used to add 1 to the WaitGroup counter.
			it waits for the 1 goroutine to finish as the number
			of goroutines increases, we need to increase the counter
		wg.Done() is used to decrement the WaitGroup counter by 1.
			it signals that 1 goroutine has finished
		wg.Wait() is used to wait for all goroutines to finish.
			it stops the main thread of the app to wait until the WaitGroup counter becomes 0
			until then the main thread of the app continues
	*/

	/*
		go keyword is used to create a new goroutine.
		goroutine is a lightweight thread of execution that can run concurrently with other goroutines.
		https://en.wikipedia.org/wiki/Goroutine
		so it creates a new thread of execution.
	*/

	for i, file := range files {

		wg.Add(1)

		go func(i int, file BulkFileUploadRequest) {

			//defer keyword is used to define a function that should be called after the surrounding function returns.

			defer wg.Done()
			url, err := UploadBufferToCloudinary(ctx, file.buffer, file.filename, randomizeFilename)
			if err != nil {
				errors <- err
				return
			}
			resultUrls[i] = url
		}(i, file)
	}

	wg.Wait()
	// Close the channel to signal that all goroutines have finished
	close(errors)

	if len(errors) > 0 {
		// Return the first error encountered
		return nil, <-errors
	}

	return resultUrls, nil
}
