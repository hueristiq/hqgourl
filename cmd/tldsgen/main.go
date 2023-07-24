package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/hueristiq/hqgolog"
	"github.com/spf13/pflag"
)

var (
	path string

	tldsTmpl = template.Must(template.New("schemes").Parse(`// This file is autogenerated. Please do not edit manually.

package tlds

// TLDs is a sorted list of all public top-level domains.
// Sources:{{range $_, $source := .Sources}}
//   - {{$source}}{{end}}
var TLDs = []string{
{{range $_, $TLD := .TLDs}}` + "\t`" + `{{$TLD}}` + "`" + `,
{{end}}}
`))
)

func init() {
	pflag.StringVarP(&path, "path", "p", "", "")

	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		h := "USAGE:\n"
		h += "  tldsgen [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += " -p, --path string                 output path\n"

		fmt.Fprintln(os.Stderr, h)
	}

	pflag.Parse()
}

func main() {
	hqgolog.Info().Msgf("Generating %s...", path)

	TLDs, sources := getTLDList()

	if err := writeTlds(TLDs, sources, path); err != nil {
		hqgolog.Fatal().Msg(err.Error())
	}
}

func getTLDList() (TLDs, sources []string) {
	var wg sync.WaitGroup

	tldSet := make(map[string]bool)
	fromURL := func(source, pat string) {
		sources = append(sources, source)

		wg.Add(1)

		go fetchFromURL(&wg, source, pat, tldSet)
	}

	fromURL("https://data.iana.org/TLD/tlds-alpha-by-domain.txt", `^[^#]+$`)
	fromURL("https://publicsuffix.org/list/effective_tld_names.dat", `^[^/.]+$`)

	wg.Wait()

	TLDs = make([]string, 0, len(tldSet))
	for tld := range tldSet {
		TLDs = append(TLDs, tld)
	}

	sort.Strings(TLDs)

	return
}

func fetchFromURL(wg *sync.WaitGroup, source, pat string, tldSet map[string]bool) {
	defer wg.Done()

	resp, err := http.Get(source)
	if err == nil && resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
	}

	if err != nil {
		panic(fmt.Errorf("%s: %s", source, err))
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	re := regexp.MustCompile(pat)

	for scanner.Scan() {
		line := scanner.Text()
		tld := re.FindString(line)
		tld = cleanTld(tld)

		if tld == "" {
			continue
		}

		tldSet[tld] = true
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("%s: %s", source, err))
	}
}

func cleanTld(tld string) string {
	tld = strings.ToLower(tld)
	if strings.HasPrefix(tld, "xn--") {
		return ""
	}

	return tld
}

func writeTlds(TLDs, sources []string, path string) error {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	return tldsTmpl.Execute(f, struct {
		TLDs    []string
		Sources []string
	}{
		TLDs:    TLDs,
		Sources: sources,
	})
}