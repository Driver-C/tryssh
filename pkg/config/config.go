package config

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/Driver-C/tryssh/pkg/utils"
)

// ConfigFileName is the default name of the configuration database file.
const (
	ConfigFileName     = "tryssh.db"
	ConfigDirName      = ".tryssh"
	KnownHostsFileName = "known_hosts"
)

// DefaultConfigPath is the absolute path to the default configuration file.
var (
	DefaultConfigPath     string
// DefaultKnownHostsPath is the absolute path to the default known_hosts file.
	DefaultKnownHostsPath string
)

func init() {
	DefaultConfigPath, DefaultKnownHostsPath = DefaultPaths()
}

// DefaultPaths returns the default configuration file path and known_hosts file path
// based on the current user's home directory.
func DefaultPaths() (configPath, knownHostsPath string) {
	usr, err := user.Current()
	if err != nil {
		configPath = filepath.Join("./", ConfigDirName, ConfigFileName)
		knownHostsPath = filepath.Join("./", ConfigDirName, KnownHostsFileName)
		return
	}
	configPath = filepath.Join(usr.HomeDir, ConfigDirName, ConfigFileName)
	knownHostsPath = filepath.Join(usr.HomeDir, ConfigDirName, KnownHostsFileName)
	return
}

// MainConfig represents the top-level configuration for the tryssh application,
// including credential lists and cached server entries.
type MainConfig struct {
	Main struct {
		Ports     []string `yaml:"ports,flow"`
		Users     []string `yaml:"users,flow"`
		Passwords []string `yaml:"passwords,flow"`
		Keys      []string `yaml:"keys,flow"`
	} `yaml:"main"`
	ServerLists []ServerListConfig `yaml:"serverList"`
}

// ServerListConfig holds the connection details for a cached server entry.
type ServerListConfig struct {
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Key      string `yaml:"key"`
	Alias    string `yaml:"alias"`
}

// String returns a safe string representation with the password masked.
func (s ServerListConfig) String() string {
	pwd := utils.MaskSecret(s.Password)
	key := utils.MaskSecret(s.Key)
	return fmt.Sprintf("%s@%s:%s (pwd:%s key:%s alias:%s)",
		s.User, s.IP, s.Port, pwd, key, s.Alias)
}