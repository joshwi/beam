package utils

import (
	"fmt"
	"math"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

var a0 = regexp.MustCompile(`[^a-zA-Z\d\/\-]+`)
var a1 = regexp.MustCompile(`\_\-|\-\_`)
var a2 = regexp.MustCompile(`\_{2,}`)

func FormatPath(filepath string) string {

	ext := path.Ext(filepath)

	dest := strings.ReplaceAll(filepath, ext, "")

	dest = a0.ReplaceAllString(dest, "_")
	dest = a1.ReplaceAllString(dest, "-")
	dest = a2.ReplaceAllString(dest, "_")

	path := dest + ext

	return path
}

func BuildRequests(query map[string]string, urls []string) ([]string, string) {

	// Enter variables to the url templates
	req_urls := []string{}
	keys := []string{}
	for _, url := range urls {
		for k, v := range query {
			keys = append(keys, k)
			re, _ := regexp.Compile(fmt.Sprintf("{%v}", k))
			url = re.ReplaceAllString(url, v)
		}
		req_urls = append(req_urls, url)
	}

	// Sort list of keys
	sort.Slice(keys, func(a, b int) bool {
		return keys[a] < keys[b]
	})
	// Get list of values http_channel alpha order of keys
	values := []string{}
	for _, v := range keys {
		values = append(values, query[v])
	}
	// Compute label for the collection
	label := strings.Join(values, "_")

	return req_urls, label
}

func ComputeMetrics(pass int, total int) string {

	rate := "0%"

	if total > 0 {
		percent := (float64(pass) / float64(total)) * 100.0
		rate = fmt.Sprintf("%v%%", math.Round(percent*100)/100)
	}

	return rate
}

func ComputeTime(total int, start time.Time, end time.Time) (string, string) {
	elapsed := end.Sub(start)
	duration := fmt.Sprintf("%v", elapsed.Round(time.Second/1000))

	average := "0 ms"

	if total > 0 {
		average = fmt.Sprintf("%v ms", int(elapsed.Milliseconds())/total)
	}

	return duration, average
}

