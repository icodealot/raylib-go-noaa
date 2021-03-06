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

// -----------------------------------------------------------------
// Parse command line arguments and read config file
// if one is provided. Command line arguments include:
//    -config <path to config file>
//    -lat <latitude>
//    -lon <longitude>
// Note: if a config file is provided the lat and lon
// arguments are ignored.
//
// Defaults to Chicago, IL (LOT) if no arguments are provided.
// -----------------------------------------------------------------
func NewConfig() (*Config, error) {
	cfg := &Config{}
	cfgPath, lat, lon, uom, err := parseFlags()
	if err != nil {
		return nil, err
	}
	if cfgPath != "" {
		cfg, err = newConfig(cfgPath)
		if err != nil {
			return nil, err
		}
	} else {
		cfg.NOAA.Latitude = lat
		cfg.NOAA.Longitude = lon
	}
	if uom == "si" {
		cfg.NOAA.Units = "si"
	}
	return cfg, nil
}

// Config file layout should define what values are
// applicable from a YAML file read from disk.
type Config struct {
	NOAA struct {
		Latitude  string `yaml:"latitude"`
		Longitude string `yaml:"longitude"`
		Units     string `yaml:"units"`
	} `yaml:"noaa"`
}

// Factory method that creates a new instance of the
// configuration values injected at runtime. This reads
// the .yml from disk and parses into a Config struct.
func newConfig(path string) (*Config, error) {
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
func validateConfigPath(path string) error {
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
func parseFlags() (string, string, string, string, error) {
	// TODO: probably worth converting this to a config/struct now.
	var path string
	var lat string
	var lon string
	var uom string
	flag.StringVar(&path, "config", "", "path to YAML config file.")
	flag.StringVar(&lat, "lat", "41.837", "latitude of noaa location") // default to Chicago
	flag.StringVar(&lon, "lon", "-87.685", "longitude of noaa location")
	flag.StringVar(&uom, "uom", "", "unit of measure can be si or us (the default)")
	flag.Parse()
	if path != "" {
		if err := validateConfigPath(path); err != nil {
			return "", "", "", "", err
		}
	}
	return path, lat, lon, uom, nil
}
