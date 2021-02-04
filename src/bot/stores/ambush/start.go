package ambush

import (
	"fmt"
	"sync"

	"bot/tasks"
)

func Start(row tasks.Row, config tasks.Config, proxies []string, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	task := &Task{
		ID:  fmt.Sprint(id),
		SKU: row.SKU,

		Email:       row.Email,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		PhoneNumber: row.PhoneNumber,
		Address1:    row.Address1,
		Address2:    row.Address2,
		State:       row.State,
		City:        row.City,
		Postcode:    row.Postcode,
		Country:     row.Country,
		CountryID:   row.CountryID,
		Currency:    row.Currency,

		Delay:   config.Delay,
		Timeout: config.Timeout,
		Webhook: config.Webhook,

		Proxies: proxies,
	}

	sizes, err := tasks.ParseSizes(row.Sizes)

	if err != nil {
		task.Error("Invalid sizes %v", row.Sizes)
	}

	task.Sizes = sizes

	task.Start()
}
