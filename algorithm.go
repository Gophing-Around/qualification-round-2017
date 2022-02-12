package main

import "fmt"

func algorithm(
	Config Config,
	videos []*Video,
	serversMap mapOfServers,
	endpointsMap mapOfEndpoints,
	requestList []*RequestGroup,
) int {

	for _, server := range serversMap {
		videoTotlaSize := 0
		for _, video := range videos {
			if videoTotlaSize+video.size > server.serverCapacity {
				continue
			}
			server.allocatedVideos = append(server.allocatedVideos, fmt.Sprintf("%d", video.id))
			videoTotlaSize += video.size
		}
	}

	return 42
}
