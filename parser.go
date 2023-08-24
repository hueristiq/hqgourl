package hqgourl

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Parse parses a raw url into a URL structure.
//
// It uses the  `net/url`'s Parse() internally, but it slightly changes its behavior:
//  1. It forces the default scheme, if the url doesnt have a scheme, to http
//  2. It favors absolute paths over relative ones, thus "example.com"
//     is parsed into url.Host instead of url.Path.
//  3. It lowercases the Host (not only the Scheme).
//
// Parse parses a raw url into a URL structure.
func Parse(rawURL string) (parsedURL *URL, err error) {
	return ParseWithDefaultScheme(rawURL, HTTP)
}

func ParseWithDefaultScheme(rawURL, defaultScheme string) (parsedURL *URL, err error) {
	parsedURL = &URL{}

	// Ensure the rawURL has a default scheme if missing
	rawURL = AddDefaultScheme(rawURL, defaultScheme)

	// Use net/url's Parse to get the base URL structure
	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("[hqgoutils/url]: %w", err)

		return
	}

	// 1. Original
	parsedURL.Original = rawURL

	// 2. Domains, 3. Ports
	parsedURL.Domain, parsedURL.Port = SplitHost(parsedURL.URL.Host)

	// 4. ETLDPlusOne - Extract ETLDPlusOne
	parsedURL.ETLDPlusOne, err = publicsuffix.EffectiveTLDPlusOne(parsedURL.Domain)
	if err != nil {
		err = fmt.Errorf("[hqgoutils/url] %w", err)

		return
	}

	// 5. RootDomain, 6. TLD - Determine the RootDomain and TLD
	parsedURL.RootDomain, parsedURL.TLD = splitETLDPlusOne(parsedURL.ETLDPlusOne)

	// 7. Subdomain - Determine the Subdomain, if any
	if subdomain := strings.TrimSuffix(parsedURL.Domain, "."+parsedURL.ETLDPlusOne); subdomain != parsedURL.Domain {
		parsedURL.Subdomain = subdomain
	}

	// 8. Extension
	parsedURL.Extension = path.Ext(parsedURL.Path)

	return
}

// AddDefaultScheme ensures a scheme is added if none exists.
func AddDefaultScheme(rawURL, scheme string) string {
	switch {
	case strings.HasPrefix(rawURL, "//"):
		return scheme + ":" + rawURL
	case strings.HasPrefix(rawURL, SchemeSeparator):
		return scheme + rawURL
	case !strings.Contains(rawURL, "//"):
		return scheme + SchemeSeparator + rawURL
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

// Used helper function splitETLDPlusOne to clearly separate the logic of splitting ETLD+1.
func splitETLDPlusOne(etldPlusOne string) (rootDomain, tld string) {
	rootDomain, tld, _ = strings.Cut(etldPlusOne, ".")

	return
}
