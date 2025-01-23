package bimg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestDeterminateImageType(t *testing.T) {
	files := []struct {
		name       string
		expected   ImageType
		shouldTest bool
	}{
		{"test.jpg", JPEG, true},
		{"test.png", PNG, true},
		{"test.webp", WEBP, true},
		{"test.gif", GIF, true},
		{"test.pdf", PDF, true},
		{"test.svg", SVG, true},
		{"test.jp2", JP2K, vipsVersionMin(8, 11)},
		{"test.jxl", JXL, vipsVersionMin(8, 11)},
		{"test.heic", HEIF, true},
		{"test2.heic", HEIF, true},
		{"test3.heic", HEIF, true},
		{"test.avif", AVIF, true},
		{"test.bmp", MAGICK, vipsVersionMin(8, 13)},
	}

	for _, file := range files {
		if !file.shouldTest {
			t.Skip("condition not met")
		}
		t.Run(file.name, func(t *testing.T) {
			img, _ := os.Open(path.Join("testdata", file.name))
			buf, _ := ioutil.ReadAll(img)
			defer img.Close()

			if VipsIsTypeSupported(file.expected) {
				value := DetermineImageType(buf)
				if value != file.expected {
					t.Fatalf("Image type is not valid: wanted %s, got: %s", ImageTypes[file.expected], ImageTypes[value])
				}
			}
		})
	}
}

func TestDeterminateImageTypeName(t *testing.T) {
	files := []struct {
		name      string
		expected  string
		condition bool
	}{
		{"test.jpg", "jpeg", true},
		{"test.png", "png", true},
		{"test.webp", "webp", true},
		{"test.gif", "gif", true},
		{"test.pdf", "pdf", vipsVersionMin(8, 12)},
		{"test.svg", "svg", true},
		{"test.jp2", "jp2k", vipsVersionMin(8, 11)},
		{"test.jxl", "jxl", vipsVersionMin(8, 11)},
		{"test.heic", "heif", true},
		{"test.avif", "avif", true},
		{"test.bmp", "magick", vipsVersionMin(8, 13)},
	}

	for _, file := range files {
		t.Run(file.name, func(t *testing.T) {
			if !file.condition {
				t.Skip("condition not met")
			}

			img, _ := os.Open(path.Join("testdata", file.name))
			buf, _ := ioutil.ReadAll(img)
			defer img.Close()

			value := DetermineImageTypeName(buf)
			if value != file.expected {
				t.Fatalf("Image type is not valid: %s != %s, got: %s", file.name, file.expected, value)
			}
		})

	}
}

func TestIsTypeSupported(t *testing.T) {
	types := []struct {
		name      ImageType
		supported bool
	}{
		{JPEG, true},
		{PNG, true},
		{WEBP, true},
		{GIF, true},
		{PDF, vipsVersionMin(8, 12)},
		{HEIF, true},
		{AVIF, true},
		{JP2K, vipsVersionMin(8, 11)},
		{JXL, vipsVersionMin(8, 11)},
	}

	for _, typ := range types {
		t.Run(ImageTypes[typ.name], func(t *testing.T) {
			if IsTypeSupported(typ.name) != typ.supported {
				t.Fatalf("Image type support is not as expected")
			}
		})
	}
}

func TestIsTypeNameSupported(t *testing.T) {
	types := []struct {
		name      string
		expected  bool
		condition bool
	}{
		{"jpeg", true, true},
		{"png", true, true},
		{"webp", true, true},
		{"gif", true, true},
		{"pdf", true, vipsVersionMin(8, 12)},
		{"heif", true, true},
		{"avif", true, true},
		{"jp2k", true, vipsVersionMin(8, 11)},
		{"jxl", true, vipsVersionMin(8, 11)},
	}

	for _, n := range types {
		t.Run(n.name, func(t *testing.T) {
			if !n.condition {
				t.Skip("condition not met")
			}
			if IsTypeNameSupported(n.name) != n.expected {
				t.Fatalf("Image type %s is not valid", n.name)
			}
		})
	}
}

func TestIsTypeSupportedSave(t *testing.T) {
	types := []ImageType{
		JPEG, PNG, WEBP, TIFF, HEIF, AVIF,
	}
	if vipsVersionMin(8, 11) {
		types = append(types, JP2K, JXL)
	}
	if vipsVersionMin(8, 12) {
		types = append(types, GIF)
	}

	for _, tt := range types {
		if IsTypeSupportedSave(tt) == false {
			t.Fatalf("Image type %s is not valid", ImageTypes[tt])
		}
	}
}

func TestIsTypeNameSupportedSave(t *testing.T) {
	types := []struct {
		name     string
		expected bool
	}{
		{"jpeg", true},
		{"png", true},
		{"webp", true},
		{"gif", vipsVersionMin(8, 12)},
		{"pdf", false},
		{"tiff", true},
		{"heif", true},
		{"avif", true},
		{"jp2k", vipsVersionMin(8, 11)},
		{"jxl", vipsVersionMin(8, 11)},
	}

	for _, n := range types {
		t.Run(n.name, func(t *testing.T) {
			if IsTypeNameSupportedSave(n.name) != n.expected {
				t.Fatalf("Image type is not valid (expected = %t)", n.expected)
			}
		})
	}
}
