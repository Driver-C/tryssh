package control

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRemotePath_PlainHostPath(t *testing.T) {
	host, path, ok := parseRemotePath("192.168.1.1:/tmp/file.txt")
	assert.True(t, ok)
	assert.Equal(t, "192.168.1.1", host)
	assert.Equal(t, "/tmp/file.txt", path)
}

func TestParseRemotePath_PlainHostPathWithPort(t *testing.T) {
	host, path, ok := parseRemotePath("example.com:22:/tmp/file.txt")
	// SplitN with 2 splits on the first colon only
	assert.True(t, ok)
	assert.Equal(t, "example.com", host)
	assert.Equal(t, "22:/tmp/file.txt", path)
}

func TestParseRemotePath_Alias(t *testing.T) {
	host, path, ok := parseRemotePath("myserver:/remote/dir")
	assert.True(t, ok)
	assert.Equal(t, "myserver", host)
	assert.Equal(t, "/remote/dir", path)
}

func TestParseRemotePath_IPv6(t *testing.T) {
	host, path, ok := parseRemotePath("[::1]:/tmp/file.txt")
	assert.True(t, ok)
	assert.Equal(t, "::1", host)
	assert.Equal(t, "/tmp/file.txt", path)
}

func TestParseRemotePath_IPv6Full(t *testing.T) {
	host, path, ok := parseRemotePath("[fe80::1%eth0]:/home/user/data")
	assert.True(t, ok)
	assert.Equal(t, "fe80::1%eth0", host)
	assert.Equal(t, "/home/user/data", path)
}

