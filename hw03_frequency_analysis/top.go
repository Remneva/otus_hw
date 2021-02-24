package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

type kv struct {
	Key   string
	Value int
}

func Top10(s string) []string {
	cache := make(map[string]int)
	var text = strings.Fields(s)
	for _, s := range text {
		cache[s]++
	}

	var slice = make([]kv, 0, len(cache))
	for k, v := range cache {
		slice = append(slice, kv{k, v})
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Value > slice[j].Value
	})
	var result = make([]string, 0, len(slice))
	for i, kv := range slice {
		if i < 10 {
			result = append(result, kv.Key)
		} else {
			break
		}
	}
	return result
}
