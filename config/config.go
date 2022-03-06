package config

// Excellent example and tutorial for YAML configuraiton file
// setup in Go projects:
// https://github.com/koddr/example-go-config-yaml

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

// Config file layout should define what values are
// applicable from a YAML file read from disk.
type Config struct {
	NOAA struct {
		Latitude  string `yaml:"latitude"`
		Longitude string `yaml:"longitude"`
	} `yaml:"noaa"`
}

// Factory method that creates a new instance of the
// configuration values injected at runtime. This reads
// the .yml from disk and parses into a Config struct.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ValidateConfigPath(path) checks to make sure the supplied
// path exists and also whether or not its a directory.
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return errors.New("config path must not be a folder")
	}
	return nil
}

// ParseFlags() reads in expected command line arguments from
// when the program was launched and returns the path if there
// were no errors. Defaults to looking for config.yml in the
// current working directory if no -config argument was sent.
func ParseFlags() (string, string, string, error) {
	var path string
	var lat string
	var lon string
	flag.StringVar(&path, "config", "", "path to YAML config file.")
	flag.StringVar(&lat, "lat", "41.837", "latitude of noaa location") // default to Chicago
	flag.StringVar(&lon, "lon", "-87.685", "longitude of noaa location")
	flag.Parse()
	if path != "" {
		if err := ValidateConfigPath(path); err != nil {
			return "", "", "", err
		}
	}
	return path, lat, lon, nil
}
