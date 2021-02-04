package main

import (
	"bot/log"
	"bot/stores/ambush"
	"bot/tasks"
	"fmt"
	"os"
	"sync"
)

var (
	format = fmt.Sprintf
)

func main() {

	log.Infoln(format("TraianBOT"), "-")
	log.Infoln(format("Welcome user.."), "-")

	config, err := tasks.ReadConfig()

	if err != nil {
		log.Error("Error reading config file", "-")
		os.Exit(0)
	}

	proxies, err := tasks.ReadProxies()

	if err != nil {
		log.Error("Error reading proxies file", "-")
		os.Exit(0)
	}

	rows, err := tasks.ReadFile("tasks.csv")

	if err != nil {
		fmt.Println(err)
		log.Error("Error reading tasks file", "-")
		os.Exit(0)
	}

	if len(rows) > 500 {
		log.Error("Max 500 tasks", "-")
		os.Exit(0)
	}

	wg := sync.WaitGroup{}

	for id, row := range rows {
		wg.Add(1)
		go ambush.Start(row, config, proxies, id, &wg)
	}

	wg.Wait()

}
