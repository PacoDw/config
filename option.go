package config

// Option represents the option to configure the service.
type Option func(*Config)

// ApplyOptions applies the options to the Config.
func ApplyOptions(s *Config, opts []Option) {
	for _, opt := range opts {
		opt(s)
	}
}

// WithFilePath sets the configuration file path.
func WithFilePath(filePath string) Option {
	return func(c *Config) {
		c.filePath = filePath
	}
}

// WithFileName sets the configuration file name without extension.
func WithFileName(fileName string) Option {
	return func(c *Config) {
		c.fileName = fileName
	}
}

// WithFileType sets the configuration file type.
func WithFileType(fileType string) Option {
	return func(c *Config) {
		c.fileType = fileType
	}
}
