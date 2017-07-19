package providers

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/deadkrolik/godbt/contract"
)

//XMLImageProvider - parser for XML-files
type XMLImageProvider struct {
}

const xmlRootElement = "dataset"

//GetXMLImageProvider - parser instance
func GetXMLImageProvider() *XMLImageProvider {
	return &XMLImageProvider{}
}

//Parse - real parse
func (p *XMLImageProvider) Parse(source interface{}) (contract.Image, error) {
	data, err := p.sourceToString(source)
	if err != nil {
		return contract.Image{}, err
	}

	var reader bytes.Buffer
	reader.WriteString(data)

	return p.toImage(&reader)
}

//CanParse - could we parse it
func (p *XMLImageProvider) CanParse(source interface{}) bool {
	_, err := p.sourceToString(source)

	return err == nil
}

//sourceToString - convert interface{} to string with xml data
func (p *XMLImageProvider) sourceToString(source interface{}) (string, error) {
	src, ok := source.(string)
	if !ok {
		return "", errors.New("source param is not a string")
	}

	if strings.Contains(src, "<?xml") {
		return src, nil
	}

	ext := strings.ToLower(filepath.Ext(src))
	if ext == ".xml" {
		bytes, err := ioutil.ReadFile(src)
		if err != nil {
			return "", err
		}

		bytesString := string(bytes)
		if strings.Contains(bytesString, "<?xml") {
			return bytesString, nil
		}
	}

	return "", errors.New("source param is not a valid source")
}

//toImage - read and parse
func (p *XMLImageProvider) toImage(file io.Reader) (contract.Image, error) {
	var image contract.Image
	decoder := xml.NewDecoder(file)
	canCollect := false

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return image, err
		}
		if token == nil {
			return image, errors.New("Token is invalid")
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == xmlRootElement {
				canCollect = true
				continue
			}
			if canCollect {
				stmt := contract.Row{Table: t.Name.Local}
				stmt.Data = make(map[string]string, len(t.Attr))
				for _, attr := range t.Attr {
					stmt.Data[attr.Name.Local] = attr.Value
				}
				image = append(image, stmt)
			}
		case xml.EndElement:
			if t.Name.Local == xmlRootElement {
				canCollect = false
				break
			}
		}
	}

	return image, nil
}
