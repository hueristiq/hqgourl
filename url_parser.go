package hqgourl

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// URL extends the standard net/url URL struct with additional domain-related fields.
type URL struct {
	// Embedding the standard URL struct for base functionalities.
	*url.URL

	Subdomain      string
	RootDomain     string
	TopLevelDomain string
	Port           int
	Extension      string
}

// URLParser encapsulates the logic for parsing URLs with additional domain-specific information.
// It uses the standard URL parsing capabilities and enhances them with subdomain, root domain,
// and TLD extraction. It also handles the addition of a default scheme if one is not present.
type URLParser struct {
	// DefaultScheme is the default URL scheme to use if not specified in the URL.
	scheme string

	// DomainParser is used to parse domain-specific details from the URL.
	dp *DomainParser
}

func (up *URLParser) DefaultScheme() (scheme string) {
	return up.scheme
}

// Parse takes a raw URL string and parses it into a URL struct.
// It enhances the standard parsing by attaching domain-specific details like subdomain, root domain, and TLD.
// The method also ensures that a default scheme is set if the URL does not specify one.
func (up *URLParser) Parse(rawURL string) (parsedURL *URL, err error) {
	parsedURL = &URL{}

	// Add default scheme if necessary
	if up.scheme != "" {
		rawURL = addScheme(rawURL, up.scheme)
	}

	// Standard URL parsing
	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("error parsing URL: %w", err)

		return
	}

	// Split host and port, and handle errors
	parsedURL.Host, parsedURL.Port, err = splitHostPort(parsedURL.Host)
	if err != nil {
		err = fmt.Errorf("error splitting host and port: %w", err)

		return
	}

	domainRegex := regexp.MustCompile(`(?i)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}`)

	if domainRegex.MatchString(parsedURL.Host) {
		// Parse domain-specific parts
		parsedURL.Subdomain, parsedURL.RootDomain, parsedURL.TopLevelDomain = up.dp.Parse(parsedURL.Host)
	}

	// Extract file extension from the path
	parsedURL.Extension = path.Ext(parsedURL.Path)

	return
}

// URLParserOptionsFunc defines a function type for configuring a URLParser.
type URLParserOptionsFunc func(*URLParser)

// URLParserInterface defines the interface for URL parsing functionality.
type URLParserInterface interface {
	DefaultScheme() (scheme string)
	Parse(rawURL string) (parsedURL *URL, err error)
}

var _ URLParserInterface = &URLParser{}

// NewURLParser creates a new URLParser with the given options.
func NewURLParser(opts ...URLParserOptionsFunc) (up *URLParser) {
	up = &URLParser{}

	// Initialize the DomainParser
	dp := NewDomainParser()
	up.dp = dp

	// Apply additional options
	for _, opt := range opts {
		opt(up)
	}

	return
}

// URLParserWithDefaultScheme returns a URLParserOptionsFunc to set a default scheme.
func URLParserWithDefaultScheme(scheme string) URLParserOptionsFunc {
	return func(up *URLParser) {
		up.scheme = scheme
	}
}

// addScheme is a helper function that adds a scheme to the URL if it's missing.
// This ensures that the URL is parsed correctly as a network address rather than a relative path.
// This makes net/url.Parse() not put both host and path into the (relative) path.
func addScheme(inURL, scheme string) (outURL string) {
	switch {
	case strings.HasPrefix(inURL, "//"):
		outURL = scheme + ":" + inURL
	case strings.HasPrefix(inURL, "://"):
		outURL = scheme + inURL
	case !strings.Contains(inURL, "//"):
		outURL = scheme + "://" + inURL
	default:
		outURL = inURL
	}

	return
}

// splitHostPort separates the host and port in a network address.
// It is designed to handle both IPv4 and IPv6 addresses and gracefully manages URLs without a port.
// Unlike net.SplitHostPort(), it doesn't remove brackets from [IPv6] hosts.
func splitHostPort(address string) (host string, port int, err error) {
	host = address

	// Check for the last colon, which should separate host and port
	i := strings.LastIndex(address, ":")
	if i == -1 {
		return
	}

	// Handle IPv6 addresses enclosed in brackets
	if strings.HasPrefix(address, "[") && strings.Contains(address[i:], "]") {
		return
	}

	// Split the host and port
	host = address[:i]

	if port, err = strconv.Atoi(address[i+1:]); err != nil {
		return
	}

	return
}
