package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"menteslibres.net/gosexy/to"
)

type Cache struct {
	id      int64
	latency int64
}

type EndPoint struct {
	latency    int64
	num_caches int64
	caches     []Cache
}

type Request struct {
	num_requests int64
	video        int64
	endpoint     int64
}

type AlgoCache struct {
	id        int64
	storage   int64
	video     []int64
	endpoints []CacheEndpoint
}

type CacheEndpoint struct {
	latencySaved int64
	videos       []AlgoVideo
	endpoint     AlgoEndpoint
}

type AlgoVideo struct {
	id           int64
	size         int64
	rating       int64
	requestCount int64
}

type AlgoEndpoint struct {
	id          int64
	baseLatency int64
	caches      map[int64]int64 // cacheId | latency
	video       []AlgoVideo
}

var algoCaches []AlgoCache
var algoEndpoints []AlgoEndpoint

var num_videos int64
var num_endpoints int64
var num_requests int64
var num_caches int64
var cache_size int64
var videos []int64
var endpoints []EndPoint
var requests []Request

var res_lines []string

func parseCache(lines []string, offset int64, ep *EndPoint) {
	parts := strings.Split(lines[offset], " ")
	c := Cache{to.Int64(parts[0]), to.Int64(parts[1])}
	ep.caches = append(ep.caches, c)
}

func parseEndpoint(lines []string, offset int64) int64 {
	parts := strings.Split(lines[offset], " ")
	ep := EndPoint{to.Int64(parts[0]), to.Int64(strings.TrimSpace(parts[1])), make([]Cache, 0)}

	var x int64

	if ep.num_caches > 0 {
		for x = 0; x < ep.num_caches; x++ {
			parseCache(lines, offset+x+1, &ep)
		}
	}
	endpoints = append(endpoints, ep)

	return ep.num_caches
}

func readVideos(lines []string) {
	strs := strings.Split(lines[1], " ")

	for i := 0; i < len(strs); i++ {
		videos = append(videos, to.Int64(strings.TrimSpace(strs[i])))
	}
}

func readRequest(lines []string) {
	for i := 0; i < len(lines)-1; i++ {
		options := strings.Split(lines[i], " ")
		requests = append(requests, Request{to.Int64(options[2]), to.Int64(options[0]), to.Int64(options[1])})
	}
}

func readFile(filename string) {
	var lines []string

	content, err := ioutil.ReadFile(filename)

	if err != nil {
		print("Failed to open ", filename, err)
		os.Exit(1)
	}

	lines = strings.Split(string(content), "\n")

	nrs := strings.Split(lines[0], " ")
	num_videos = to.Int64(nrs[0])
	num_endpoints = to.Int64(nrs[1])
	num_requests = to.Int64(nrs[2])
	num_caches = to.Int64(nrs[3])
	cache_size = to.Int64(nrs[4])

	var offset int64 = 2
	var x int64

	readVideos(lines)

	// Parse the endpoints and their
	for x = 0; x < num_endpoints; x++ {
		offset += parseEndpoint(lines, offset) + 1
	}

	readRequest(lines[offset:])
}

func writeOutput(filename string) {
	tempFilename := filename + ".out"

	file, err := os.OpenFile(tempFilename, os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		println("Failed to open ", tempFilename, err.Error())
		os.Exit(1)
	}

	for _, line := range res_lines {
		file.WriteString(line)
		file.WriteString("\n")
	}
}

func convert() {
	algoCaches = []AlgoCache{}
	var id int64
	for id = 0; id < num_caches; id++ {
		algoCache := new(AlgoCache)
		algoCache.id = id
		algoCache.storage = cache_size
		algoCache.video = []int64{}
		algoCaches = append(algoCaches, *algoCache)
	}

	algoEndpoints = make([]AlgoEndpoint, len(endpoints))
	for i, e := range endpoints {
		ep := new(AlgoEndpoint)
		ep.id = to.Int64(i)
		ep.baseLatency = e.latency
		ep.caches = make(map[int64]int64, e.num_caches)
		for _, c := range e.caches {
			ep.caches[c.id] = c.latency
		}
		algoEndpoints[i] = *ep
	}

	for _, req := range requests {
		vid := new(AlgoVideo)
		vid.id = req.video
		vid.requestCount = req.num_requests
		vid.size = videos[req.video]

		algoEndpoints[req.endpoint].video = append(algoEndpoints[req.endpoint].video, *vid)
	}
}

