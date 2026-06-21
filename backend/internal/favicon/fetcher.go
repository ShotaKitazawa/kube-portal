package favicon

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Fetcher struct {
	client *http.Client
	cache  sync.Map // string (origin) -> string (favicon URL or "")
}

func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

// Resolve returns the favicon URL for the given origin (e.g. "https://example.com").
// Returns "" if no favicon is found. Results are cached by origin.
func (f *Fetcher) Resolve(ctx context.Context, origin string) string {
	if v, ok := f.cache.Load(origin); ok {
		return v.(string)
	}
	result := f.resolve(ctx, origin)
	f.cache.Store(origin, result)
	return result
}

func (f *Fetcher) resolve(ctx context.Context, origin string) string {
	if u := f.fromHTML(ctx, origin); u != "" {
		return u
	}
	faviconURL := origin + "/favicon.ico"
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, faviconURL, nil)
	if err != nil {
		return ""
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return ""
	}
	resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return faviconURL
	}
	return ""
}

func (f *Fetcher) fromHTML(ctx context.Context, origin string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, origin, nil)
	if err != nil {
		return ""
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return ""
	}
	href := findFaviconHref(doc)
	if href == "" {
		return ""
	}
	return toAbsolute(origin, href)
}

func findFaviconHref(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "link" {
		if hasFaviconRel(attrVal(n, "rel")) {
			if href := attrVal(n, "href"); href != "" {
				return href
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if h := findFaviconHref(c); h != "" {
			return h
		}
	}
	return ""
}

// hasFaviconRel checks whether the rel attribute (space-separated tokens) contains "icon".
func hasFaviconRel(rel string) bool {
	for _, r := range strings.Fields(rel) {
		if strings.EqualFold(r, "icon") {
			return true
		}
	}
	return false
}

func attrVal(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func toAbsolute(origin, href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}
	base, err := url.Parse(origin)
	if err != nil {
		return ""
	}
	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return base.ResolveReference(ref).String()
}
