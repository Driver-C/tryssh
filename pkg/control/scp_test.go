package control

import (
	"testing"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewScpController(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("/tmp/file.txt", "192.168.1.1:/tmp/dest.txt", cfg)
	assert.NotNil(t, ctrl)
	assert.Equal(t, "/tmp/file.txt", ctrl.source)
	assert.Equal(t, "192.168.1.1:/tmp/dest.txt", ctrl.destination)
	assert.Equal(t, cfg, ctrl.configuration)
}

func TestScpController_TryCopy_SourceContainsColon(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("server1:/remote/path/file.txt", "/local/path/", cfg)

	// Verify controller fields are set
	assert.Equal(t, "server1:/remote/path/file.txt", ctrl.source)
	assert.Equal(t, "/local/path/", ctrl.destination)

	// Verify alias resolution
	resolved := config.ResolveAlias("server1", cfg)
	assert.Equal(t, "192.168.1.1", resolved)
}

func TestScpController_TryCopy_DestContainsColon(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("/local/path/file.txt", "192.168.1.1:/remote/path/", cfg)

	assert.Equal(t, "/local/path/file.txt", ctrl.source)
	assert.Equal(t, "192.168.1.1:/remote/path/", ctrl.destination)
}

func TestScpController_TryCopy_NoColonAnywhere(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("/local/file.txt", "/local/dest.txt", cfg)

	// Neither source nor dest contains ":", TryCopy should return early.
	assert.Equal(t, "/local/file.txt", ctrl.source)
	assert.Equal(t, "/local/dest.txt", ctrl.destination)
	assert.Equal(t, false, ctrl.cacheIsFound)
}

func TestScpController_TryCopy_AliasInSource(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("server2:/remote/file.txt", "/local/dest.txt", cfg)

	// Verify controller created correctly
	assert.Equal(t, "server2:/remote/file.txt", ctrl.source)

	// Verify alias resolution for source
	resolved := config.ResolveAlias("server2", cfg)
	assert.Equal(t, "192.168.1.2", resolved)
}

func TestScpController_TryCopy_AliasInDest(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("/local/file.txt", "server1:/remote/path/", cfg)

	// Verify controller created correctly
	assert.Equal(t, "server1:/remote/path/", ctrl.destination)

	// Verify alias resolution for destination
	resolved := config.ResolveAlias("server1", cfg)
	assert.Equal(t, "192.168.1.1", resolved)
}

func TestScpController_Fields(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("src", "dst", cfg)

	assert.Equal(t, "", ctrl.destIP)
	assert.Equal(t, false, ctrl.cacheIsFound)
	assert.Equal(t, 0, ctrl.cacheIndex)
	assert.Equal(t, 0, ctrl.concurrency)
}

func TestScpController_TryCopy_IpInSource(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("192.168.1.1:/remote/file.txt", "/local/dest.txt", cfg)

	assert.Contains(t, ctrl.source, "192.168.1.1")

	server, idx, found := config.SelectServerCache("", "192.168.1.1", cfg)
	assert.True(t, found)
	assert.NotNil(t, server)
	assert.Equal(t, 0, idx)
}

func TestScpController_TryCopy_IpInDest(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("/local/file.txt", "192.168.1.2:/remote/path/", cfg)

	assert.Contains(t, ctrl.destination, "192.168.1.2")

	server, idx, found := config.SelectServerCache("", "192.168.1.2", cfg)
	assert.True(t, found)
	assert.NotNil(t, server)
	assert.Equal(t, 1, idx)
}

func TestScpController_TryCopy_UnknownIp(t *testing.T) {
	cfg := newTestMainConfig()
	ctrl := NewScpController("/local/file.txt", "10.0.0.1:/remote/path/", cfg)

	assert.Contains(t, ctrl.destination, "10.0.0.1")

	server, _, found := config.SelectServerCache("", "10.0.0.1", cfg)
	assert.False(t, found)
	assert.Nil(t, server)
}
