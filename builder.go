package main

type mapOfServers map[int]*Server
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
	servers   map[int]*Server
}

type RequestGroup struct {
	videoId    int
	endpointId int
	nRequests  int

	video    *Video
	endpoint *Endpoint
}

type Server struct {
	id                 int
	endpointLatencyMap map[int]int
}

func buildInput(inputSet string) (Config, []*Video, mapOfServers, mapOfEndpoints, []*RequestGroup) {

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

	serversMap := make(map[int]*Server)

	endpointsMap := make(map[int]*Endpoint)

	index := 2
	for i := 0; i < endpoints; i++ {
		endpointConfigLineParts := splitSpaces(lines[index])

		nCaches := toint(endpointConfigLineParts[1])
		endpoint := &Endpoint{
			id:        i,
			dcLatency: toint(endpointConfigLineParts[0]),
			nCaches:   nCaches,
			servers:   make(map[int]*Server),
		}
		index++
		for k := 0; k < nCaches; k++ {
			cacheConfigParts := splitSpaces(lines[index])
			index++

			serverId := toint(cacheConfigParts[0])

			server, ok := serversMap[serverId]
			if !ok {
				server = &Server{
					id:                 serverId,
					endpointLatencyMap: make(map[int]int),
				}
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

	return config, videos, serversMap, endpointsMap, requestsList
}

func buildOutput(result int) string {
	return "42"
}
