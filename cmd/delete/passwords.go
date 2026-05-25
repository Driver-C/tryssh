package delete

import (
	"fmt"
	"os"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/control"
	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewPasswordsCommand creates and returns the cobra command for deleting a password entry.
func NewPasswordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "passwords",
		Args:    cobra.NoArgs,
		Short:   "Delete an alternative password",
		Long:    "Delete an alternative password (interactive prompt)",
		Aliases: []string{"password", "pass", "pwd"},
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Print("Enter password to delete: ")
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

			configuration, err := config.LoadConfig()
			if err != nil {
				utils.Fatalln(err)
			}
			controller := control.NewDeleteController(control.TypePasswords, password, configuration)
			controller.ExecuteDelete()
		},
	}
	return cmd
}
