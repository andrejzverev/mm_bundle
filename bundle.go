package bundle

// type IBundle interface {
// 	CreateBundle(bundleName string) error
// 	OpenBundle(bundleName string) error
// 	CloseBundle() error
// 	AddContent(fileName string, b []byte) error
// 	GetContent(fileName string) ([]byte, error)
// }

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
)

type BundleTAR struct {
	bundleFile *os.File
	tarWriter  *tar.Writer
	tarReader  *tar.Reader
	isOpen     bool
}

func (b *BundleTAR) CreateBundle(bundleName string) error {
	if b.bundleFile != nil {
		return fmt.Errorf("cat not create bundle already created")
	}

	file, err := os.Create(bundleName)
	if err != nil {
		return fmt.Errorf("failed to create bundle: %s", err)
	}
	b.bundleFile = file
	b.tarWriter = tar.NewWriter(file)
	b.isOpen = true
	return nil
}

func (b *BundleTAR) OpenBundle(bundleName string) error {
	file, err := os.Open(bundleName)
	if err != nil {
		return fmt.Errorf("failed to open bundle: %s", err)
	}
	b.bundleFile = file
	b.tarReader = tar.NewReader(file)
	b.isOpen = true
	return nil

}

func (b *BundleTAR) CloseBundle() error {
	if b.isOpen == true {
		if b.tarWriter != nil {
			if err := b.tarWriter.Close(); err != nil {
				return fmt.Errorf("failed to close bundle: %s", err)
			}
		}
		if b.bundleFile != nil {
			if err := b.bundleFile.Close(); err != nil {
				return fmt.Errorf("failed to close bundle: %s", err)
			}
		}
	}
	b.isOpen = false
	return nil
}

func (b *BundleTAR) AddContent(fileName string, data []byte) error {
	if b.isOpen == false || b.tarWriter == nil {
		return fmt.Errorf("bundle is uknown to proccess")
	}
	header := &tar.Header{
		Name: fileName,
		Mode: 0600,
		Size: int64(len(data)),
	}
	if err := b.tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header to bundle: %s", err)
	}
	_, err := b.tarWriter.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write content for %s to bundle: %s", fileName, err)
	}
	return nil
}

func (b *BundleTAR) GetContent(fileName string) ([]byte, error) {
	var content bytes.Buffer

	if b.isOpen == false || b.tarReader == nil {
		return nil, fmt.Errorf("bundle is uknown to proccess")
	}
	for {
		header, err := b.tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read bundle header: %s", err)
		}

		if header.Name == fileName {
			_, err := io.Copy(&content, b.tarReader)
			if err != nil {
				return nil, fmt.Errorf("failed to read content for %s from bundle: %s", fileName, err)
			}
			return content.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("unable to find file %s in bundle", fileName)
}
