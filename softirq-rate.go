package main

// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var rateMap *map[string][]int64 = nil
var curTimestamp = time.Time{}
var prevRateMap *map[string][]int64 = nil
var prevTimestamp = time.Time{}
var numCPU = 1 // default
var refreshInterval = 3
var pattern *regexp.Regexp = nil

func buildPattern() *regexp.Regexp {
	var pstring = "\\s+(\\w+):"
	for i := 0; i < numCPU; i++ {
		pstring += "\\s+(\\d+)"
	}
	pattern := regexp.MustCompile(pstring)
	return pattern
}

func process() {
	data, err := ioutil.ReadFile("/proc/softirqs")
	if err != nil {
		fmt.Printf("Got error when reading /proc/softirqs, err: %v", err)
		return
	}

	str := string(data)
	curTimestamp = time.Now()
	lines := strings.Split(str, "\n")
	validLines := lines[1:]
	curMap := map[string][]int64{}
	for _, line := range validLines {
		processLine(line, &curMap)
	}
	rateMap = &curMap
	result := getDiff()
	if result != nil {
		printTable(result)
	}
	prevRateMap = &curMap
	prevTimestamp = curTimestamp
}

func processLine(line string, curMap *map[string][]int64) {
	match := pattern.FindStringSubmatch(line)
	if match == nil {
		return
	}

	name := match[1]
	counters := []int64{}
	for i := 0; i < numCPU && len(match) > i+2; i++ {
		counter, err := strconv.Atoi(match[i+2])
		if err != nil {
			fmt.Printf("Got error when parsing number %v, err: %v", match[i+2], err)
		}
		counters = append(counters, int64(counter))
	}
	(*curMap)[name] = counters
}

func getDiff() *map[string][]int64 {
	if prevRateMap == nil || rateMap == nil {
		return nil
	}

	diffMap := map[string][]int64{}

	for name, counters := range *rateMap {
		prevCounters := (*prevRateMap)[name]
		diffCounters := getCounterDiff(counters, prevCounters, name)
		diffMap[name] = diffCounters
	}

	return &diffMap
}

func getCounterDiff(cur []int64, prev []int64, name string) []int64 {
	length := len(cur)
	result := []int64{}
	duration := int64(curTimestamp.Sub(prevTimestamp)) / 1000000000 // in seconds
	for i := 0; i < length; i++ {
		result = append(result, (cur[i]-prev[i])/duration)
	}
	return result
}

func clear() {
	print("\033[H\033[2J")
}

func printTable(result *map[string][]int64) {

	clear()

	fmt.Printf("Refresh Interval: every %v seconds\n", refreshInterval)

	print("\n")

	var header = fmt.Sprintf("%10s", "")

	for i := 0; i < numCPU; i++ {
		header += fmt.Sprintf("%15s", fmt.Sprintf("CPU%v", i))
	}

	header += "\n"

	fmt.Printf(header)

	keys := make([]string, 0, len(*result))

	for k := range *result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		var row = fmt.Sprintf("%10s", name)
		counters := (*result)[name]
		for _, v := range counters {
			row += fmt.Sprintf("%15s", fmt.Sprintf("%v/s", v))
		}
		row += "\n"
		fmt.Printf(row)
	}
}

func main() {
	intervalPtr := flag.Int("interval", 3, "refresh interval")

	flag.Parse()

	refreshInterval = *intervalPtr
	numCPU = runtime.NumCPU()
	pattern = buildPattern()

	duration := time.Duration(refreshInterval) * time.Second

	clear()

	process()

	for {
		time.Sleep(duration)
		process()
	}

}
