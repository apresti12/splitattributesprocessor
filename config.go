package splitAttributesProcessor

// Create a configuration struct if you want any custom settings
type Config struct {
	Delimiter    string `mapstructure:"delimiter"`
	AttributeKey string `mapstructure:"attribute_key"`
}

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
