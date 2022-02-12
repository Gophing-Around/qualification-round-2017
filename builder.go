package main

import (
	"fmt"
	"strings"
)

type mapOfServers map[int]*CacheServer
type mapOfEndpoints map[int]*Endpoint

type Config struct {
	videos               int
	endpoints            int
	requestDescription   int
	cacheServers         int
	cacheServersCapacity int
}

type Video struct {
	id   int
	size int
}

type Endpoint struct {
	id        int
	dcLatency int
	nCaches   int
	servers   map[int]*CacheServer
}

type RequestGroup struct {
	videoId    int
	endpointId int
	nRequests  int

	video    *Video
	endpoint *Endpoint
}

type CacheServer struct {
	id                 int
	endpointLatencyMap map[int]int
	serverCapacity     int

	potentialRequests []*RequestGroup
	allocatedVideos   []string // list video id
	allocatedVideoMap map[int]bool
}

func buildInput(inputSet string) (
	Config,
	[]*Video,
	mapOfServers,
	mapOfEndpoints,
	[]*RequestGroup,
	[]*CacheServer,
) {

	lines := splitNewLines(inputSet)

	configLine := splitSpaces(lines[0])

	endpoints := toint(configLine[1])
	requests := toint(configLine[2])
	config := Config{
		videos:               toint(configLine[0]),
		endpoints:            endpoints,
		requestDescription:   requests,
		cacheServers:         toint(configLine[3]),
		cacheServersCapacity: toint(configLine[4]),
	}

	videoSizes := splitSpaces(lines[1])
	videos := make([]*Video, len(videoSizes))
	for id, size := range videoSizes {
		videos[id] = &Video{id: id, size: toint(size)}
	}

	serversMap := make(map[int]*CacheServer)
	serversList := make([]*CacheServer, 0)

	endpointsMap := make(map[int]*Endpoint)

	index := 2
	for i := 0; i < endpoints; i++ {
		endpointConfigLineParts := splitSpaces(lines[index])

		nCaches := toint(endpointConfigLineParts[1])
		endpoint := &Endpoint{
			id:        i,
			dcLatency: toint(endpointConfigLineParts[0]),
			nCaches:   nCaches,
			servers:   make(map[int]*CacheServer),
		}
		index++
		for k := 0; k < nCaches; k++ {
			cacheConfigParts := splitSpaces(lines[index])
			index++

			serverId := toint(cacheConfigParts[0])

			server, ok := serversMap[serverId]
			if !ok {
				server = &CacheServer{
					id:                 serverId,
					endpointLatencyMap: make(map[int]int),
					serverCapacity:     config.cacheServersCapacity,

					allocatedVideoMap: make(map[int]bool),
				}
				serversList = append(serversList, server)
			}

			server.endpointLatencyMap[i] = toint(cacheConfigParts[1])
			serversMap[serverId] = server

			endpoint.servers[serverId] = server
		}
		endpointsMap[i] = endpoint
	}

	requestsList := make([]*RequestGroup, 0)
	for i := 0; i < requests; i++ {
		requestConfigLineParts := splitSpaces(lines[index])
		index++

		videoId := toint(requestConfigLineParts[0])
		endpointId := toint(requestConfigLineParts[1])
		endpoint, _ := endpointsMap[endpointId]
		req := &RequestGroup{
			videoId:    videoId,
			video:      videos[videoId],
			endpointId: endpointId,
			endpoint:   endpoint,
			nRequests:  toint(requestConfigLineParts[2]),
		}
		requestsList = append(requestsList, req)
	}

	for _, request := range requestsList {
		for _, server := range request.endpoint.servers {
			server.potentialRequests = append(server.potentialRequests, request)
		}
	}

	return config, videos, serversMap, endpointsMap, requestsList, serversList
}

func buildOutput(servers mapOfServers) string {
	result := fmt.Sprintf("%d", len(servers))

	for _, server := range servers {
		serverLine := fmt.Sprintf("%d %s", server.id, strings.Join(server.allocatedVideos, " "))
		result += fmt.Sprintf("\n%s", serverLine)
	}
	return result
}
