package config

import "strings"

// Parses hostnames from a string like "host1:alias1,alias2|host2:alias3,alias4"
func parseHostnames(hostnames string) map[string][]string {
	domainAliases := make(map[string][]string)
	if strings.TrimSpace(hostnames) == "" {
		return domainAliases
	}

	groups := strings.Split(hostnames, "|")
	for _, group := range groups {
		if group == "" {
			continue
		}
		parts := strings.SplitN(group, ":", 2)
		hostname := strings.TrimSpace(parts[0])
		var aliases []string
		if len(parts) == 2 {
			for _, a := range strings.Split(parts[1], ",") {
				a = strings.TrimSpace(a)
				if a != "" {
					aliases = append(aliases, a)
				}
			}
		}
		if hostname != "" {
			domainAliases[hostname] = aliases
		}
	}
	return domainAliases
}

// Parses backends from a string like "host1=http://service1:80|host2=https://service2:443"
func parseBackends(backends string) map[string]string {
	result := make(map[string]string)
	if strings.TrimSpace(backends) == "" {
		return result
	}
	pairs := strings.Split(backends, "|")
	for _, p := range pairs {
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		host := strings.TrimSpace(kv[0])
		url := strings.TrimSpace(kv[1])
		if host != "" && url != "" {
			result[host] = url
		}
	}
	return result
}
