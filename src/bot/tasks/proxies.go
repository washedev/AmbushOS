package tasks

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ReadProxies() ([]string, error) {

	fileContent, err := ioutil.ReadFile("proxies.txt")

	if err != nil {
		return nil, err
	}

	raw := strings.Split(string(fileContent), "\r\n")

	proxies := make([]string, 0)

	for _, proxy := range raw {

		spl := strings.Split(proxy, ":")

		if len(spl) == 2 {
			proxies = append(proxies, fmt.Sprintf("http://%v:%v", spl[0], spl[1]))
		} else if len(spl) == 4 {
			proxies = append(proxies, fmt.Sprintf("http://%v:%v@%v:%v", spl[2], spl[3], spl[0], spl[1]))
		}
	}

	return proxies, nil
}
