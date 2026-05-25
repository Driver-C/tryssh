package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_InvalidCommand(t *testing.T) {
	// Override os.Args to run with an invalid subcommand
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "__nonexistent__"}
	result := run()
	assert.Equal(t, 1, result, "invalid command should return 1")
}

func TestRun_VersionCommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "version"}
	result := run()
	assert.Equal(t, 0, result, "version command should return 0")
}

func TestRun_NoArgs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Running with no subcommand shows help and returns 0 (cobra default)
	os.Args = []string{"tryssh"}
	result := run()
	// Root command without subcommand: cobra shows help and returns 0 by default
	assert.Equal(t, 0, result)
}

func TestRun_HelpFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "--help"}
	result := run()
	assert.Equal(t, 0, result, "help flag should return 0")
}

func TestRun_SSHMissingArgs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "ssh"}
	result := run()
	assert.Equal(t, 1, result, "ssh without args should return 1")
}

func TestRun_SCPMissingArgs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "scp"}
	result := run()
	assert.Equal(t, 1, result, "scp without args should return 1")
}

func TestRun_AliasMissingSubcommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "alias"}
	result := run()
	// Cobra shows help for parent commands without subcommands and returns 0
	assert.Equal(t, 0, result, "alias without subcommand shows help and returns 0")
}

func TestRun_CreateMissingSubcommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "create"}
	result := run()
	assert.Equal(t, 0, result, "create without subcommand shows help and returns 0")
}

func TestRun_DeleteMissingSubcommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "delete"}
	result := run()
	assert.Equal(t, 0, result, "delete without subcommand shows help and returns 0")
}

func TestRun_GetMissingSubcommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "get"}
	result := run()
	assert.Equal(t, 0, result, "get without subcommand shows help and returns 0")
}

func TestRun_SSHHelpFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "ssh", "--help"}
	result := run()
	assert.Equal(t, 0, result, "ssh --help should return 0")
}

func TestRun_SCPHelpFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "scp", "--help"}
	result := run()
	assert.Equal(t, 0, result, "scp --help should return 0")
}

func TestRun_PruneHelpFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"tryssh", "prune", "--help"}
	result := run()
	assert.Equal(t, 0, result, "prune --help should return 0")
}
