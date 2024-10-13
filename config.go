package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	// defaultFilePath is the default configuration file path.
	defaultFilePath = "./"

	// defaultFileName is the default configuration file name without extension.
	defaultFileName = ".env"

	// defaultFileType is the default configuration file type.
	defaultFileType = "yaml"
)

// Config is a wrapper around viper.
type Config struct {
	v *viper.Viper

	// filePath is the configuration file path.
	filePath string

	// fileName is the configuration file name without extension.
	fileName string

	// fileType is the configuration file type.
	fileType string
}

// New creates a new Config.
func New(opts ...Option) *Config {
	c := &Config{
		v:        viper.New(),
		filePath: defaultFilePath,
		fileName: defaultFileName,
		fileType: defaultFileType,
	}

	// apply options
	ApplyOptions(c, opts)

	// Set the config file
	c.v.AddConfigPath(c.filePath)
	c.v.SetConfigName(c.fileName)
	c.v.SetConfigType(c.fileType)

	// Enable VIPER to read Environment Variables
	c.v.AutomaticEnv()

	// Try to read the config file
	c.v.ReadInConfig()

	return c
}

// Unmarshal reads the configuration from the environment variables and the config file.
func (c *Config) Unmarshal(config interface{}) error {
	// Get all settings from Viper (from both env and the file) and apply global env settings
	allSettings := applyGlobalEnvSettings(c.v.AllSettings())

	// Decode settings into the provided config structure
	if err := decodeConfig(allSettings, config); err != nil {
		return err
	}

	// Use Viper's Unmarshal to handle environment variables with the custom DecodeHook
	if err := c.v.Unmarshal(config, viper.DecodeHook(mapstructureDecodeHook(config))); err != nil {
		return err
	}

	// Validate required fields using go-playground/validator
	if err := validateConfig(config); err != nil {
		return err
	}

	// Set default values for any missing fields
	return defaults.Set(config)
}

// decodeConfig decodes the provided settings map into the given config structure.
func decodeConfig(settings map[string]interface{}, config interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true, // Allow flexible type matching
		ZeroFields:       true, // Zero fields before decoding
		Result:           config,
		TagName:          "env", // Use `env` tags for field mapping
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}

	return decoder.Decode(settings)
}

// applyGlobalEnvSettings applies global environment variables to all settings.
func applyGlobalEnvSettings(allSettings map[string]interface{}) map[string]interface{} {
	// Get all global environment variables
	globalEnvs := make(map[string]interface{})
	for k, v := range allSettings {
		switch v := v.(type) {
		case map[string]interface{}:
			continue
		default:
			globalEnvs[k] = v
		}
	}

	// Apply global environment variables to all settings
	for _, v := range allSettings {
		switch v := v.(type) {
		case map[string]interface{}:
			for gKey, gVal := range globalEnvs {
				if _, ok := v[gKey]; !ok {
					v[gKey] = gVal
				}
			}
		}
	}

	return allSettings
}

// mapstructureDecodeHook handles custom decoding logic for environment variables
func mapstructureDecodeHook(config interface{}) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// If it's a map, try to match it with the structure name
		if f.Kind() == reflect.Map && data != nil {
			v, ok := data.(map[string]interface{})[strings.ToLower(t.Name())]
			if !ok {
				return data, nil
			}

			// Decode the map into the structure using mapstructure
			if err := decodeConfig(v.(map[string]interface{}), config); err != nil {
				return nil, err
			}

			return config, nil
		}

		return data, nil
	}
}

// validateConfig validates the provided config structure using go-playground/validator
func validateConfig(config interface{}) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("validation error: field '%s' is %s", err.Field(), err.Tag()))
		}

		return fmt.Errorf("errors: %s", strings.Join(errorMessages, ", "))
	}

	return nil
}
