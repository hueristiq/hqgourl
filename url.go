package hqgourl

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// URL represents a parsed URL (technically, a URI reference).
//
// The general form represented is:
//
//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]
//
// URLs that do not start with a slash after the scheme are interpreted as:
//
//	scheme:opaque[?query][#fragment]
//
// https://sub.example.com:8080/path/to/file.txt
type URL struct {
	*url.URL
	// Scheme      string    // e.g https
	// Opaque      string    // encoded opaque data
	// User        *Userinfo // username and password information
	// Host        string    // e.g. sub.example.com, sub.example.com:8080
	// Path        string    // path (relative paths may omit leading slash) e.g /path/to/file.txt
	// RawPath     string    // encoded path hint (see EscapedPath method)
	// OmitHost    bool      // do not emit empty host (authority)
	// ForceQuery  bool      // append a query ('?') even if RawQuery is empty
	// RawQuery    string    // encoded query values, without '?'
	// Fragment    string    // fragment for references, without '#'
	// RawFragment string    // encoded fragment hint (see EscapedFragment method)
	Domain      string // e.g. sub.example.com
	ETLDPlusOne string // e.g. example.com
	Subdomain   string // e.g. sub
	RootDomain  string // e.g. example
	TLD         string // e.g. com
	Port        string // e.g. 8080
	Extension   string // e.g. txt
}

// Parse parses a raw url into a URL structure.
//
// It uses the  `net/url`'s Parse() internally, but it slightly changes its behavior:
//  1. It forces the default scheme, if the url doesnt have a scheme, to http
//  2. It favors absolute paths over relative ones, thus "example.com"
//     is parsed into url.Host instead of url.Path.
//  3. It lowercases the Host (not only the Scheme).
func Parse(rawURL string) (parsedURL *URL, err error) {
	const defaultScheme string = "http"

	rawURL = AddDefaultScheme(rawURL, defaultScheme)

	parsedURL = &URL{}

	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("[hqgoutils/url]: %w", err)

		return
	}

	// Host = Domain + Port
	parsedURL.Domain, parsedURL.Port = SplitHost(parsedURL.URL.Host)

	// ETLDPlusOne
	parsedURL.ETLDPlusOne, err = publicsuffix.EffectiveTLDPlusOne(parsedURL.Domain)
	if err != nil {
		err = fmt.Errorf("[hqgoutils/url] %w", err)

		return
	}

	// RootDomain + TLD
	i := strings.Index(parsedURL.ETLDPlusOne, ".")
	parsedURL.RootDomain = parsedURL.ETLDPlusOne[0:i]
	parsedURL.TLD = parsedURL.ETLDPlusOne[i+1:]

	// Subdomain
	if rest := strings.TrimSuffix(parsedURL.Domain, "."+parsedURL.ETLDPlusOne); rest != parsedURL.Domain {
		parsedURL.Subdomain = rest
	}

	// Extension
	parsedURL.Extension = path.Ext(parsedURL.Path)

	return
}

// AddDefaultScheme ensures a scheme is added if none exists.
func AddDefaultScheme(rawURL, scheme string) string {
	switch {
	case strings.HasPrefix(rawURL, "//"):
		return scheme + ":" + rawURL
	case strings.Contains(rawURL, "://") && !strings.HasPrefix(rawURL, "http"):
		return scheme + rawURL
	case !strings.Contains(rawURL, "://"):
		return scheme + "://" + rawURL
	default:
		return rawURL
	}
}

// splitHost splits the host into domain and port.
func SplitHost(host string) (domain string, port string) {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i], host[i+1:]
		} else if host[i] < '0' || host[i] > '9' {
			domain = host
		}
	}
	return
}
