package config

import (
	"fmt"
	"os"

	"github.com/Driver-C/tryssh/pkg/utils"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from the default paths.
func LoadConfig() (*MainConfig, error) {
	return LoadConfigFromPath(DefaultConfigPath, DefaultKnownHostsPath)
}

// LoadConfigFromPath loads the configuration from the specified paths,
// generating a new config file if one does not exist.
func LoadConfigFromPath(configPath, knownHostsPath string) (*MainConfig, error) {
	c := new(MainConfig)

	if utils.CheckFileIsExist(configPath) {
		conf, readErr := os.ReadFile(configPath) //nolint:gosec // G304: path is from known config constant
		if readErr != nil {
			return nil, fmt.Errorf("configuration file load failed: %w", readErr)
		}
		if unmarshalErr := yaml.Unmarshal(conf, c); unmarshalErr != nil {
			return nil, fmt.Errorf("configuration file parsing failed: %w", unmarshalErr)
		}
		if err := decryptConfig(c); err != nil {
			return nil, fmt.Errorf("configuration decryption failed: %w", err)
		}
	} else {
		if genErr := generateConfig(configPath); genErr != nil {
			return nil, genErr
		}
	}

	if !utils.CheckFileIsExist(knownHostsPath) {
		if createErr := utils.CreateFile(knownHostsPath, 0600); createErr != nil {
			return nil, fmt.Errorf("the known_hosts file creation failed: %w", createErr)
		}
	}
	return c, nil
}

func generateConfig(configPath string) error {
	if err := utils.FileYamlMarshalAndWrite(configPath, &MainConfig{}); err != nil {
		return fmt.Errorf("failed to generate configuration file: %w", err)
	}
	return nil
}

// UpdateConfig writes the configuration to the default config path.
func UpdateConfig(conf *MainConfig) error {
	return UpdateConfigAtPath(DefaultConfigPath, conf)
}

// UpdateConfigAtPath writes the configuration to the specified config path.
func UpdateConfigAtPath(configPath string, conf *MainConfig) error {
	toSave, encErr := encryptConfigForSave(conf)
	if encErr != nil {
		return encErr
	}
	return utils.FileYamlMarshalAndWrite(configPath, toSave)
}

// decryptConfig decrypts all encrypted fields in the config using the master key.
// Only prompts for the master password if encrypted content is detected.
func decryptConfig(c *MainConfig) error {
	hasEncrypted := false
	for _, pwd := range c.Main.Passwords {
		if utils.IsEncrypted(pwd) {
			hasEncrypted = true
			break
		}
	}
	if !hasEncrypted {
		for _, s := range c.ServerLists {
			if utils.IsEncrypted(s.Password) {
				hasEncrypted = true
				break
			}
		}
	}
	if !hasEncrypted {
		return nil
	}

	key, err := utils.GetMasterKey()
	if err != nil || key == nil {
		return fmt.Errorf("encrypted fields found but no master key provided")
	}

	for i, pwd := range c.Main.Passwords {
		if utils.IsEncrypted(pwd) {
			decrypted, err := utils.Decrypt(pwd, key)
			if err != nil {
				return fmt.Errorf("failed to decrypt password[%d]: %w", i, err)
			}
			c.Main.Passwords[i] = decrypted
		}
	}

	for i := range c.ServerLists {
		if utils.IsEncrypted(c.ServerLists[i].Password) {
			decrypted, err := utils.Decrypt(c.ServerLists[i].Password, key)
			if err != nil {
				return fmt.Errorf("failed to decrypt server cache password[%d]: %w", i, err)
			}
			c.ServerLists[i].Password = decrypted
		}
	}
	return nil
}

// encryptConfigForSave creates a copy with encrypted passwords for saving to disk.
// Uses the cached master key only — does not prompt interactively.
func encryptConfigForSave(conf *MainConfig) (*MainConfig, error) {
	key, err := utils.GetCachedMasterKey()
	if err != nil || key == nil {
		// No master key — save as plaintext
		return conf, nil
	}

	cp := *conf
	cp.Main.Passwords = make([]string, len(conf.Main.Passwords))
		for i, pwd := range conf.Main.Passwords {
		if utils.IsEncrypted(pwd) {
			cp.Main.Passwords[i] = pwd
		} else {
			enc, err := utils.Encrypt(pwd, key)
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt password[%d]: %w", i, err)
			}
			cp.Main.Passwords[i] = enc
		}
	}

	cp.ServerLists = make([]ServerListConfig, len(conf.ServerLists))
	for i, s := range conf.ServerLists {
		cp.ServerLists[i] = s
		if !utils.IsEncrypted(s.Password) {
			enc, err := utils.Encrypt(s.Password, key)
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt server cache password[%d]: %w", i, err)
			}
			cp.ServerLists[i].Password = enc
		}
	}
	return &cp, nil
}
