package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"fmt"
	"sort"
	"strings"
)

type kv struct {
	Key   string
	Value int
}

func Top10(s string) []string {

	var result []string
	var slice []kv
	cache := make(map[string]int)

	var text = strings.Fields(s)
	for _, s := range text {
		cache[s]++
	}
	for k, v := range cache {
		slice = append(slice, kv{k, v})
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Value > slice[j].Value
	})
	slice2 := make([]kv, 10)
	for i, kv := range slice {
		if i < 10 {
			slice2[i] = slice[i]
			result = append(result, kv.Key)
		}
	}
	for _, k := range result {
		fmt.Printf("%s\t", k)
	}
	return result
}
