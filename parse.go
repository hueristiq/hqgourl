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
// Parse parses a raw url into a URL structure.
func Parse(rawURL string) (parsedURL *URL, err error) {
	parsedURL = &URL{}
	
	const defaultScheme = "http"

	// Ensure the rawURL has a default scheme if missing
	rawURL = AddDefaultScheme(rawURL, defaultScheme)

	// Use net/url's Parse to get the base URL structure
	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("[hqgoutils/url]: %w", err)

		return
	}

	parsedURL.Domain, parsedURL.Port = SplitHost(parsedURL.URL.Host)

	// ETLDPlusOne - Extract ETLDPlusOne
	parsedURL.ETLDPlusOne, err = publicsuffix.EffectiveTLDPlusOne(parsedURL.Domain)
	if err != nil {
		err = fmt.Errorf("[hqgoutils/url] %w", err)

		return
	}

	// RootDomain + TLD - Determine the RootDomain and TLD
	parsedURL.RootDomain, parsedURL.TLD = splitETLDPlusOne(parsedURL.ETLDPlusOne)

	// Subdomain - Determine the Subdomain, if any
	if subdomain := strings.TrimSuffix(parsedURL.Domain, "."+parsedURL.ETLDPlusOne); subdomain != parsedURL.Domain {
		parsedURL.Subdomain = subdomain
	}

	// Extension
	parsedURL.Extension = path.Ext(parsedURL.Path)

	return
}

// Used helper function splitETLDPlusOne to clearly separate the logic of splitting ETLD+1.
func splitETLDPlusOne(etldPlusOne string) (rootDomain, tld string) {
	i := strings.Index(etldPlusOne, ".")
	rootDomain = etldPlusOne[:i]
	tld = etldPlusOne[i+1:]
	
	return
}

// AddDefaultScheme ensures a scheme is added if none exists.
func AddDefaultScheme(rawURL, scheme string) string {
	switch {
	case strings.HasPrefix(rawURL, "//"):
		return scheme + ":" + rawURL
	case strings.HasPrefix(rawURL, "://"):
		return scheme + rawURL
	case !strings.Contains(rawURL, "//"):
		return scheme + "://" + rawURL
	default:
		return rawURL
	}
}

// SplitHost splits the host into domain and port.
func SplitHost(host string) (domain, port string) {
	if i := strings.LastIndex(host, ":"); i != -1 {
		domain = host[:i]
		port = host[i+1:]
		
		return
	}

	domain = host

	return
}
