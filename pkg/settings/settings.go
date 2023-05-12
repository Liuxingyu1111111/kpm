package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"kusionstack.io/kpm/pkg/env"
	"kusionstack.io/kpm/pkg/errors"
)

// The config.json used to persist user information
const CONFIG_JSON_PATH = ".kpm/config/config.json"

// The kpm.json used to describe the default configuration of kpm.
const KPM_JSON_PATH = ".kpm/config/kpm.json"

// The kpm configuration
type KpmConf struct {
	DefaultOciRegistry string
	DefaultOciRepo     string
}

const DEFAULT_REGISTRY = "ghcr.io"
const DEFAULT_REPO = "KusionStack"

// This is a singleton that loads kpm settings from 'kpm.json'
// and is only initialized on the first call by 'Init()' or 'GetSettings()'
var kpm_settings *Settings = nil

// DefaultKpmConf create a default configuration for kpm.
func DefaultKpmConf() KpmConf {
	return KpmConf{
		DefaultOciRegistry: DEFAULT_REGISTRY,
		DefaultOciRepo:     DEFAULT_REPO,
	}
}

type Settings struct {
	CredentialsFile string
	// the default configuration for kpm.
	Conf KpmConf
}

// DefaultOciRepo return the default OCI registry 'ghcr.io'.
func (settings *Settings) DefaultOciRegistry() string {
	return settings.Conf.DefaultOciRegistry
}

// DefaultOciRepo return the default OCI repo 'KusionStack'.
func (settings *Settings) DefaultOciRepo() string {
	return settings.Conf.DefaultOciRepo
}

// GetFullJsonPath returns the full path of 'config.json' and 'kpm.json' file path under '$HOME/.kpm/config/'
func GetFullJsonPath(jsonFileName string) (string, error) {
	home, err := env.GetAbsPkgPath()
	if err != nil {
		return "", errors.InternalBug
	}

	return filepath.Join(home, jsonFileName), nil
}

// Init returns default kpm settings load from '$KCL_PKG_PATH/.kpm/config/kpm.json'
// and '$KCL_PKG_PATH/.kpm/config/config.json'.
func Init() (*Settings, error) {
	if kpm_settings == nil {
		credentialsFile, err := GetFullJsonPath(CONFIG_JSON_PATH)
		if err != nil {
			return nil, err
		}

		conf, err := loadOrCreateDefaultKpmJson()
		if err != nil {
			return nil, err
		}

		kpm_settings = &Settings{
			CredentialsFile: credentialsFile,
			Conf:            *conf,
		}

		return kpm_settings, nil
	} else {
		return kpm_settings, nil
	}
}

// GetSettings will return the kpm setting singleton.
func GetSettings() (*Settings, error) {
	if kpm_settings == nil {
		return Init()
	} else {
		return kpm_settings, nil
	}
}

// loadOrCreateDefaultKpmJson will load the 'kpm.json' file from '$KCL_PKG_PATH/.kpm/config',
// and create a default 'kpm.json' file if the file does not exist.
func loadOrCreateDefaultKpmJson() (*KpmConf, error) {
	kpmConfpath, err := GetFullJsonPath(KPM_JSON_PATH)
	if err != nil {
		return nil, err
	}

	defaultKpmConf := DefaultKpmConf()

	b, err := ioutil.ReadFile(kpmConfpath)
	// if the file '$KCL_PKG_PATH/.kpm/config/kpm.json' does not exist
	if os.IsNotExist(err) {
		// create the default kpm.json.
		err = os.MkdirAll(filepath.Dir(kpmConfpath), 0755)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(defaultKpmConf)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(kpmConfpath, b, 0644)
		if err != nil {
			return nil, err
		}
		return &defaultKpmConf, nil
	} else {
		err = json.Unmarshal(b, &defaultKpmConf)
		if err != nil {
			return nil, err
		}
		return &defaultKpmConf, nil
	}
}
