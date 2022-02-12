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

	globalAllocationMap := make(map[int]bool)

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

			return gainA*requestA.nRequests/videoASize > gainB*requestB.nRequests/videoBSize
		})

		for _, potentialRequest := range server.potentialRequests {
			video := potentialRequest.video
			videoAlreadyAllocatedInServer := server.allocatedVideoMap[potentialRequest.videoId]
			videoAlreadyAllocatedGlobally := globalAllocationMap[potentialRequest.videoId]

			if videoAlreadyAllocatedGlobally || videoAlreadyAllocatedInServer || videoTotlaSize+video.size > server.serverCapacity {
				continue
			}

			server.allocatedVideos = append(server.allocatedVideos, fmt.Sprintf("%d", video.id))
			videoTotlaSize += video.size

			server.allocatedVideoMap[potentialRequest.videoId] = true
			globalAllocationMap[potentialRequest.videoId] = true
		}

		for _, potentialRequest := range server.potentialRequests {
			video := potentialRequest.video
			videoAlreadyAllocatedInServer := server.allocatedVideoMap[potentialRequest.videoId]

			if videoAlreadyAllocatedInServer || videoTotlaSize+video.size > server.serverCapacity {
				continue
			}

			server.allocatedVideos = append(server.allocatedVideos, fmt.Sprintf("%d", video.id))
			videoTotlaSize += video.size

			server.allocatedVideoMap[potentialRequest.videoId] = true
		}
	}

	return 42
}
