package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultPort    = 5432
	configDirName  = ".sql-proxy"
	configFileName = "config.yaml"
	minValidPort   = 1
	maxValidPort   = 65535
)

var (
	ErrMalformedConfig = errors.New("malformed config file")
	ErrMissingInstance = errors.New("missing instance")
	ErrInvalidPort     = errors.New("invalid port")
	ErrConflictingIPMode = errors.New("conflicting ip mode flags")
)

// Settings holds the runtime configuration resolved for startup.
type Settings struct {
	Port       int
	Instance   string
	UsePrivateIP bool
	ConfigDir  string
	ConfigFile string
}

// Init resolves settings using flags > config file > defaults.
func Init() (Settings, error) {
	return InitWithArgs(os.Args[1:])
}

// InitWithArgs resolves settings using provided CLI args for testability.
func InitWithArgs(args []string) (Settings, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Settings{}, fmt.Errorf("resolve user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, configDirName)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return Settings{}, fmt.Errorf("create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, configFileName)

	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("private_ip", false)

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !errors.Is(err, os.ErrNotExist) {
			return Settings{}, fmt.Errorf("%w: %v", ErrMalformedConfig, err)
		}
	}

	fs := pflag.NewFlagSet("sql-proxy", pflag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.IntP("port", "p", viper.GetInt("port"), "Local port to bind")
	fs.StringP("instance", "i", viper.GetString("instance"), "Cloud SQL instance")
	fs.Bool("private-ip", viper.GetBool("private_ip"), "Use private IP connectivity")
	fs.Bool("public-ip", false, "Use public IP connectivity (default)")
	if err := fs.Parse(args); err != nil {
		return Settings{}, fmt.Errorf("parse flags: %w", err)
	}
	if err := viper.BindPFlags(fs); err != nil {
		return Settings{}, fmt.Errorf("bind flags: %w", err)
	}

	settings := Settings{
		Port:       viper.GetInt("port"),
		Instance:   viper.GetString("instance"),
		UsePrivateIP: viper.GetBool("private-ip"),
		ConfigDir:  configDir,
		ConfigFile: configFile,
	}

	privateChanged := fs.Lookup("private-ip") != nil && fs.Lookup("private-ip").Changed
	publicChanged := fs.Lookup("public-ip") != nil && fs.Lookup("public-ip").Changed
	if privateChanged && publicChanged {
		return Settings{}, fmt.Errorf("%w: use either --private-ip or --public-ip", ErrConflictingIPMode)
	}
	if publicChanged {
		settings.UsePrivateIP = false
	}

	if err := validate(settings); err != nil {
		return Settings{}, err
	}

	return settings, nil
}

func validate(s Settings) error {
	if s.Instance == "" {
		return fmt.Errorf("%w: provide --instance or set instance in %s", ErrMissingInstance, configFileName)
	}
	if s.Port < minValidPort || s.Port > maxValidPort {
		return fmt.Errorf("%w: port must be between %d and %d", ErrInvalidPort, minValidPort, maxValidPort)
	}
	return nil
}

// UserFacingError returns clear, actionable messages for known startup failures.
func UserFacingError(err error) string {
	switch {
	case errors.Is(err, ErrMalformedConfig):
		return "Invalid configuration file format. Fix ~/.sql-proxy/config.yaml and try again."
	case errors.Is(err, ErrMissingInstance):
		return "Missing instance. Use --instance or set instance in ~/.sql-proxy/config.yaml."
	case errors.Is(err, ErrInvalidPort):
		return "Invalid port. Use a value between 1 and 65535."
	case errors.Is(err, ErrConflictingIPMode):
		return "Invalid IP mode. Use either --private-ip or --public-ip."
	default:
		return fmt.Sprintf("Startup failed: %v", err)
	}
}
