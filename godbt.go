package godbt

import (
	"errors"

	"github.com/deadkrolik/godbt/contract"
	"github.com/deadkrolik/godbt/installers"
)

//Tester - simple DB tester
type Tester struct {
	imageManager *ImageManager
	installer    contract.Installer
}

//GetTester - tester instance
func GetTester(config contract.InstallerConfig) (*Tester, error) {
	var installer contract.Installer
	var err error

	switch config.Type {
	case "mysql":
		installer, err = installers.GetInstallerMysql(config)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown loader name: " + config.Type)
	}

	return &Tester{
		imageManager: getImageManager(),
		installer:    installer,
	}, nil
}

//GetImageManager - Image manager
func (tester *Tester) GetImageManager() *ImageManager {
	return tester.imageManager
}

//GetInstaller - installer instance
func (tester *Tester) GetInstaller() contract.Installer {
	return tester.installer
}
