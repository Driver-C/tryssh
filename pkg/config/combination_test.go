package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCombination_FullCredentials(t *testing.T) {
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22", "2222"}
	conf.Main.Users = []string{"root", "admin"}
	conf.Main.Passwords = []string{"pass1", "pass2"}
	conf.Main.Keys = []string{"/key1", "/key2"}

	ip := "192.168.1.1"
	user := "root"

	ch := GenerateCombination(ip, user, conf)

	var results [][]interface{}
	for combo := range ch {
		results = append(results, combo)
	}

	// Expected: 1 ip x 2 ports x 1 user (specified) x 2 passwords x 2 keys = 8
	assert.Len(t, results, 8)

	// Verify each combination has 5 elements: ip, port, user, password, key
	for _, combo := range results {
		assert.Len(t, combo, 5)
		assert.Equal(t, ip, combo[0])
		assert.Equal(t, user, combo[2])
	}
}

func TestGenerateCombination_SpecifiedUser(t *testing.T) {
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Users = []string{"root", "admin", "deploy"}
	conf.Main.Passwords = []string{"pass1"}
	conf.Main.Keys = []string{"/key1"}

	ip := "10.0.0.1"
	user := "deploy"

	ch := GenerateCombination(ip, user, conf)

	var results [][]interface{}
	for combo := range ch {
		results = append(results, combo)
	}

	// When user is specified, only that single user is used
	// Expected: 1 ip x 1 port x 1 user (specified) x 1 password x 1 key = 1
	assert.Len(t, results, 1)
	assert.Equal(t, ip, results[0][0])
	assert.Equal(t, "22", results[0][1])
	assert.Equal(t, user, results[0][2])
	assert.Equal(t, "pass1", results[0][3])
	assert.Equal(t, "/key1", results[0][4])
}

func TestGenerateCombination_EmptyUser_UsesAllUsers(t *testing.T) {
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Users = []string{"root", "admin"}
	conf.Main.Passwords = []string{"pass1"}
	conf.Main.Keys = []string{"/key1"}

	ip := "10.0.0.1"
	user := ""

	ch := GenerateCombination(ip, user, conf)

	var results [][]interface{}
	for combo := range ch {
		results = append(results, combo)
	}

	// When user is empty, all users from config should be used
	// Expected: 1 ip x 1 port x 2 users x 1 password x 1 key = 2
	assert.Len(t, results, 2)

	usersSeen := map[string]bool{}
	for _, combo := range results {
		assert.Equal(t, ip, combo[0])
		usersSeen[combo[2].(string)] = true
	}
	assert.True(t, usersSeen["root"])
	assert.True(t, usersSeen["admin"])
}

func TestGenerateCombination_EmptyCredentials(t *testing.T) {
	conf := &MainConfig{}
	conf.Main.Ports = nil
	conf.Main.Users = nil
	conf.Main.Passwords = nil
	conf.Main.Keys = nil

	ip := "192.168.1.1"
	user := ""

	// When all slices are nil/empty, the cartesian product with nil slices
	// should still produce results (the cartesian library handles nil as empty)
	ch := GenerateCombination(ip, user, conf)

	var results [][]interface{}
	for combo := range ch {
		results = append(results, combo)
	}

	// With nil slices passed to InterfaceSlice which returns nil,
	// cartesian.Iter with nil inputs yields 0 combinations
	assert.Len(t, results, 0)
}

func TestGenerateCombination_CartesianProductCorrectness(t *testing.T) {
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22", "2222"}
	conf.Main.Users = []string{"root"}
	conf.Main.Passwords = []string{"pass1", "pass2", "pass3"}
	conf.Main.Keys = []string{"/key1"}

	ip := "192.168.1.100"
	user := ""

	ch := GenerateCombination(ip, user, conf)

	var results [][]interface{}
	for combo := range ch {
		results = append(results, combo)
	}

	// Expected: 1 ip x 2 ports x 1 user x 3 passwords x 1 key = 6
	assert.Len(t, results, 6)

	// Verify all combinations are unique (cartesian product correctness)
	seen := map[string]bool{}
	for _, combo := range results {
		key := combo[0].(string) + ":" + combo[1].(string) + ":" +
			combo[2].(string) + ":" + combo[3].(string) + ":" + combo[4].(string)
		assert.False(t, seen[key], "duplicate combination found: %s", key)
		seen[key] = true
	}
	assert.Len(t, seen, 6)

	// Verify ports are both present
	ports := map[string]bool{}
	passwords := map[string]bool{}
	for _, combo := range results {
		ports[combo[1].(string)] = true
		passwords[combo[3].(string)] = true
	}
	assert.True(t, ports["22"])
	assert.True(t, ports["2222"])
	assert.True(t, passwords["pass1"])
	assert.True(t, passwords["pass2"])
	assert.True(t, passwords["pass3"])
}

func TestGenerateCombination_SinglePortMultipleKeys(t *testing.T) {
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Users = []string{"root"}
	conf.Main.Passwords = []string{"pass1"}
	conf.Main.Keys = []string{"/key1", "/key2", "/key3"}

	ip := "10.0.0.1"
	user := "root"

	ch := GenerateCombination(ip, user, conf)

	var results [][]interface{}
	for combo := range ch {
		results = append(results, combo)
	}

	// Expected: 1 ip x 1 port x 1 user x 1 password x 3 keys = 3
	assert.Len(t, results, 3)

	keysSeen := map[string]bool{}
	for _, combo := range results {
		keysSeen[combo[4].(string)] = true
	}
	assert.True(t, keysSeen["/key1"])
	assert.True(t, keysSeen["/key2"])
	assert.True(t, keysSeen["/key3"])
}
