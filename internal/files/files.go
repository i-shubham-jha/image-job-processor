package files

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"retail_pulse/internal/logger"

	"github.com/google/uuid"
)

// ImageHolder is a struct that holds an image and its metadata.
type ImageHolder struct {
	ID     string
	Image  image.Image
	Width  int
	Height int
	Format string // Format of the image (e.g., "png", "jpeg")
}

// DownloadImage downloads an image from the specified URL and returns an ImageHolder.
func DownloadImage(url string) (*ImageHolder, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: received status code %d", resp.StatusCode)
	}

	// Read the response body
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Decode the image
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get image dimensions
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Generate a new UUID and convert it to a string
	id := uuid.New().String()

	logger.GetLogger().Log(fmt.Sprintf("Downloaded image from %v", url))

	return &ImageHolder{
		ID:     id,
		Image:  img,
		Width:  width,
		Height: height,
		Format: format,
	}, nil
}

// SaveImage saves the image to a file in the specified format (PNG or JPEG).
func (ih *ImageHolder) SaveImage(documentID, storeID string) error {
	if documentID == "" || storeID == "" {
		return fmt.Errorf("documentID and storeID must be provided")
	}

	// Create the directory structure
	dirPath := filepath.Join("./files", documentID, storeID)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Construct the file path
	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", ih.ID, ih.Format))

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Determine the format based on the Format field
	switch ih.Format {
	case "png":
		logger.GetLogger().Log(fmt.Sprintf("Saving file %v.png", ih.ID))
		return png.Encode(outFile, ih.Image)
	case "jpeg":
		logger.GetLogger().Log(fmt.Sprintf("Saving file %v.jpeg", ih.ID))
		return jpeg.Encode(outFile, ih.Image, nil)
	default:
		return fmt.Errorf("unsupported image format: %s", ih.Format)
	}
}
