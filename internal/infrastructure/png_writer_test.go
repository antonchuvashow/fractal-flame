package infrastructure

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPngWriter(t *testing.T) {
	t.Parallel()

	writer := NewPngWriter()
	assert.NotNil(t, writer)
	assert.IsType(t, &PngWriter{}, writer)
}

func TestPngWriter_Write(t *testing.T) {
	t.Parallel()

	writer := NewPngWriter()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	t.Run("should write a PNG file successfully", func(t *testing.T) {
		t.Parallel()

		// Setup: Create a temporary directory for the output file
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.png")

		// Execute
		err := writer.Write(img, filePath)

		// Assert
		require.NoError(t, err)
		_, err = os.Stat(filePath)
		assert.False(t, os.IsNotExist(err), "Expected output file to be created")
	})

	t.Run("should return an error for an invalid file path", func(t *testing.T) {
		t.Parallel()

		// Setup: An invalid path that cannot be created
		invalidPath := "/non_existent_dir/test.png"

		// Execute
		err := writer.Write(img, invalidPath)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})
}
