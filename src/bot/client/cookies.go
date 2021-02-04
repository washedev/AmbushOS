package client

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

type ExportableCookieJar struct {
	jar        *cookiejar.Jar
	allCookies map[url.URL][]*http.Cookie
	sync.RWMutex
}

func NewExportableCookieJar() *ExportableCookieJar {
	realJar, _ := cookiejar.New(nil)

	e := &ExportableCookieJar{
		jar:        realJar,
		allCookies: make(map[url.URL][]*http.Cookie),
	}

	return e
}

func (jar *ExportableCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.Lock()
	defer jar.Unlock()
	jar.allCookies[*u] = cookies
	jar.jar.SetCookies(u, cookies)
}

func (jar *ExportableCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return jar.jar.Cookies(u)
}

func (jar *ExportableCookieJar) ExportAllCookies() map[url.URL][]*http.Cookie {
	jar.RLock()
	defer jar.RUnlock()

	copied := make(map[url.URL][]*http.Cookie)

	for u, c := range jar.allCookies {
		copied[u] = c
	}

	return copied
}
