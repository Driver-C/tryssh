// Package config handles loading, saving, and managing SSH configuration.
package config

// SelectServerCache searches the server list for a cached entry matching the given
// IP and optional user. It returns the matching config, its index, and whether a
// match was found.
func SelectServerCache(user string, ip string, conf *MainConfig) (*ServerListConfig, int, bool) {
	for index, server := range conf.ServerLists {
		if server.IP == ip {
			if user != "" {
				if server.User == user {
					return &conf.ServerLists[index], index, true
				}
			} else {
				return &conf.ServerLists[index], index, true
			}
		}
	}
	return nil, 0, false
}

// ResolveAlias looks up the given alias in the server list and returns the
// corresponding IP address. If no match is found, the alias string itself is returned.
func ResolveAlias(alias string, conf *MainConfig) string {
	for _, server := range conf.ServerLists {
		if server.Alias == alias {
			return server.IP
		}
	}
	return alias
}

// FindAlias returns all server list entries that have the specified alias.
func FindAlias(alias string, conf *MainConfig) []ServerListConfig {
	if alias == "" {
		return nil
	}
	var result []ServerListConfig
	for _, server := range conf.ServerLists {
		if server.Alias == alias {
			result = append(result, server)
		}
	}
	return result
}