func createResult() {
	cachesAffected := 0
	for _, c := range algoCaches {
		if len(c.video) > 0 {
			cachesAffected++
		}
	}
	res_lines = make([]string, cachesAffected+1)
	res_lines[0] = to.String(cachesAffected)
	pos := 1
	for id, c := range algoCaches {
		if len(c.video) > 0 {
			vids := to.String(id) + " "
			for _, v := range c.video {
				vids += to.String(v) + " "
			}
			res_lines[pos] = vids
			pos++
		}
	}
}

func magic() {
	convert()
	loopEndpoints()
	findLonlies()
	realMagic()
	createResult()
}

func main() {
	if len(os.Args) == 1 {
		println("Usage: go run magic.go filename")
		os.Exit(1)
	}

	filename := os.Args[1]

	readFile(filename)
	magic()
	writeOutput(filename)
}

// Algo

func loopEndpoints() {
	for _, endpoint := range algoEndpoints {
		for cacheID, ls := range endpoint.caches {
			cache := algoCaches[cacheID]
			cacheEp := new(CacheEndpoint)
			cacheEp.latencySaved = endpoint.baseLatency - ls
			cacheEp.endpoint = endpoint
			cacheEp.videos = make([]AlgoVideo, len(endpoint.video))
			copy(cacheEp.videos, endpoint.video)
			for _, vid := range cacheEp.videos {
				vid.rating = vid.requestCount * cacheEp.latencySaved
			}

			cache.endpoints = append(cache.endpoints, *cacheEp)
			algoCaches[cacheID] = cache
		}
	}
}

func realMagic() {
	for id, cache := range algoCaches {
		for _, endp := range cache.endpoints {
			endp.videos = QuickSort(endp.videos)
			pos := 0
			for cache.storage > 0 && pos < len(endp.videos) {
				vid := endp.videos[pos]
				if cache.storage-vid.size > 0 {
					cache.video = append(cache.video, vid.id)
					cache.storage -= vid.size
					for cid := range endp.endpoint.caches {
						epCaches := algoCaches[cid]
						for _, epCEndp := range epCaches.endpoints {
							if epCEndp.endpoint.id == endp.endpoint.id {
								delPos := 0
								for index, epceVid := range epCEndp.videos {
									index = index - delPos
									if epceVid.id == vid.id {
										epCEndp.videos = remove(epCEndp.videos, index)
										delPos++
									}
								}
							}
						}
					}
				}
				pos++
			}
		}
		algoCaches[id] = cache
	}
}

func findLonlies() {
	for _, cache := range algoCaches {
		if len(cache.endpoints) == 1 {
			endp := cache.endpoints[0]
			endp.videos = QuickSort(endp.videos)
			pos := 0
			for cache.storage > 0 && pos < len(endp.videos) {
				vid := endp.videos[pos]
				if cache.storage-vid.size > 0 {
					cache.video = append(cache.video, vid.id)
					cache.storage -= vid.size
					for cid := range endp.endpoint.caches {
						epCaches := algoCaches[cid]
						for _, epCEndp := range epCaches.endpoints {
							if epCEndp.endpoint.id == endp.endpoint.id {
								for index, epceVid := range epCEndp.videos {
									if epceVid.id == vid.id {
										epCEndp.videos = remove(epCEndp.videos, index)
									}
								}
							}
						}
					}
				}
				pos++
			}
		}
	}
}

func remove(slice []AlgoVideo, s int) []AlgoVideo {
	return append(slice[:s], slice[s+1:]...)
}

func QuickSort(slice []AlgoVideo) []AlgoVideo {
	length := len(slice)

	if length <= 1 {
		sliceCopy := make([]AlgoVideo, length)
		copy(sliceCopy, slice)
		return sliceCopy
	}

	m := slice[rand.Intn(length)]

	less := make([]AlgoVideo, 0, length)
	middle := make([]AlgoVideo, 0, length)
	more := make([]AlgoVideo, 0, length)

	for _, item := range slice {
		switch {
		case item.rating < m.rating:
			less = append(less, item)
		case item.rating == m.rating:
			middle = append(middle, item)
		case item.rating > m.rating:
			more = append(more, item)
		}
	}

	less, more = QuickSort(less), QuickSort(more)

	less = append(less, middle...)
	less = append(less, more...)

	return less
}
