package deploykeychain

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var repoNameRegexp = regexp.MustCompile(`^'?/?(.*)/(.*).git'?$`)

// ErrNoRepositoryFound occurs when IdentifyRepository is unable to find any argument identifying the repo being worked on.
var ErrNoRepositoryFound = errors.New("unable to identify repository")

// ErrNoKeyAvailable occurs when DetermineKeyFile is unable to locate a key to use for the repo being worked on.
var ErrNoKeyAvailable = errors.New("no key available for this repository")

// Config represents the structure of a config file.
type Config struct {
	// KeyPath is the path to a folder on disk where keys should be searched for, if no custom key file is provided.
	KeyPath string `mapstructure:"key_path"`
	// KeyNameFormat is the format of key file names, for repositories with no custom key file provided.
	KeyNameFormat string `mapstructure:"key_name_format"`
	// Keys is a map of repository names to the path to their key on disk.
	Keys map[string]string `mapstructure:"keys"`
	// SSHCommand is the name of the command to call to connect via SSH.
	SSHCommand string `mapstructure:"ssh_command"`
	// FallbackKey is the key that should be used if no other key is found. Leave blank to error out if no key can be found.
	FallbackKey string `mapstructure:"fallback_key"`
}

var debug bool

// Log logs a message to stderr when debug mode is enabled.
func Log(message string) {
	if debug {
		fmt.Fprintln(os.Stderr, message)
	}
}

// IdentifyRepository attempts to identify the repository being interacted with from the arguments provided to SSH by Git.
func IdentifyRepository(args []string) (string, string, error) {
	for _, argument := range args {
		argParts := strings.Split(argument, " ")
		for _, part := range argParts {
			match := repoNameRegexp.FindStringSubmatch(part)
			if len(match) == 3 {
				Log(fmt.Sprintf("Found repository details: %+v", match[1:]))
				return match[1], match[2], nil
			}
		}
	}

	return "", "", ErrNoRepositoryFound
}

// LoadConfig loads & populates the config for this tool.
func LoadConfig() (Config, error) {
	debugVar := os.Getenv("DEPLOY_KEYCHAIN_DEBUG")
	if debugVar == "" {
		debugVar = "false"
	}
	debug, _ = strconv.ParseBool(debugVar)

	viper.SetConfigName("deploy-keychain")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.deploy-keychain")
	viper.AddConfigPath(".")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("error getting user's home directory: %w", err)
	}

	viper.SetDefault("key_path", path.Join(homeDir, ".ssh", "deploy-keys"))
	viper.SetDefault("key_name_format", "{{.account}}-{{.repository}}.pem")
	viper.SetDefault("keys", make(map[string]string))
	viper.SetDefault("ssh_command", "ssh")

	err = viper.ReadInConfig()
	if err != nil {
		Log(fmt.Sprintf("Unable to load config file: %s\n", err))
	}

	config := new(Config)

	err = viper.Unmarshal(config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}

	Log(fmt.Sprintf("Config loaded: %+v", config))

	return *config, nil
}

// DetermineKeyFile will, given a repository, attempt to determine a SSH key to use for the repository.
func DetermineKeyFile(config Config, account, repository string) (string, error) {
	keyFile, found := config.Keys[fmt.Sprintf("%s/%s", account, repository)]
	if found {
		Log(fmt.Sprintf("Found key via custom keys map: %s", keyFile))
		return keyFile, nil
	}

	repositoryKeyNameTemplate, err := template.New("").Parse(config.KeyNameFormat)
	if err != nil {
		return "", fmt.Errorf("unable to create template for provided key name format: %w", err)
	}

	keyNameBuf := new(bytes.Buffer)
	repositoryKeyNameTemplate.Execute(keyNameBuf, map[string]string{"account": account, "repository": repository})
	keyName := keyNameBuf.String()
	keyPath := path.Join(config.KeyPath, keyName)

	Log(fmt.Sprintf("Generated key name: %s (Full path: %s)", keyName, keyPath))

	if _, err := os.Stat(keyPath); err == nil {
		Log(fmt.Sprintf("Found key via generated key name: %s", keyPath))
		return keyPath, nil
	}

	if config.FallbackKey != "" {
		Log(fmt.Sprintf("Using fallback key: %s", config.FallbackKey))
		return config.FallbackKey, nil
	}

	return "", ErrNoKeyAvailable
}
