package main

import (
	"fmt"
	"sort"
)

type GroupedPotentialRequest struct {
	score   int
	videoId int
	video   *Video
}

func algorithm(
	Config Config,
	videos []*Video,
	serversMap mapOfServers,
	endpointsMap mapOfEndpoints,
	requestList []*RequestGroup,
	serversList []*CacheServer,
) int {

	globalAllocationMap := make(map[int]int)

	sort.Slice(serversList, func(a, b int) bool {
		serverA := serversList[a]
		serverB := serversList[b]

		potReqA := 0
		for _, req := range serverA.potentialRequests {
			potReqA += req.nRequests
		}

		potReqB := 0
		for _, req := range serverB.potentialRequests {
			potReqB += req.nRequests
		}

		return potReqA > potReqB
	})

	for _, server := range serversList {
		videoTotlaSize := 0

		groupedPotentialRequests := make([]*GroupedPotentialRequest, 0)
		groupedPotentialRequestsMap := make(map[int]*GroupedPotentialRequest)

		for _, potRequest := range server.potentialRequests {
			videoId := potRequest.videoId

			group, ok := groupedPotentialRequestsMap[videoId]
			cacheLatency := server.endpointLatencyMap[potRequest.endpointId]
			dcLatenty := potRequest.endpoint.dcLatency

			gain := dcLatenty - cacheLatency
			score := (gain * potRequest.nRequests) / potRequest.video.size

			if !ok {
				// ccrea
				newGroup := &GroupedPotentialRequest{
					score:   score,
					videoId: videoId,
					video:   potRequest.video,
				}
				groupedPotentialRequests = append(groupedPotentialRequests, newGroup)
				groupedPotentialRequestsMap[videoId] = newGroup
				continue
			}

			// patch group
			group.score = (group.score + score) / len(server.potentialRequests)
		}

		sort.Slice(groupedPotentialRequests, func(a int, b int) bool {
			requestA := groupedPotentialRequests[a]
			requestB := groupedPotentialRequests[b]
			return requestA.score > requestB.score
		})
		for _, potentialRequest := range groupedPotentialRequests {
			video := potentialRequest.video
			videoAlreadyAllocatedInServer := server.allocatedVideoMap[potentialRequest.videoId]
			_, videoAlreadyAllocatedGlobally := globalAllocationMap[potentialRequest.videoId]

			if videoAlreadyAllocatedInServer || videoTotlaSize+video.size > server.serverCapacity {
				continue
			}

			if videoAlreadyAllocatedGlobally { // && globalLatencyPerRequest < gainedLatencyPerRequest {
				continue
			}

			score := potentialRequest.score
			server.allocatedVideos = append(server.allocatedVideos, fmt.Sprintf("%d", video.id))
			videoTotlaSize += video.size

			server.allocatedVideoMap[potentialRequest.videoId] = true
			globalAllocationMap[potentialRequest.videoId] = score
		}
		for _, potentialRequest := range groupedPotentialRequests {
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

func algorithm2(
	Config Config,
	videos []*Video,
	serversMap mapOfServers,
	endpointsMap mapOfEndpoints,
	requestList []*RequestGroup,
	serversList []*CacheServer,
) int {

	globalAllocationMap := make(map[int]int)

	sort.Slice(serversList, func(a, b int) bool {
		serverA := serversList[a]
		serverB := serversList[b]

		potReqA := 0
		for _, req := range serverA.potentialRequests {
			potReqA += req.nRequests
		}

		potReqB := 0
		for _, req := range serverB.potentialRequests {
			potReqB += req.nRequests
		}

		return potReqA > potReqB
	})

	for _, server := range serversList {
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
			_, videoAlreadyAllocatedGlobally := globalAllocationMap[potentialRequest.videoId]

			if videoAlreadyAllocatedInServer || videoTotlaSize+video.size > server.serverCapacity {
				continue
			}

			cacheALatency := server.endpointLatencyMap[potentialRequest.endpointId]
			gainedLatencyPerRequest := (potentialRequest.endpoint.dcLatency - cacheALatency) * potentialRequest.nRequests
			if videoAlreadyAllocatedGlobally { // && globalLatencyPerRequest < gainedLatencyPerRequest {
				continue
			}

			server.allocatedVideos = append(server.allocatedVideos, fmt.Sprintf("%d", video.id))
			videoTotlaSize += video.size

			server.allocatedVideoMap[potentialRequest.videoId] = true
			globalAllocationMap[potentialRequest.videoId] = gainedLatencyPerRequest
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
