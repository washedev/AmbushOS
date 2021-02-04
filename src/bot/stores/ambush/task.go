package ambush

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"bot/client"
	"bot/log"
)

type Task struct {
	ID string

	SKU   string
	Sizes []string

	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	Address1    string
	Address2    string
	State       string
	City        string
	Postcode    string
	Country     string
	CountryID   string
	Currency    string

	Region string

	Delay   int
	Timeout int
	Webhook string

	Proxies []string

	Client    *http.Client
	CookieJar client.ExportableCookieJar

	BagID                    string
	MerchantID               int
	Scale                    string
	OrderID                  int
	ShippingPrice            int
	ShippingFormattedPrice   string
	ShippingCostType         int
	ShippingDescription      string
	ShippingID               int
	ShippingName             string
	ShippingType             string
	MinEstimatedDeliveryHour int
	MaxEstimatedDeliveryHour int
	GrandTotal               int
	CTX                      string
	PaymentIntentID          string
	RedirectURL              string
	PayPalURL                string
	CheckoutURL              string

	ProductName    string
	ProductImage   string
	ProductPrice   string
	ProductSize    string
	ProductVariant string
}

func (t *Task) Debug(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Debugln(s, t.ID)
}

func (t *Task) Info(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Infoln(s, t.ID)
}

func (t *Task) Warn(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Warningln(s, t.ID)
}

func (t *Task) Error(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	log.Errorln(s, t.ID)
}

func (t *Task) HandleError(err error) bool {
	if err != nil {
		t.Error(err.Error())
		return true
	}
	return false
}

func (t *Task) Sleep() {
	time.Sleep(time.Millisecond * time.Duration(t.Delay))
}

func (t *Task) NewClient() (*http.Client, error) {
	return &http.Client{}, nil
}

func (t *Task) SetupClient() {
	jar := client.NewExportableCookieJar()
	t.Client.Jar = jar
	t.CookieJar = *jar

	t.Client.Timeout = time.Millisecond * time.Duration(t.Timeout)

	t.Rotate()
}

func newRoundTripper(u string) (*http.Transport, error) {
	proxyURL, err := url.Parse(u)

	if err != nil {
		return nil, err
	}

	return &http.Transport{Proxy: http.ProxyURL(proxyURL)}, nil
}

func (t *Task) Rotate() {
	var roundTripper http.RoundTripper
	var err error

	if len(t.Proxies) > 0 {
		proxyURL := t.Proxies[rand.Intn(len(t.Proxies))]
		roundTripper, err = newRoundTripper(proxyURL)
		if err != nil {
			t.Error(err.Error())
			return
		}
	} else {
		roundTripper = &http.Transport{}
	}

	t.Client.Transport = roundTripper
}

func (task *Task) SetAllowRedirects(allow bool) {
	if allow {
		task.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return nil
		}
	} else {
		task.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

func (t *Task) SleepAndRotate() {
	t.Sleep()
	t.Rotate()
}
