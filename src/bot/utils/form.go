package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type Form struct {
	Values []map[string]string
}

func (form *Form) Set(key, value string) {
	form.Values = append(form.Values, map[string]string{key: value})
}

func (form *Form) Pop(keyToRemove string) {
	i := 0
	var elementToRemove string
	for _, y := range form.Values {
		for key, _ := range y {
			if keyToRemove == key {
				elementToRemove = fmt.Sprintf("%v", form.Values[i])
			}
		}
		i += 1
	}
	if elementToRemove == "" {
		return
	}
	for i := 0; i < len(form.Values); i++ {
		value := fmt.Sprintf("%v", form.Values[i])
		if value == elementToRemove {
			copy(form.Values[i:], form.Values[i+1:])
			form.Values[len(form.Values)-1] = nil
			form.Values = form.Values[:len(form.Values)-1]
			return
		}
		continue
	}
}

func (form *Form) Encode() string {
	var encoded string
	for _, pair := range form.Values {
		for key, value := range pair {
			encoded += fmt.Sprintf("%v=%v&", url.QueryEscape(key), url.QueryEscape(value))
		}
	}

	return strings.Trim(encoded, "&")
}
