package bundle_test

import (
	"bytes"
	"os"
	"testing"

	bundle "github.com/andrejzverev/mm_bundle"
)

func TestCreateBundle(t *testing.T) {
	b := &bundle.BundleTAR{}
	err := b.CreateBundle("test_bundle.bundle")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer func() {
		err := b.CloseBundle()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		err = os.Remove("test_bundle.bundle")
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()
}

func TestAddGetContent(t *testing.T) {
	b := &bundle.BundleTAR{}

	err := b.CreateBundle("test_bundle.bundle")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	defer func() {
		err := b.CloseBundle()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		err = os.Remove("test_bundle.bundle")
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	content := []byte("TEST DATA")
	err = b.AddContent("testfile.txt", content)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = b.OpenBundle("test_bundle.bundle")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	extractedContent, err := b.GetContent("testfile.txt")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !bytes.Equal(extractedContent, content) {
		t.Errorf("expected: %s, got: %s",
			content, extractedContent)
	}

}

func TestGetContentFileNotFound(t *testing.T) {
	b := &bundle.BundleTAR{}

	err := b.CreateBundle("test_bundle.bundle")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	defer func() {
		err := b.CloseBundle()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		err = os.Remove("test_bundle.bundle")
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()
	content := []byte("TEST DATA")
	err = b.AddContent("file1.txt", content)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = b.OpenBundle("test_bundle.bundle")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	_, err = b.GetContent("nonexistentfile.txt")
	if err == nil {
		t.Fatalf("must get error because no suchfile")
	}
}

func TestGetContentWithoutOpen(t *testing.T) {
	b := &bundle.BundleTAR{}
	_, err := b.GetContent("file1.txt")
	if err == nil {
		t.Fatalf("must get error because trying to read without open first")
	}
}

func TestCreateBundleTwice(t *testing.T) {
	b := &bundle.BundleTAR{}

	err := b.CreateBundle("test_bundle.bundle")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	defer func() {
		err := b.CloseBundle()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		err = os.Remove("test_bundle.bundle")
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		os.Remove("test_bundle2.bundle")

	}()

	err = b.CreateBundle("test_bundle2.bundle")
	if err == nil {
		t.Fatalf("err: %s", err)
	}
}
