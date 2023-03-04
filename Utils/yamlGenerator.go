package util

import (
	"gopkg.in/yaml.v3"
)

type Pubspec struct {
	Name             string               `yaml:"name"`
	Description      string               `yaml:"description"`
	Publish_to       string               `yaml:"publish_to"`
	Version          string               `yaml:"version"`
	Environment      Sdk_version          `yaml:"environment"`
	Dependencies     Flutter              `yaml:"dependencies"`
	Dev_dependencies Flutter_test         `yaml:"dev_dependencies"`
	Flutter          Uses_material_design `yaml:"flutter"`
	Flutter_icons    Android              `yaml:"flutter_icons"`
}

type Sdk_version struct {
	Sdk string `yaml:"sdk"`
}

type Flutter struct {
	Flutter                 Sdk
	Flutter_dotenv          string `yaml:"flutter_dotenv"`
	Flutter_config          string `yaml:"flutter_config"`
	Cupertino_icons         string `yaml:"cupertino_icons"`
	Webview_flutter         string `yaml:"webview_flutter"`
	Lottie                  string `yaml:"lottie"`
	Url_launcher            string `yaml:"url_launcher"`
	Flutter_launcher_icons  string `yaml:"flutter_launcher_icons"`
	Change_app_package_name string `yaml:"change_app_package_name"`
}

type Flutter_test struct {
	Flutter_test  Sdk    `yaml:"flutter_test"`
	Flutter_lints string `yaml:"flutter_lints"`
}

type Uses_material_design struct {
	Uses_material_design bool     `yaml:"uses-material-design"`
	Assets               []string `yaml:"assets"`
}

type Android struct {
	Android    string `yaml:"android"`
	Ios        bool   `yaml:"ios"`
	Image_path string `yaml:"image_path"`
}

type Sdk struct {
	Sdk string `yaml:"sdk"`
}

func YamlGenerator(version string) []byte{

	pubspecfile := Pubspec{
		Name:        "boost_msme_app_builder",
		Description: "A new Flutter project",
		Publish_to:  "none",
		Version:     "1.0.0+"+version,
		Environment: Sdk_version{
			Sdk: ">=2.16.2 <3.0.0",
		},
		Dependencies: Flutter{
			Flutter: Sdk{
				Sdk: "flutter",
			},
			Flutter_dotenv:          "^5.0.2",
			Flutter_config:          "^2.0.0",
			Cupertino_icons:         "^1.0.2",
			Webview_flutter:         "^3.0.4",
			Lottie:                  "^1.2.2",
			Url_launcher:            "^6.1.0",
			Flutter_launcher_icons:  "^0.9.2",
			Change_app_package_name: "^1.0.0",
		},
		Dev_dependencies: Flutter_test{
			Flutter_test: Sdk{
				Sdk: "flutter",
			},
			Flutter_lints: "^1.0.0",
		},
		Flutter: Uses_material_design{
			Uses_material_design: true,
			Assets:               []string{"assets/anim/", "assets/logo/", ".env"},
		},
		Flutter_icons: Android{
			Android:    "launcher_icon",
			Ios:        true,
			Image_path: "assets/logo/icon.png",
		},
	}

	yamlData, _ := yaml.Marshal(&pubspecfile)

	return yamlData
}