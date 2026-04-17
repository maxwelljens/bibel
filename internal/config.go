package bible

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

// findBibleFileInDataDir looks for Bible JSON files in the XDG data directory
// It returns the path to the first JSON file found, or an error if none found
func findBibleFileInDataDir() (string, error) {
	dataDir := filepath.Join(xdg.DataHome, "bibel")
	
	// Check if the directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return "", fmt.Errorf("Bible data directory not found: %s\nPlease create this directory and add Bible JSON files to it", dataDir)
	}

	// Read the directory
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return "", fmt.Errorf("error reading Bible data directory %s: %w", dataDir, err)
	}

	// Look for JSON files
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			return filepath.Join(dataDir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no JSON Bible files found in %s\nPlease add Bible JSON files to this directory", dataDir)
}

// Config represents the application configuration
type Config struct {
	// Reading mode: "evangelion", "new_testament", "old_testament", "bible"
	ReadingMode string `toml:"reading_mode" mapstructure:"reading_mode"`
	// Output mode: "tui", "formatted", "plain"
	OutputMode string `toml:"output_mode" mapstructure:"output_mode"`

	// Bible data file path (default: "books/pol_nbg.json")
	BiblePath string `toml:"bible_path" mapstructure:"bible_path"`

	// TUI settings
	TUI struct {
		// Whether to show quit message (default: true)
		ShowQuitMessage bool `toml:"show_quit_message" mapstructure:"show_quit_message"`

		// Box border style: "rounded", "double", "single", "hidden"
		BorderStyle string `toml:"border_style" mapstructure:"border_style"`

		// Box border colour (default: adaptive to terminal)
		BorderColour string `toml:"border_colour" mapstructure:"border_colour"`

		// Header colour (default: adaptive to terminal)
		HeaderColour string `toml:"header_colour" mapstructure:"header_colour"`

		// Text colour (default: adaptive to terminal)
		TextColour string `toml:"text_colour" mapstructure:"text_colour"`

		// Quit message colour (default: adaptive to terminal)
		QuitColour string `toml:"quit_colour" mapstructure:"quit_colour"`
	} `toml:"tui" mapstructure:"tui"`

	// Formatter settings
	Formatter struct {
		// Whether to use ANSI colours in formatted mode (default: true)
		UseColours bool `toml:"use_colours" mapstructure:"use_colours"`

		// Header format template (default: "{book} {chapter}\nw. {first_verse}-{second_verse}")
		HeaderFormat string `toml:"header_format" mapstructure:"header_format"`
	} `toml:"formatter" mapstructure:"formatter"`

	// Date progression settings
	DateProgression struct {
		// Verses per day (default: 12)
		VersesPerDay int `toml:"verses_per_day" mapstructure:"verses_per_day"`

		// Start date for progression (default: "1 January" of current year)
		StartDate string `toml:"start_date" mapstructure:"start_date"`
	} `toml:"date_progression" mapstructure:"date_progression"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	cfg := &Config{
		ReadingMode: "evangelion",
		OutputMode: "tui",
		BiblePath:  "", // Empty means use XDG data directory
	}

	cfg.TUI.ShowQuitMessage = true
	cfg.TUI.BorderStyle = "rounded"
	cfg.TUI.BorderColour = "" // Empty means adaptive
	cfg.TUI.HeaderColour = "" // Empty means adaptive
	cfg.TUI.TextColour = ""   // Empty means adaptive
	cfg.TUI.QuitColour = ""   // Empty means adaptive

	cfg.Formatter.UseColours = true
	cfg.Formatter.HeaderFormat = "{book} {chapter}\nw. {first_verse}-{second_verse}"

	cfg.DateProgression.VersesPerDay = 12
	cfg.DateProgression.StartDate = "" // Empty means 1 January of current year

	return cfg
}

// LoadConfig loads configuration from file and environment
func LoadConfig() (*Config, error) {
	// Set up viper
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name

	// Add XDG config paths
	configPath := filepath.Join(xdg.ConfigHome, "bibel")
	viper.AddConfigPath(configPath)

	// Also check current directory for local config
	viper.AddConfigPath(".")

	// Set defaults
	defaultCfg := DefaultConfig()
	viper.SetDefault("reading_mode", defaultCfg.ReadingMode)
	viper.SetDefault("output_mode", defaultCfg.OutputMode)
	viper.SetDefault("bible_path", defaultCfg.BiblePath)
	viper.SetDefault("tui.show_quit_message", defaultCfg.TUI.ShowQuitMessage)
	viper.SetDefault("tui.border_style", defaultCfg.TUI.BorderStyle)
	viper.SetDefault("tui.border_colour", defaultCfg.TUI.BorderColour)
	viper.SetDefault("tui.header_colour", defaultCfg.TUI.HeaderColour)
	viper.SetDefault("tui.text_colour", defaultCfg.TUI.TextColour)
	viper.SetDefault("tui.quit_colour", defaultCfg.TUI.QuitColour)
	viper.SetDefault("formatter.use_colours", defaultCfg.Formatter.UseColours)
	viper.SetDefault("formatter.header_format", defaultCfg.Formatter.HeaderFormat)
	viper.SetDefault("date_progression.verses_per_day", defaultCfg.DateProgression.VersesPerDay)
	viper.SetDefault("date_progression.start_date", defaultCfg.DateProgression.StartDate)

	// Read in config file (if it exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; we'll use defaults
			fmt.Fprintf(os.Stderr, "Note: No configuration file found at %s/config.toml\n", configPath)
			fmt.Fprintf(os.Stderr, "Using default configuration. Run 'bibel --generate-config' to create a default config file.\n")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Apply defaults for any empty fields
	defaultCfg = DefaultConfig()
	if cfg.ReadingMode == "" {
		cfg.ReadingMode = defaultCfg.ReadingMode
	}
	if cfg.OutputMode == "" {
		cfg.OutputMode = defaultCfg.OutputMode
	}
	if cfg.BiblePath == "" {
		cfg.BiblePath = defaultCfg.BiblePath
	}
	if cfg.TUI.BorderStyle == "" {
		cfg.TUI.BorderStyle = defaultCfg.TUI.BorderStyle
	}
	if cfg.DateProgression.VersesPerDay == 0 {
		cfg.DateProgression.VersesPerDay = defaultCfg.DateProgression.VersesPerDay
	}

	// If BiblePath is still empty (default or from config), try to find a Bible file in XDG data directory
	if cfg.BiblePath == "" {
		biblePath, err := findBibleFileInDataDir()
		if err != nil {
			return nil, fmt.Errorf("failed to find Bible data file: %w", err)
		}
		cfg.BiblePath = biblePath
		fmt.Fprintf(os.Stderr, "Using Bible file from XDG data directory: %s\n", cfg.BiblePath)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(cfg *Config) error {
	// Set up viper with current config
	viper.Set("reading_mode", cfg.ReadingMode)
	viper.Set("output_mode", cfg.OutputMode)
	viper.Set("bible_path", cfg.BiblePath)
	viper.Set("tui.show_quit_message", cfg.TUI.ShowQuitMessage)
	viper.Set("tui.border_style", cfg.TUI.BorderStyle)
	viper.Set("tui.border_colour", cfg.TUI.BorderColour)
	viper.Set("tui.header_colour", cfg.TUI.HeaderColour)
	viper.Set("tui.text_colour", cfg.TUI.TextColour)
	viper.Set("tui.quit_colour", cfg.TUI.QuitColour)
	viper.Set("formatter.use_colours", cfg.Formatter.UseColours)
	viper.Set("formatter.header_format", cfg.Formatter.HeaderFormat)
	viper.Set("date_progression.verses_per_day", cfg.DateProgression.VersesPerDay)
	viper.Set("date_progression.start_date", cfg.DateProgression.StartDate)

	// Ensure config directory exists
	configPath := filepath.Join(xdg.ConfigHome, "bibel")
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Write config file
	configFile := filepath.Join(configPath, "config.toml")
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// GenerateDefaultConfig creates a default configuration file
func GenerateDefaultConfig() error {
	return SaveConfig(DefaultConfig())
}

// GetConfigPath returns the path to the configuration file

// Validate checks if configuration values are valid
func (c *Config) Validate() error {
	// Validate reading mode
	validReadingModes := map[string]bool{
		"evangelion": true,
		"new_testament": true,
		"old_testament": true,
		"bible": true,
	}
	// Allow empty reading mode (will use default)
	if c.ReadingMode != "" && !validReadingModes[c.ReadingMode] {
		return fmt.Errorf("invalid reading mode: %s, must be one of: evangelion, new_testament, old_testament, bible", c.ReadingMode)
	}

	// Validate output mode
	validOutputModes := map[string]bool{
		"tui":       true,
		"formatted": true,
		"plain":     true,
	}
	if !validOutputModes[c.OutputMode] {
		return fmt.Errorf("invalid output mode: %s, must be one of: tui, formatted, plain", c.OutputMode)
	}

	// Validate TUI border style
	validBorderStyles := map[string]bool{
		"rounded": true,
		"double":  true,
		"single":  true,
		"hidden":  true,
	}
	if !validBorderStyles[c.TUI.BorderStyle] {
		return fmt.Errorf("invalid border style: %s, must be one of: rounded, double, single, hidden", c.TUI.BorderStyle)
	}

	// Validate verses per day
	if c.DateProgression.VersesPerDay <= 0 {
		return fmt.Errorf("verses per day must be positive, got: %d", c.DateProgression.VersesPerDay)
	}

	return nil
}
func GetConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "bibel", "config.toml")
}