package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/spf13/viper"
)

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
}

var config Config
var debug bool

func log(message string) {
	if debug {
		fmt.Fprintln(os.Stderr, message)
	}
}

// LoadConfig loads & populates the config for this tool.
func LoadConfig() error {
	viper.SetConfigName("deploy-keychain")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.deploy-keychain")
	viper.AddConfigPath(".")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user's home directory: %w", err)
	}

	viper.SetDefault("key_path", path.Join(homeDir, ".ssh", "deploy-keys"))
	viper.SetDefault("key_name_format", "{{.account}}-{{.repository}}.pem")
	viper.SetDefault("keys", make(map[string]string))
	viper.SetDefault("ssh_command", "ssh")

	err = viper.ReadInConfig()
	if err != nil {
		log(fmt.Sprintf("Unable to load config file: %s\n", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	log(fmt.Sprintf("Config loaded: %+v", config))

	return nil
}

func main() {
	debugVar := os.Getenv("DEPLOY_KEYCHAIN_DEBUG")
	if debugVar == "" {
		debugVar = "false"
	}
	debug, _ = strconv.ParseBool(debugVar)

	LoadConfig()

	log(fmt.Sprintf("Args: %s", os.Args))
}
