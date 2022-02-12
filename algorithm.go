package main

import (
	"fmt"
	"sort"
)

func algorithm(
	Config Config,
	videos []*Video,
	serversMap mapOfServers,
	endpointsMap mapOfEndpoints,
	requestList []*RequestGroup,
) int {

	for _, server := range serversMap {
		videoTotlaSize := 0

		sort.Slice(server.potentialRequests, func(a int, b int) bool {
			requestA := server.potentialRequests[a]
			requestB := server.potentialRequests[b]

			dcLatencyA := requestA.endpoint.dcLatency
			dcLatencyB := requestB.endpoint.dcLatency

			cacheALatency := server.endpointLatencyMap[requestA.endpointId]
			cacheBLatency := server.endpointLatencyMap[requestB.endpointId]

			gainA := dcLatencyA - cacheALatency
			gainB := dcLatencyB - cacheBLatency

			videoASize := requestA.video.size
			videoBSize := requestB.video.size

			return gainA*requestA.nRequests-videoASize > gainB*requestB.nRequests-videoBSize
		})

		for _, potentialRequest := range server.potentialRequests {
			video := potentialRequest.video
			videoAlreadyAllocated := server.allocatedVideoMap[potentialRequest.videoId]

			if videoAlreadyAllocated || videoTotlaSize+video.size > server.serverCapacity {
				continue
			}

			// potentialRequest.endpoint.servers

			server.allocatedVideos = append(server.allocatedVideos, fmt.Sprintf("%d", video.id))
			videoTotlaSize += video.size
			server.allocatedVideoMap[potentialRequest.videoId] = true
		}
	}

	return 42
}