func TestParseRemotePath_IPv6Long(t *testing.T) {
	host, path, ok := parseRemotePath("[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:/data")
	assert.True(t, ok)
	assert.Equal(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", host)
	assert.Equal(t, "/data", path)
}

func TestParseRemotePath_IPv6NoColonAfterBracket(t *testing.T) {
	// [host] with no colon after bracket -- rest is empty, host returned but ok=false
	host, path, ok := parseRemotePath("[::1]")
	assert.False(t, ok)
	assert.Equal(t, "::1", host)
	assert.Equal(t, "", path)
}

func TestParseRemotePath_IPv6ColonOnlyNoPath(t *testing.T) {
	// [host]: with no path after colon
	host, path, ok := parseRemotePath("[::1]:")
	assert.False(t, ok)
	assert.Equal(t, "::1", host)
	assert.Equal(t, "", path)
}

func TestParseRemotePath_IPv6UnclosedBracket(t *testing.T) {
	host, path, ok := parseRemotePath("[::1/tmp/file.txt")
	assert.False(t, ok)
	assert.Equal(t, "", host)
	assert.Equal(t, "", path)
}

func TestParseRemotePath_EmptyString(t *testing.T) {
	host, path, ok := parseRemotePath("")
	assert.False(t, ok)
	assert.Equal(t, "", host)
	assert.Equal(t, "", path)
}

func TestParseRemotePath_NoColon(t *testing.T) {
	host, path, ok := parseRemotePath("just-a-hostname")
	assert.False(t, ok)
	assert.Equal(t, "", host)
	assert.Equal(t, "", path)
}

func TestParseRemotePath_ColonButEmptyPath(t *testing.T) {
	host, path, ok := parseRemotePath("192.168.1.1:")
	assert.False(t, ok)
	assert.Equal(t, "", host)
	assert.Equal(t, "", path)
}

func TestParseRemotePath_PathWithMultipleColons(t *testing.T) {
	host, path, ok := parseRemotePath("host:path:with:colons")
	// SplitN(..., 2) splits only on the first colon
	assert.True(t, ok)
	assert.Equal(t, "host", host)
	assert.Equal(t, "path:with:colons", path)
}

func TestParseRemotePath_IPv6PathWithMultipleColons(t *testing.T) {
	host, path, ok := parseRemotePath("[::1]:/path:weird:name")
	assert.True(t, ok)
	assert.Equal(t, "::1", host)
	assert.Equal(t, "/path:weird:name", path)
}

func TestParseRemotePath_IPv6BracketRestWithoutColon(t *testing.T) {
	// When [host] is followed by something that is not a colon, it still returns ok=true
	// because the code returns host, rest, true for non-colon rest.
	host, path, ok := parseRemotePath("[::1]extra")
	assert.True(t, ok)
	assert.Equal(t, "::1", host)
	assert.Equal(t, "extra", path)
}

func TestParseRemotePath_SimplePath(t *testing.T) {
	host, path, ok := parseRemotePath("10.0.0.1:/home/user/.ssh/config")
	assert.True(t, ok)
	assert.Equal(t, "10.0.0.1", host)
	assert.Equal(t, "/home/user/.ssh/config", path)
}

func TestFormatRemotePath_PlainHost(t *testing.T) {
	result := formatRemotePath("192.168.1.1", "/tmp/file.txt")
	assert.Equal(t, "192.168.1.1:/tmp/file.txt", result)
}

func TestFormatRemotePath_IPv6(t *testing.T) {
	result := formatRemotePath("::1", "/tmp/file.txt")
	assert.Equal(t, "[::1]:/tmp/file.txt", result)
}

func TestFormatRemotePath_IPv6Full(t *testing.T) {
	result := formatRemotePath("fe80::1%eth0", "/home/user/data")
	assert.Equal(t, "[fe80::1%eth0]:/home/user/data", result)
}

func TestFormatRemotePath_LongIPv6(t *testing.T) {
	result := formatRemotePath("2001:0db8:85a3::8a2e:0370:7334", "/data")
	assert.Equal(t, "[2001:0db8:85a3::8a2e:0370:7334]:/data", result)
}

func TestFormatRemotePath_DomainName(t *testing.T) {
	result := formatRemotePath("example.com", "/remote/path")
	assert.Equal(t, "example.com:/remote/path", result)
}

func TestFormatRemotePath_EmptyHost(t *testing.T) {
	result := formatRemotePath("", "/path")
	assert.Equal(t, ":/path", result)
}

func TestFormatRemotePath_EmptyPath(t *testing.T) {
	result := formatRemotePath("host", "")
	assert.Equal(t, "host:", result)
}

// Table-driven test for parseRemotePath
func TestParseRemotePath_Table(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantHost  string
		wantPath  string
		wantOk    bool
	}{
		{"plain IPv4 with path", "192.168.1.1:/tmp/file", "192.168.1.1", "/tmp/file", true},
		{"plain host with path", "myhost:/data", "myhost", "/data", true},
		{"IPv6 with path", "[::1]:/tmp/file", "::1", "/tmp/file", true},
		{"IPv6 full with path", "[fe80::1]:/home", "fe80::1", "/home", true},
		{"empty string", "", "", "", false},
		{"no colon", "justhost", "", "", false},
		{"colon empty path", "host:", "", "", false},
		{"IPv6 no close bracket", "[::1path", "", "", false},
		{"IPv6 empty rest", "[::1]", "::1", "", false},
		{"IPv6 colon only", "[::1]:", "::1", "", false},
		{"path with colons", "host:a:b:c", "host", "a:b:c", true},
		{"IPv6 bracket rest no colon", "[::1]extra", "::1", "extra", true},
		{"localhost with path", "localhost:/var/log", "localhost", "/var/log", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, path, ok := parseRemotePath(tt.input)
			assert.Equal(t, tt.wantHost, host)
			assert.Equal(t, tt.wantPath, path)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

// Table-driven test for formatRemotePath
func TestFormatRemotePath_Table(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		path     string
		expected string
	}{
		{"plain IPv4", "10.0.0.1", "/tmp", "10.0.0.1:/tmp"},
		{"IPv6 loopback", "::1", "/tmp", "[::1]:/tmp"},
		{"IPv6 full", "2001:db8::1", "/data", "[2001:db8::1]:/data"},
		{"domain", "example.com", "/file", "example.com:/file"},
		{"empty host", "", "/path", ":/path"},
		{"empty path", "host", "", "host:"},
		{"host with zone", "fe80::1%en0", "/home", "[fe80::1%en0]:/home"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRemotePath(tt.host, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}
