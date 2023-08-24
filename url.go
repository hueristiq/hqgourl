package hqgourl

import "net/url"

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
	// Scheme      string    // e.g. -> https
	// Opaque      string    // encoded opaque data
	// User        *Userinfo // username and password information
	// Host        string    // e.g. -> sub.example.com, sub.example.com:8080
	// Path        string    // path (relative paths may omit leading slash) e.g. -> /path/to/file.txt
	// RawPath     string    // encoded path hint (see EscapedPath method)
	// OmitHost    bool      // do not emit empty host (authority)
	// ForceQuery  bool      // append a query ('?') even if RawQuery is empty
	// RawQuery    string    // encoded query values, without '?'
	// Fragment    string    // fragment for references, without '#'
	// RawFragment string    // encoded fragment hint (see EscapedFragment method)
	*url.URL

	Original    string
	Domain      string // e.g. -> sub.example.com
	Port        string // e.g. -> 8080
	ETLDPlusOne string // e.g. -> example.com
	RootDomain  string // e.g. -> example
	TLD         string // e.g. -> com
	Subdomain   string // e.g. -> sub
	Extension   string // e.g. -> .txt
}
