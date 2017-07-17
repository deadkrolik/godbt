package godbt

import (
	"errors"
	"fmt"
	"github.com/deadkrolik/godbt/contract"
	"github.com/deadkrolik/godbt/providers"
	"strings"
)

//ImageManager - manager for Image type
type ImageManager struct {
	providers map[string]contract.ImageProvider
}

//GetImageManager - get an instance
func getImageManager() *ImageManager {
	imageProviders := make(map[string]contract.ImageProvider)
	imageProviders["xml"] = providers.GetXMLImageProvider()

	return &ImageManager{
		providers: imageProviders,
	}
}

//LoadImage - setup Image to real DB
func (manager *ImageManager) LoadImage(args ...interface{}) (contract.Image, error) {
	var source interface{}
	var modifiers contract.ModifiersList
	var ok bool

	if len(args) == 0 {
		return contract.Image{}, errors.New("LoadImage requires at least 1 param")
	}
	if len(args) > 0 {
		source = args[0]
	}
	if len(args) > 1 {
		modifiers, ok = args[1].(contract.ModifiersList)
		if !ok {
			return contract.Image{}, errors.New("Second param for LoadImage must have `map[string]Modifier` type")
		}
	}

	parser := manager.getLoaderBySource(source)
	if parser == nil {
		return contract.Image{}, errors.New("Can't find dataset parser")
	}

	image, err := parser.Parse(source)
	if err != nil {
		return contract.Image{}, err
	}

	if len(modifiers) > 0 {
		image = manager.applyModifiers(image, modifiers)
	}

	return image, nil
}

//GetImagesDiff - simple images diff
func (manager *ImageManager) GetImagesDiff(left, right contract.Image) []string {
	var diffs []string

	min := len(left)
	if len(right) > len(left) {
		diffs = append(diffs, fmt.Sprintf(
			"RightImage is bigger than LeftImage (`%d` > `%d`)",
			len(right), len(left),
		))
		min = len(left)
	}
	if len(left) > len(right) {
		diffs = append(diffs, fmt.Sprintf(
			"LeftImage is bigger than RightImage (`%d` > `%d`)",
			len(left), len(right),
		))
		min = len(right)
	}

	for i := 0; i < min; i++ {
		if left[i].Table != right[i].Table {
			diffs = append(diffs, fmt.Sprintf(
				"position `%d`: table names are not equal (`%s` != `%s`)",
				i, left[i].Table, right[i].Table,
			))
			continue
		}

		diffs = append(diffs, manager.compareMap(
			fmt.Sprintf("position `%d`:", i),
			left[i].Data,
			"LeftImage",
			right[i].Data,
			"RightImage",
		)...)

		diffs = append(diffs, manager.compareMap(
			fmt.Sprintf("position `%d`:", i),
			right[i].Data,
			"RightImage",
			left[i].Data,
			"LeftImage",
		)...)
	}

	return diffs
}

//compareMap - compare by keys
func (manager *ImageManager) compareMap(prefix string, m1 map[string]string, m1name string, m2 map[string]string, m2name string) []string {
	var diffs []string

	for k1, v1 := range m1 {
		v2, ok := m2[k1]
		if !ok {
			diffs = append(diffs, fmt.Sprintf(
				"%s key `%s` in %s is not exists in %s",
				prefix, k1, m1name, m2name,
			))
			continue
		}

		if v1 != v2 {
			diffs = append(diffs, fmt.Sprintf(
				"%s key `%s` in %s is not equal to such key in %s (`%s` != `%s`)",
				prefix, k1, m1name, m2name, v1, v2,
			))
			continue
		}
	}

	return diffs
}

//getLoaderBySource - choosing a parser for dataset source
func (manager *ImageManager) getLoaderBySource(source interface{}) contract.ImageProvider {
	for _, parser := range manager.providers {
		if parser.CanParse(source) {
			return parser
		}
	}

	return nil
}

//applyModifiers - changin fields values
func (manager *ImageManager) applyModifiers(statements contract.Image, modifiers contract.ModifiersList) contract.Image {
	for sIndex, stmt := range statements {
		for k, v := range stmt.Data {
			for key, modifier := range modifiers {
				if strings.Contains(v, key) {
					newValue := modifier(stmt.Table, k, v)
					statements[sIndex].Data[k] = newValue
				}
			}
		}
	}

	return statements
}
