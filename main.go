package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	files := []string{
		// Uncomment the line with the desired files (add other lines if needed)
		// "a",
		"a", "b", "c", "d", "e",
		// "a", "b",
		// "a", "b", "e", "f",
		// "c",
		// "d",
	}

	for _, fileName := range files {
		fmt.Printf("****************** INPUT: %s\n", fileName)
		inputSet := readFile(fmt.Sprintf("./inputFiles/%s.in", fileName))

		config, videos, serversMap, endpointsMap, requestsList := buildInput(inputSet)
		// fmt.Printf("Config %+v\n", config)
		// for i, video := range videos {
		// 	fmt.Printf("video %d %+v\n", i, video)
		// }
		// for i, server := range serversMap {
		// 	fmt.Printf("serversMap %d %+v\n", i, server)
		// }
		// for i, endpoint := range endpointsMap {
		// 	fmt.Printf("endpointsMap %d %+v\n", i, endpoint)
		// }
		// for i, req := range requestsList {
		// 	fmt.Printf("requestsList %d %+v\n", i, req)
		// }
		// printInputMetrics(input)

		algorithm(config, videos, serversMap, endpointsMap, requestsList)

		output := buildOutput(serversMap)
		// printResultMetrics(result)

		ioutil.WriteFile(fmt.Sprintf("./result/%s.out", fileName), []byte(output), 0644)
	}
}
