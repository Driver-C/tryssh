// Package encrypt provides the encrypt subcommand for encrypting configuration passwords.
package encrypt

import (
	"fmt"
	"os"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// readPasswordFn reads a password without echo. Injectable for testing.
var readPasswordFn = defaultReadPassword

// fatalFn is called on fatal errors. Injectable for testing.
var fatalFn = defaultFatal

func defaultReadPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return pwd, err
}

func defaultFatal(args ...interface{}) {
	utils.Fatalln(args...)
}

// NewEncryptCommand creates and returns the cobra command for encrypting config passwords.
func NewEncryptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt passwords in the configuration file",
		Long:  "Encrypt all plaintext passwords in the configuration file using a master password. Encrypted passwords are prefixed with 'enc:'.",
		Run: func(_ *cobra.Command, _ []string) {
			runEncrypt()
		},
	}
	return cmd
}

func runEncrypt() {
	// Check if env var is already set — skip interactive prompt
	if os.Getenv("TRYSSH_MASTER_KEY") != "" {
		executeEncrypt()
		return
	}

	// Interactive: prompt for password twice
	pwd1, err := readPasswordFn("Enter master password: ")
	if err != nil {
		fatalFn("Failed to read password:", err)
		return
	}
	if len(pwd1) == 0 {
		fatalFn("Password cannot be empty")
		return
	}
	if len(pwd1) < 4 {
		fatalFn("Password must be at least 4 characters")
		return
	}

	pwd2, err := readPasswordFn("Confirm master password: ")
	if err != nil {
		fatalFn("Failed to read password:", err)
		return
	}

	if string(pwd1) != string(pwd2) {
		fatalFn("Passwords do not match")
		return
	}

	// Set env var so GetMasterKey picks it up without another prompt
	pwdStr := string(pwd1)
	for i := range pwd1 {
		pwd1[i] = 0
	}
	for i := range pwd2 {
		pwd2[i] = 0
	}
	if err := os.Setenv("TRYSSH_MASTER_KEY", pwdStr); err != nil {
		fatalFn("Failed to set master key:", err)
		return
	}

	executeEncrypt()
}

// executeEncrypt loads the config, encrypts plaintext passwords, and writes back.
func executeEncrypt() {
	configuration, err := config.LoadConfig()
	if err != nil {
		fatalFn(err)
		return
	}

	count := countPlaintextPasswords(configuration)
	if count == 0 {
		fmt.Println("No passwords to encrypt")
		utils.ClearMasterKey()
		return
	}

	if err := config.UpdateConfig(configuration); err != nil {
		fatalFn("Failed to write encrypted config:", err)
		return
	}

	utils.ClearMasterKey()
	fmt.Printf("Encrypted %d password(s) successfully\n", count)
	fmt.Println("Remember your master password — it cannot be recovered if lost")
}

// countPlaintextPasswords counts password fields that are non-empty and not yet encrypted.
func countPlaintextPasswords(c *config.MainConfig) int {
	count := 0
	for _, pwd := range c.Main.Passwords {
		if pwd != "" && !utils.IsEncrypted(pwd) {
			count++
		}
	}
	for _, s := range c.ServerLists {
		if s.Password != "" && !utils.IsEncrypted(s.Password) {
			count++
		}
	}
	return count
}
