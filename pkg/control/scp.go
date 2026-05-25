package control

import (
	"strings"
	"time"

	"github.com/Driver-C/tryssh/pkg/config"
	"github.com/Driver-C/tryssh/pkg/launcher"
	"github.com/Driver-C/tryssh/pkg/utils"
)

// ScpController manages SCP file copy operations using cached credentials or credential combinations.
type ScpController struct {
	source        string
	destination   string
	configuration *config.MainConfig
	cacheIsFound  bool
	cacheIndex    int
	destIP        string
	concurrency   int
	sshTimeout    time.Duration
	recursive     bool
}

// parseRemotePath parses a remote path string (host:path or [host]:path) and returns
// the host/alias part and the path part.
func parseRemotePath(s string) (host, path string, ok bool) {
	if strings.HasPrefix(s, "[") {
		closeBracket := strings.Index(s, "]")
		if closeBracket < 0 {
			return "", "", false
		}
		host = s[1:closeBracket]
		rest := s[closeBracket+1:]
		if rest == "" {
			return host, "", false
		}
		if rest[0] == ':' {
			if len(rest) == 1 {
				return host, "", false
			}
			return host, rest[1:], true
		}
		return host, rest, true
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) == 2 && parts[1] != "" {
		return parts[0], parts[1], true
	}
	return "", "", false
}

// formatRemotePath builds a "host:path" string, wrapping IPv6 addresses in brackets.
func formatRemotePath(host, path string) string {
	if strings.Contains(host, ":") {
		return "[" + host + "]:" + path
	}
	return host + ":" + path
}

// TryCopy attempts to copy files to/from the target server, first using cached
// credentials and then by trying all credential combinations.
func (cc *ScpController) TryCopy(user string, concurrency int, recursive bool, sshTimeout time.Duration) {
	cc.sshTimeout = sshTimeout
	cc.concurrency = concurrency
	cc.recursive = recursive

	if host, path, ok := parseRemotePath(cc.source); ok {
		cc.destIP = config.ResolveAlias(host, cc.configuration)
		cc.source = formatRemotePath(cc.destIP, path)
	} else if host, path, ok := parseRemotePath(cc.destination); ok {
		cc.destIP = config.ResolveAlias(host, cc.configuration)
		cc.destination = formatRemotePath(cc.destIP, path)
	} else {
		utils.Errorln("Unable to determine SCP direction: no valid remote path found in source or destination")
		return
	}
	var targetServer *config.ServerListConfig
	targetServer, cc.cacheIndex, cc.cacheIsFound = config.SelectServerCache(user, cc.destIP, cc.configuration)

	if cc.cacheIsFound {
		utils.Infof("The cache for %s is found, which will be used to try.\n", cc.destIP)
		cc.tryCopyWithCache(user, targetServer)
	} else {
		utils.Warnf("The cache for %s could not be found. Start trying to login.\n\n", cc.destIP)
		cc.tryCopyWithoutCache(user)
	}
}

func (cc *ScpController) tryCopyWithCache(user string, targetServer *config.ServerListConfig) {
	lan := &launcher.ScpLauncher{
		SSHConnector: *launcher.GetSSHConnectorFromConfig(targetServer),
		Src:          cc.source,
		Dest:         cc.destination,
		Recursive:    cc.recursive,
	}
	lan.SSHTimeout = sshClientTimeoutWhenLogin
	if !lan.Launch() {
		utils.Errorf("Failed to log in with cached information. Start trying to login again.\n\n")
		cc.tryCopyWithoutCache(user)
	}
}

func (cc *ScpController) tryCopyWithoutCache(user string) {
	combinations := config.GenerateCombination(cc.destIP, user, cc.configuration)
	launchers := launcher.NewScpLaunchersByCombinations(combinations, cc.source, cc.destination,
		cc.recursive, cc.sshTimeout)
	connectors := make([]launcher.Connector, len(launchers))
	for i, l := range launchers {
		connectors[i] = l
	}
	hitLaunchers := ConcurrencyTryToConnect(cc.concurrency, connectors)
	if len(hitLaunchers) > 0 {
		utils.Infoln("Login succeeded. The cache will be added.")
		hitLauncher := hitLaunchers[0].(*launcher.ScpLauncher)
		newServerCache := launcher.GetConfigFromSSHConnector(&hitLauncher.SSHConnector)
		if cc.cacheIsFound {
			newServerCache.Alias = cc.configuration.ServerLists[cc.cacheIndex].Alias
			utils.Infoln("The old cache will be deleted.")
			cc.configuration.ServerLists = append(
				cc.configuration.ServerLists[:cc.cacheIndex], cc.configuration.ServerLists[cc.cacheIndex+1:]...)
		}
		cc.configuration.ServerLists = append(cc.configuration.ServerLists, *newServerCache)
		if err := config.UpdateConfig(cc.configuration); err == nil {
			utils.Infoln("Cache added.")
			if cc.sshTimeout > sshClientTimeoutWhenLogin {
				hitLauncher.SSHTimeout = cc.sshTimeout
			} else {
				hitLauncher.SSHTimeout = sshClientTimeoutWhenLogin
			}
			if !hitLauncher.Launch() {
				utils.Errorf("Login failed.\n")
			}
		} else {
			utils.Errorf("Cache added failed.\n\n")
		}
	} else {
		utils.Errorf("There is no password combination that can log in.\n")
	}
}

// NewScpController creates a new ScpController for the given source, destination, and configuration.
func NewScpController(source string, destination string, configuration *config.MainConfig) *ScpController {
	return &ScpController{
		source:        source,
		destination:   destination,
		configuration: configuration,
	}
}
