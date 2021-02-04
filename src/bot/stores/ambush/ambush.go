package ambush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"bot/utils"
)

type apiUsersMeResponse struct {
	BagID string `json:"bagId"`
}

type apiCommerceV1Products struct {
	BreadCrumbs []struct {
		Text string `json:"text"`
		Link string `json:"link"`
	} `json:"breadCrumbs"`
	ImageGroups []struct {
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"imageGroups"`
	Sizes []struct {
		SizeID          string `json:"sizeId"`
		SizeDescription string `json:"sizeDescription"`
		Scale           string `json:"scale"`
		Variants        []struct {
			MerchantID     int    `json:"merchantId"`
			FormattedPrice string `json:"formattedPrice"`
		} `json:"variants"`
	} `json:"sizes"`
}

type apiCommerceV1BagsPayload struct {
	MerchantID       int    `json:"merchantId"`
	ProductID        string `json:"productId"`
	Quantity         int    `json:"quantity"`
	Scale            string `json:"scale"`
	Size             string `json:"size"`
	CustomAttributes string `json:"customAttributes"`
}

type apiCommerceV1BagsResponse struct {
	BagSummary struct {
		GrandTotal float64 `json:"grandTotal"`
	} `json:"BagSummary"`
}

type apiCheckoutV1OrdersPayload struct {
	BagID            string `json:"bagId"`
	GuestUserEmail   string `json:"guestUserEmail"`
	UsePaymentIntent bool   `json:"usePaymentIntent"`
}

type apiCheckoutV1OrdersResponse struct {
	ID int `json:"id"`
}

type apiCheckoutV1OrdersResponse2 struct {
	CheckoutOrder struct {
		GrandTotal      float64 `json:"grandTotal"`
		PaymentIntentID string  `json:"paymentIntentId"`
	}
	ShippingOptions []struct {
		Price            float64 `json:"price"`
		FormattedPrice   string  `json:"formattedPrice"`
		ShippingCostType int     `json:"shippingCostType"`
		ShippingService  struct {
			Description              string  `json:"description"`
			ID                       int     `json:"id"`
			Name                     string  `json:"name"`
			Type                     string  `json:"type"`
			MinEstimatedDeliveryHour float64 `json:"minEstimatedDeliveryHour"`
			MaxEstimatedDeliveryHour float64 `json:"maxEstimatedDeliveryHour"`
		} `json:"shippingService"`
	} `json:"shippingOptions"`
}

type apiCheckoutV1OrderChargesResponse struct {
	RedirectURL string `json:"redirectUrl"`
}

func (t *Task) Start() {
	var err error

	if err = utils.CheckProfile(t.Email, t.FirstName, t.LastName, t.PhoneNumber, t.Address1, t.City, t.Postcode, t.Country); err != nil {
		t.HandleError(err)
		return
	}

	t.Email, t.FirstName, t.LastName, t.PhoneNumber, t.Address1 = utils.WrapProfile(t.Email, t.FirstName, t.LastName, t.PhoneNumber, t.Address1)

	t.Client, err = t.NewClient()

	if t.HandleError(err) {
		return
	}

	t.SetupClient()

	t.CreateSession()
}

func (t *Task) CreateSession() {

	t.Warn("Creating session")

	for {

		u := "https://www.ambushdesign.com/en-it/api/users/me"

		req, err := http.NewRequest("GET", u, nil)

		if err != nil {
			t.Error("Error creating session - 0 - %v", err.Error())
			continue
		}

		req.Header.Set("Cache-Control", "max-age=0")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Accept-Language", "en-US,en-GB;q=0.9,en;q=0.8,it;q=0.7")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error creating session - 1 - %v", err.Error())
			t.Rotate()
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			t.Error("Error creating session - 2 - %v", err.Error())
			t.Sleep()
			continue
		}

		switch resp.StatusCode {
		case 200:
			body := new(apiUsersMeResponse)

			if err := json.Unmarshal(b, &body); err != nil {
				t.Error("Error creating session - 3 - %v", err.Error())
				t.SleepAndRotate()
				continue
			}

			t.BagID = body.BagID

			t.CheckStock()
			return

		case 400:
			t.Error("Error creating session - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error creating session - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error creating session - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error creating session - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error creating session - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) CheckStock() {

	t.Warn("Checking stock")

	for {

		u := fmt.Sprintf("https://www.ambushdesign.com/api/commerce/v1/products/%v", t.SKU)

		req, err := http.NewRequest("GET", u, nil)

		if err != nil {
			t.Error("Error checking stock - 0 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Encoding", "deflate, br")
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("X-Newrelic-ID", "VQUCV1ZUGwIFVlBRDgcA")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error checking stock - 1 - %v", err.Error())
			t.Rotate()
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			t.Error("Error checking stock - 2 - %v", err.Error())
			t.Sleep()
			continue
		}

		switch resp.StatusCode {
		case 200:
			body := new(apiCommerceV1Products)

			if err := json.Unmarshal(b, &body); err != nil {
				t.Error("Error checking stock - 3 - %v", err.Error())
				t.SleepAndRotate()
				continue
			}

			for _, product := range body.BreadCrumbs {
				if product.Link == "" {
					t.ProductName = product.Text
					break
				}
			}

			if len(body.ImageGroups) > 0 {
				if len(body.ImageGroups[0].Images) > 0 {
					t.ProductImage = body.ImageGroups[0].Images[0].URL
				}
			}

			sizes := make([]map[string]string, 0)

			for _, size := range body.Sizes {
				if len(size.Variants) > 0 {
					t.ProductPrice = size.Variants[0].FormattedPrice
					t.MerchantID = size.Variants[0].MerchantID
				}

				sizes = append(sizes, map[string]string{
					"size":   size.SizeDescription,
					"sizeId": size.SizeID,
					"scale":  size.Scale,
				})
			}

			if len(sizes) == 0 {
				t.Warn("Waiting for restock")
				t.Sleep()
				continue
			}

			size := sizes[rand.Intn(len(sizes))]

			t.ProductSize = size["size"]
			t.ProductVariant = size["sizeId"]

			t.Scale = size["scale"]

			t.Info("Stock found - %v - %v - %v", t.ProductName, t.ProductPrice, t.ProductSize)

			t.AddToCart()
			return

		case 400:
			t.Error("Error checking stock - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error checking stock - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error checking stock - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error checking stock - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error checking stock - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) AddToCart() {

	t.Warn("Carting %v", t.ProductName)

	for {

		u := fmt.Sprintf("https://www.ambushdesign.com/api/commerce/v1/bags/%v/items", t.BagID)

		p := apiCommerceV1BagsPayload{
			MerchantID: t.MerchantID,
			ProductID:  t.SKU,
			Quantity:   1,
			Scale:      t.Scale,
			Size:       t.ProductVariant,
		}

		payload, err := json.Marshal(&p)

		if err != nil {
			t.Error("Error adding to cart - 0 - %v", err.Error())
			t.Sleep()
			continue
		}

		req, err := http.NewRequest("POST", u, bytes.NewReader(payload))

		if err != nil {
			t.Error("Error adding to cart - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://www.ambushdesign.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error adding to cart - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			t.Error("Error adding to cart - 3 - %v", err.Error())
			t.Sleep()
			continue
		}

		switch resp.StatusCode {
		case 200:
			body := new(apiCommerceV1BagsResponse)

			if err := json.Unmarshal(b, &body); err != nil {
				t.Error("Error adding to cart - 4 - %v", err.Error())
				t.SleepAndRotate()
				continue
			}

			if body.BagSummary.GrandTotal > 0 {
				t.Info("%v added to cart", t.ProductName)
				t.SubmitGuest()
				return
			} else {
				t.Error("Product OOS")
				t.CheckStock()
				return
			}

		case 400:
			t.Error("Error adding to cart - Product OOS [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error adding to cart - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error adding to cart - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error adding to cart - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error adding to cart - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) SubmitGuest() {

	t.Warn("Submitting guest")

	for {

		u := "https://www.ambushdesign.com/api/checkout/v1/orders"

		p := apiCheckoutV1OrdersPayload{
			BagID:            t.BagID,
			GuestUserEmail:   t.Email,
			UsePaymentIntent: true,
		}

		payload, err := json.Marshal(&p)

		if err != nil {
			t.Error("Error submitting guest - 0 - %v", err.Error())
			t.Sleep()
			continue
		}

		req, err := http.NewRequest("POST", u, bytes.NewReader(payload))

		if err != nil {
			t.Error("Error submitting guest - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://www.ambushdesign.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error submitting guest - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			t.Error("Error submitting guest - 3 - %v", err.Error())
			t.Sleep()
			continue
		}

		switch resp.StatusCode {
		case 201:
			body := new(apiCheckoutV1OrdersResponse)

			if err := json.Unmarshal(b, &body); err != nil {
				t.Error("Error submitting guest - 4 - %v", err.Error())
				t.SleepAndRotate()
				continue
			}

			t.OrderID = body.ID
			t.SubmitShipping()
			return

		case 400:
			t.Error("Error submitting guest - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error submitting guest - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error submitting guest - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error submitting guest - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error submitting guest - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) SubmitShipping() {

	t.Warn("Setting shipping")

	for {

		u := fmt.Sprintf("https://www.ambushdesign.com/api/checkout/v1/orders/%v", t.OrderID)

		p := map[string]interface{}{
			"shippingAddress": map[string]interface{}{
				"firstName": t.FirstName,
				"lastName":  t.LastName,
				"country": map[string]string{
					"name": utils.GetFullCountry(t.Country),
					"id":   t.CountryID,
				},
				"addressLine1": t.Address1,
				"addressLine2": t.Address2,
				"addressLine3": "",
				"city": map[string]string{
					"name": t.City,
				},
				"state": map[string]string{
					"name": t.State,
				},
				"zipCode":   t.Postcode,
				"phone":     t.PhoneNumber,
				"vatNumber": "",
			},
			"billingAddress": map[string]interface{}{
				"firstName": t.FirstName,
				"lastName":  t.LastName,
				"country": map[string]string{
					"name": utils.GetFullCountry(t.Country),
					"id":   t.CountryID,
				},
				"addressLine1": t.Address1,
				"addressLine2": t.Address2,
				"addressLine3": "",
				"city": map[string]string{
					"name": t.City,
				},
				"state": map[string]string{
					"name": t.State,
				},
				"zipCode":   t.Postcode,
				"phone":     t.PhoneNumber,
				"vatNumber": "",
			},
		}

		payload, err := json.Marshal(&p)

		if err != nil {
			t.Error("Error setting shipping - 0 - %v", err.Error())
			t.Sleep()
			continue
		}

		req, err := http.NewRequest("PATCH", u, bytes.NewReader(payload))

		if err != nil {
			t.Error("Error setting shipping - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://www.ambushdesign.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error setting shipping - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			t.Error("Error setting shipping - 3 - %v", err.Error())
			t.Sleep()
			continue
		}

		switch resp.StatusCode {
		case 200:
			body := new(apiCheckoutV1OrdersResponse2)

			if err := json.Unmarshal(b, &body); err != nil {
				t.Error("Error setting shipping - 4 - %v", err.Error())
				t.SleepAndRotate()
				continue
			}

			t.ShippingPrice = int(body.ShippingOptions[0].Price)
			t.ShippingFormattedPrice = body.ShippingOptions[0].FormattedPrice
			t.ShippingCostType = body.ShippingOptions[0].ShippingCostType
			t.ShippingDescription = body.ShippingOptions[0].ShippingService.Description
			t.ShippingID = body.ShippingOptions[0].ShippingService.ID
			t.ShippingName = body.ShippingOptions[0].ShippingService.Name
			t.ShippingType = body.ShippingOptions[0].ShippingService.Type
			t.MinEstimatedDeliveryHour = int(body.ShippingOptions[0].ShippingService.MinEstimatedDeliveryHour)
			t.MaxEstimatedDeliveryHour = int(body.ShippingOptions[0].ShippingService.MaxEstimatedDeliveryHour)
			t.GrandTotal = int(body.CheckoutOrder.GrandTotal)
			t.PaymentIntentID = body.CheckoutOrder.PaymentIntentID

			t.Info("Shipping set")
			t.SubmitDelivery()
			return

		case 400:
			t.Error("Error setting shipping - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error setting shipping - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error setting shipping - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error setting shipping - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error setting shipping - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) SubmitDelivery() {

	t.Warn("Setting delivery")

	for {

		u := fmt.Sprintf("https://www.ambushdesign.com/api/checkout/v1/orders/%v", t.OrderID)

		p := map[string]interface{}{
			"shippingOption": map[string]interface{}{
				"discount":         0,
				"merchants":        []int{t.MerchantID},
				"price":            t.ShippingPrice,
				"formattedPrice":   t.ShippingFormattedPrice,
				"shippingCostType": t.ShippingCostType,
				"shippingService": map[string]interface{}{
					"description":              t.ShippingDescription,
					"id":                       t.ShippingID,
					"name":                     t.ShippingName,
					"type":                     t.ShippingType,
					"minEstimatedDeliveryHour": t.MinEstimatedDeliveryHour,
					"maxEstimatedDeliveryHour": t.MaxEstimatedDeliveryHour,
					"trackingCodes":            []string{},
				},
				"shippingWithoutCapped": 0,
				"baseFlatRate":          0,
			},
		}

		payload, err := json.Marshal(&p)

		if err != nil {
			t.Error("Error setting delivery - 0 - %v", err.Error())
			t.Sleep()
			continue
		}

		req, err := http.NewRequest("PATCH", u, bytes.NewReader(payload))

		if err != nil {
			t.Error("Error setting delivery - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://www.ambushdesign.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error setting delivery - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		resp.Body.Close()

		switch resp.StatusCode {
		case 200:
			t.Info("Delivery set")
			t.SubmitPayment()
			return
		case 400:
			t.Error("Error setting delivery - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error setting delivery - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error setting delivery - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error setting delivery - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error setting delivery - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) SubmitPayment() {

	t.Warn("Setting payment")

	ctx := utils.GetCookie(&t.CookieJar, "ctx")
	t.CTX, _ = utils.Extract(ctx, "%3a", "%2c")
	i, _ := strconv.Atoi(t.CTX)

	for {

		u := fmt.Sprintf("https://www.ambushdesign.com/api/payment/v1/intents/%v/instruments", t.PaymentIntentID)

		p := map[string]interface{}{
			"method":      "PayPal",
			"option":      "PayPalExp",
			"createToken": false,
			"payer": map[string]interface{}{
				"id":        i,
				"firstName": t.FirstName,
				"lastName":  t.LastName,
				"email":     t.Email,
				"birthDate": nil,
				"address": map[string]interface{}{
					"city": map[string]interface{}{
						"countryId": t.CountryID,
						"id":        0,
						"name":      t.City,
					},
					"country": map[string]interface{}{
						"alpha2Code":  t.Country,
						"alpha3Code":  t.Country,
						"culture":     "it-IT",
						"id":          t.CountryID,
						"name":        utils.GetFullCountry(t.Country),
						"nativeName":  utils.GetFullCountry(t.Country),
						"region":      "Europe",
						"regionId":    0,
						"continentId": 3,
					},
					"id":       "00000000-0000-0000-0000-000000000000",
					"lastName": t.LastName,
					"state": map[string]interface{}{
						"countryId": 0,
						"id":        0,
						"code":      t.State,
						"name":      t.State,
					},
					"userId":                   0,
					"isDefaultBillingAddress":  false,
					"isDefaultShippingAddress": false,
					"addressLine1":             t.Address1,
					"addressLine2":             t.Address2,
					"firstName":                t.FirstName,
					"phone":                    t.PhoneNumber,
					"vatNumber":                "",
					"zipCode":                  t.Postcode,
				},
			},
			"amounts": []map[string]interface{}{{
				"value": t.GrandTotal,
			}},
			"data": map[string]interface{}{},
		}

		payload, err := json.Marshal(&p)

		if err != nil {
			t.Error("Error setting payment - 0 - %v", err.Error())
			t.Sleep()
			continue
		}

		req, err := http.NewRequest("POST", u, bytes.NewReader(payload))

		if err != nil {
			t.Error("Error setting payment - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://www.ambushdesign.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error setting payment - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		resp.Body.Close()

		switch resp.StatusCode {
		case 201:
			t.Info("Payment set")
			t.CheckCharge()
			return
		case 400:
			t.Error("Error setting payment - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error setting payment - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error setting payment - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error setting payment - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error setting payment - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) CheckCharge() {

	t.Warn("Checking order")

	for {

		u := fmt.Sprintf("https://www.ambushdesign.com/api/checkout/v1/orders/%v/charges", t.OrderID)

		p := map[string]string{
			"cancelUrl": "https://www.ambushdesign.com/en-it/commerce/checkout",
			"returnUrl": fmt.Sprintf("https://www.ambushdesign.com/en-it/checkoutdetails?orderid=%v", t.OrderID),
		}

		payload, err := json.Marshal(&p)

		if err != nil {
			t.Error("Error checking order - 0 - %v", err.Error())
			t.Sleep()
			continue
		}

		req, err := http.NewRequest("POST", u, bytes.NewReader(payload))

		if err != nil {
			t.Error("Error checking order - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://www.ambushdesign.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error checking order - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			t.Error("Error checking order - 3 - %v", err.Error())
			t.Sleep()
			continue
		}

		switch resp.StatusCode {
		case 201:

			text := string(b)

			if strings.Contains(text, "Error") {
				t.Error("Payment declined")
				t.FailedWebhook()
				return
			} else if strings.Contains(text, "Processing") {
				body := new(apiCheckoutV1OrderChargesResponse)

				if err := json.Unmarshal(b, &body); err != nil {
					t.Error("Error checking order - 4 - %v", err.Error())
					t.Sleep()
					continue
				}

				t.RedirectURL = body.RedirectURL
				t.SubmitPayPal()
				return

			} else {
				t.Error("Error checking order - 5")
				t.SleepAndRotate()
				continue
			}

		case 400:
			t.Error("Error checking order - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error checking order - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error checking order - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error checking order - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error checking order - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}

func (t *Task) SubmitPayPal() {

	t.Warn("Submitting PayPal")

	for {

		req, err := http.NewRequest("GET", t.RedirectURL, nil)

		if err != nil {
			t.Error("Error submitting PayPal - 1 - %v", err.Error())
			continue
		}

		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("FF-Country", t.Country)
		req.Header.Set("FF-Currency", t.Currency)
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://www.ambushdesign.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.ambushdesign.com/")

		t.SetAllowRedirects(false)

		resp, err := t.Client.Do(req)

		if err != nil {
			t.Error("Error submitting PayPal - 2 - %v", err.Error())
			t.Rotate()
			continue
		}

		switch resp.StatusCode {
		case 302:

			if strings.Contains(resp.Header.Get("Location"), "paypal") {

				t.PayPalURL = resp.Header.Get("Location")
				t.PaypalWebhook()
				return
			}

			return
		case 400:
			t.Error("Error submitting PayPal - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		case 403:
			t.Error("Error submitting PayPal - Access Denied [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 429:
			t.Error("Error submitting PayPal - Rate limited [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		case 500, 501, 502, 503:
			t.Error("Error submitting PayPal - Site error [%v]", resp.StatusCode)
			t.Sleep()
			continue
		default:
			t.Error("Error submitting PayPal - Unhandled error [%v]", resp.StatusCode)
			t.SleepAndRotate()
			continue
		}
	}
}
