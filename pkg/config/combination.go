package config

import (
	"fmt"

	"github.com/Driver-C/tryssh/pkg/utils"
	"github.com/schwarmco/go-cartesian-product"
)

// GenerateCombination produces a channel of credential combinations for the given
// IP address and optional username using the main configuration values.
// Passwords and keys are treated as alternatives: if only one is configured,
// the other is padded with an empty string so the cartesian product still produces results.
func GenerateCombination(ip string, user string, conf *MainConfig) chan []interface{} {
	ips := utils.ToInterfaceSlice([]string{ip})
	users := utils.ToInterfaceSlice([]string{user})
	ports := utils.ToInterfaceSlice(conf.Main.Ports)
	if user == "" {
		users = utils.ToInterfaceSlice(conf.Main.Users)
	}
	passwords := utils.ToInterfaceSlice(conf.Main.Passwords)
	keys := utils.ToInterfaceSlice(conf.Main.Keys)

	if len(passwords) == 0 && len(keys) == 0 {
		utils.Warnln("No passwords or keys configured — no credential combinations can be generated.")
		fmt.Println("Hint: Use 'tryssh create passwords' or 'tryssh create keys' to add credentials.")
	}

	// Passwords and keys are alternatives, not jointly required.
	// Pad with empty string so the cartesian product produces results when only one is configured.
	if len(passwords) == 0 {
		passwords = utils.ToInterfaceSlice([]string{""})
	}
	if len(keys) == 0 {
		keys = utils.ToInterfaceSlice([]string{""})
	}

	return cartesian.Iter(ips, ports, users, passwords, keys)
}
