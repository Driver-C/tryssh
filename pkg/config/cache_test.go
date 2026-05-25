package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestConfig() *MainConfig {
	conf := &MainConfig{}
	conf.Main.Ports = []string{"22"}
	conf.Main.Users = []string{"root"}
	conf.Main.Passwords = []string{"pass1"}
	conf.Main.Keys = []string{"/key1"}
	conf.ServerLists = []ServerListConfig{
		{
			IP:       "192.168.1.1",
			Port:     "22",
			User:     "root",
			Password: "pass1",
			Key:      "/key1",
			Alias:    "server1",
		},
		{
			IP:       "192.168.1.2",
			Port:     "2222",
			User:     "admin",
			Password: "pass2",
			Key:      "/key2",
			Alias:    "server2",
		},
		{
			IP:       "192.168.1.1",
			Port:     "22",
			User:     "admin",
			Password: "pass3",
			Key:      "/key3",
			Alias:    "server1-admin",
		},
		{
			IP:       "10.0.0.1",
			Port:     "22",
			User:     "deploy",
			Password: "deploy-pass",
			Alias:    "",
		},
	}
	return conf
}

func TestSelectServerCache_MatchingIP(t *testing.T) {
	conf := newTestConfig()

	// Search with IP only (empty user) -- should return the first matching IP
	server, index, found := SelectServerCache("", "192.168.1.1", conf)
	assert.True(t, found)
	assert.NotNil(t, server)
	assert.Equal(t, 0, index)
	assert.Equal(t, "192.168.1.1", server.IP)
	assert.Equal(t, "root", server.User)
}

func TestSelectServerCache_MatchingIPAndUser(t *testing.T) {
	conf := newTestConfig()

	// Search with both IP and user -- should match the specific server
	server, index, found := SelectServerCache("admin", "192.168.1.1", conf)
	assert.True(t, found)
	assert.NotNil(t, server)
	assert.Equal(t, 2, index)
	assert.Equal(t, "192.168.1.1", server.IP)
	assert.Equal(t, "admin", server.User)
	assert.Equal(t, "server1-admin", server.Alias)
}

func TestSelectServerCache_MatchingIPAndUserSecondServer(t *testing.T) {
	conf := newTestConfig()

	server, index, found := SelectServerCache("admin", "192.168.1.2", conf)
	assert.True(t, found)
	assert.NotNil(t, server)
	assert.Equal(t, 1, index)
	assert.Equal(t, "192.168.1.2", server.IP)
	assert.Equal(t, "admin", server.User)
	assert.Equal(t, "server2", server.Alias)
}

func TestSelectServerCache_NoMatch(t *testing.T) {
	conf := newTestConfig()

	// IP that doesn't exist
	server, index, found := SelectServerCache("", "10.10.10.10", conf)
	assert.False(t, found)
	assert.Nil(t, server)
	assert.Equal(t, 0, index)

	// User that doesn't match the IP entries
	server, index, found = SelectServerCache("nonexistent", "192.168.1.1", conf)
	assert.False(t, found)
	assert.Nil(t, server)
	assert.Equal(t, 0, index)
}

func TestSelectServerCache_IPMatchesButUserDoesNot(t *testing.T) {
	conf := newTestConfig()

	// 10.0.0.1 exists with user "deploy", but we search for "root"
	server, index, found := SelectServerCache("root", "10.0.0.1", conf)
	assert.False(t, found)
	assert.Nil(t, server)
	assert.Equal(t, 0, index)
}

func TestSelectServerCache_EmptyServerList(t *testing.T) {
	conf := &MainConfig{}
	conf.ServerLists = nil

	server, index, found := SelectServerCache("", "192.168.1.1", conf)
	assert.False(t, found)
	assert.Nil(t, server)
	assert.Equal(t, 0, index)
}

func TestResolveAlias_ExistingAlias(t *testing.T) {
	conf := newTestConfig()

	result := ResolveAlias("server1", conf)
	assert.Equal(t, "192.168.1.1", result)

	result = ResolveAlias("server2", conf)
	assert.Equal(t, "192.168.1.2", result)

	result = ResolveAlias("server1-admin", conf)
	assert.Equal(t, "192.168.1.1", result)
}

func TestResolveAlias_NonExistingAlias(t *testing.T) {
	conf := newTestConfig()

	// When alias is not found, it returns the alias itself
	result := ResolveAlias("unknown-server", conf)
	assert.Equal(t, "unknown-server", result)
}

func TestResolveAlias_EmptyAlias(t *testing.T) {
	conf := newTestConfig()

	// Empty alias matches the server with Alias="" (10.0.0.1), returning its IP
	result := ResolveAlias("", conf)
	assert.Equal(t, "10.0.0.1", result)
}

func TestResolveAlias_EmptyServerList(t *testing.T) {
	conf := &MainConfig{}
	conf.ServerLists = nil

	result := ResolveAlias("server1", conf)
	assert.Equal(t, "server1", result)
}

func TestFindAlias_ExistingAlias(t *testing.T) {
	conf := newTestConfig()

	results := FindAlias("server1", conf)
	assert.Len(t, results, 1)
	assert.Equal(t, "server1", results[0].Alias)
	assert.Equal(t, "192.168.1.1", results[0].IP)
	assert.Equal(t, "root", results[0].User)
}

func TestFindAlias_NonExistingAlias(t *testing.T) {
	conf := newTestConfig()

	results := FindAlias("nonexistent", conf)
	assert.Empty(t, results)
}

func TestFindAlias_EmptyAlias(t *testing.T) {
	conf := newTestConfig()

	// Empty alias should not match (the function checks alias != "")
	results := FindAlias("", conf)
	assert.Empty(t, results)
}

func TestFindAlias_MultipleMatches(t *testing.T) {
	conf := newTestConfig()
	// Add another server with the same alias "server1"
	conf.ServerLists = append(conf.ServerLists, ServerListConfig{
		IP:       "192.168.1.100",
		Port:     "22",
		User:     "deploy",
		Password: "deploy-pass",
		Key:      "/key-deploy",
		Alias:    "server1",
	})

	results := FindAlias("server1", conf)
	assert.Len(t, results, 2)

	ips := map[string]bool{}
	for _, s := range results {
		assert.Equal(t, "server1", s.Alias)
		ips[s.IP] = true
	}
	assert.True(t, ips["192.168.1.1"])
	assert.True(t, ips["192.168.1.100"])
}

func TestFindAlias_EmptyServerList(t *testing.T) {
	conf := &MainConfig{}
	conf.ServerLists = nil

	results := FindAlias("server1", conf)
	assert.Empty(t, results)
}

func TestFindAlias_ServerWithEmptyAlias(t *testing.T) {
	conf := newTestConfig()

	// "deploy" has an empty alias in our test config, searching for it should return nothing
	// but searching for a non-empty alias should still work
	results := FindAlias("", conf)
	assert.Empty(t, results)
}
