package providers

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/deadkrolik/godbt/contract"
)

func TestSourceToString(t *testing.T) {
	var (
		err  error
		path string
	)
	p := GetXMLImageProvider()
	dir, _ := filepath.Abs("./testdata")

	path = dir + "/xml_valid.xml"
	_, err = p.sourceToString(path)
	if err != nil {
		t.Fatal("sourceToPath should not return error for the partially valid xml file")
	}

	path = dir + "/xml_invalid.xml"
	_, err = p.sourceToString(path)
	if err == nil {
		t.Fatal("sourceToPath should return error for invalid xml file")
	}

	_, err = p.sourceToString("not xml")
	if err == nil {
		t.Fatal("sourceToPath should return error for invalid xml string")
	}

	_, err = p.sourceToString(`<?xml version="1.0" ?><dataset> ...`)
	if err != nil {
		t.Fatal("sourceToPath should not return error for partially valid xml string")
	}
}

func TestToImage(t *testing.T) {
	var (
		err    error
		reader bytes.Buffer
		image  contract.Image
		ok     bool
	)
	p := GetXMLImageProvider()

	reader.Reset()
	reader.WriteString("invalid xml<<<<<>>>>>>>>>>>>")
	_, err = p.toImage(&reader)
	if err == nil {
		t.Fatal("toImage should return error for invalid xml payload")
	}

	reader.Reset()
	reader.WriteString(`<?xml version="1.0" ?>
<dataset>
    <keys1 id="1" data="test1"/>
    <keys2 id="2" data="test2"/>
</dataset>`)
	image, err = p.toImage(&reader)
	if err != nil {
		t.Fatal("toImage should not return error for valid xml payload")
	}

	if len(image) != 2 {
		t.Fatalf("Image len should be 2, got %d", len(image))
	}

	if image[0].Table != "keys1" || image[1].Table != "keys2" {
		t.Fatalf("Image tables should be keys1 and keys2, not %s and %s", image[0].Table, image[1].Table)
	}

	id1, ok := image[0].Data["id"]
	if !ok {
		t.Fatalf("Image[0] doesn't contain key `id`")
	}
	if id1 != "1" {
		t.Fatalf("Image[0].id is not equal to `1`, but `%s`", id1)
	}
}
