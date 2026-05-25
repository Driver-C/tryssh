package create

import (
	"fmt"
	"os"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewPasswordsCommand creates and returns the cobra command for creating a password entry.
func NewPasswordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "passwords",
		Args:    cobra.NoArgs,
		Short:   "Create an alternative password",
		Long:    "Create an alternative password (interactive prompt)",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Print("Enter password: ")
			pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				utils.Fatalln("Failed to read password:", err)
			}
			password := string(pwdBytes)
			for i := range pwdBytes {
				pwdBytes[i] = 0
			}
			if password == "" {
				utils.Fatalln("Password cannot be empty")
			}

			fmt.Print("Confirm password: ")
			confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				utils.Fatalln("Failed to read password:", err)
			}
			match := string(confirmBytes) == password
			for i := range confirmBytes {
				confirmBytes[i] = 0
			}
			if !match {
				utils.Fatalln("Passwords do not match")
			}

			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewCreateController(control.TypePasswords, password, configuration)
			controller.ExecuteCreate()
		},
	}
	return cmd
}
