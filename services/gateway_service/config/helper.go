package config

import (
	"net/url"
	"strings"
)

// Parses hostnames from a string like "host1:alias1,alias2|host2:alias3,alias4"
func parseHostnames(hostnames string) map[string][]string {
	domainAliases := make(map[string][]string)
	if strings.TrimSpace(hostnames) == "" {
		return domainAliases
	}

	groups := strings.Split(hostnames, "|")
	for _, group := range groups {
		group = strings.TrimSpace(group)
		if group == "" {
			continue
		}
		parts := strings.SplitN(group, ":", 2)
		hostname := strings.ToLower(strings.TrimSpace(parts[0]))
		var aliases []string
		if len(parts) == 2 {
			for _, a := range strings.Split(parts[1], ",") {
				a = strings.ToLower(strings.TrimSpace(a))
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
// - Trims whitespace and lowercases host keys
// - Ignores invalid or unparsable URL values
// - Tolerates accidental newlines/commas around entries
func parseBackends(backends string) map[string]string {
	result := make(map[string]string)
	if strings.TrimSpace(backends) == "" {
		return result
	}
	// Normalize newlines to separators we already handle
	normalized := strings.ReplaceAll(backends, "\n", "|")
	normalized = strings.ReplaceAll(normalized, "\r", "|")
	pairs := strings.Split(normalized, "|")
	for _, p := range pairs {
		p = strings.TrimSpace(strings.Trim(p, ","))
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		host := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])
		if host == "" || val == "" {
			continue
		}
		// Validate URL and ensure scheme is present
		u, err := url.Parse(val)
		if err != nil {
			continue
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			continue
		}
		if u.Host == "" {
			continue
		}
		result[host] = u.String()
	}
	return result
}
