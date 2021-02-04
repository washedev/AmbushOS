package utils

import (
	"fmt"
	"strings"

	"bot/client"
	"net/http"
	"net/url"
)

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
	URL    string `json:"url"`
	Secure bool   `json:"secure"`
}

func GetCookie(jar *client.ExportableCookieJar, name string) string {

	for _, cookies := range jar.ExportAllCookies() {
		for _, cookie := range cookies {
			if name == cookie.Name {
				return cookie.Value
			}
		}
	}

	return ""
}

func SetCookie(jar *client.ExportableCookieJar, name, value, domain, URL string) (*client.ExportableCookieJar, error) {

	new := client.NewExportableCookieJar()

	for host, cookies := range jar.ExportAllCookies() {

		cks := make([]*http.Cookie, 0)

		for _, cookie := range cookies {
			if name != cookie.Name {
				cks = append(cks, cookie)
			}
		}

		new.SetCookies(&host, cks)
	}

	cookie := http.Cookie{
		Name:   name,
		Value:  value,
		Domain: domain,
		Path:   "/",
	}

	u, err := url.Parse(URL)

	if err != nil {
		return nil, err
	}

	new.SetCookies(u, []*http.Cookie{&cookie})

	return new, nil
}

func PopCookie(jar *client.ExportableCookieJar, name string) *client.ExportableCookieJar {

	new := client.NewExportableCookieJar()

	for host, cookies := range jar.ExportAllCookies() {

		cks := make([]*http.Cookie, 0)

		for _, cookie := range cookies {
			if name != cookie.Name {
				cks = append(cks, cookie)
			}
		}

		new.SetCookies(&host, cks)
	}

	return new
}

func ExtractCookies(jar *client.ExportableCookieJar) []*Cookie {

	r := make([]*Cookie, 0)

	for u, cookies := range jar.ExportAllCookies() {
		for _, cookie := range cookies {

			c := &Cookie{
				Name:   cookie.Name,
				Value:  cookie.Value,
				Domain: cookie.Domain,
				Path:   cookie.Path,
				Secure: true,
			}

			if strings.HasPrefix(u.Host, ".") {
				c.URL = fmt.Sprintf("https://%v", strings.Trim(u.Host, "."))
			} else {
				c.URL = fmt.Sprintf("https://%v", u.Host)
			}

			r = append(r, c)
		}
	}

	return r
}
