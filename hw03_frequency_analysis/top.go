package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"fmt"
	"sort"
	"strings"
)

func Top10(s string) []string {

	var result []string
	cache := make(map[string]int)

	type kv struct {
		Key   string
		Value int
	}

	var text = strings.Fields(s) // делим строку на саб строки по пробелам и кладем в слайс
	for _, s := range text {     // кладем их в мапку, ключив мапе уникальны, если значения дублируются, value увеличивается
		cache[s]++
	}

	var slice []kv
	for k, v := range cache { // из мапы перекладываем в слайс для соритровки
		slice = append(slice, kv{k, v})
	}

	sort.Slice(slice, func(i, j int) bool { // сортируем слайс
		return slice[i].Value > slice[j].Value
	})

	slice2 := make([]kv, 10, 10) // создаем слайс с длиной 10

	for i, kv := range slice { // берем первые 10 ключей из первого слайса
		if i < 10 {
			slice2[i] = slice[i]
			result = append(result, kv.Key)
		}

	}
	for _, k := range result { // проверка результата
		fmt.Printf("%s\t", k)
	}

	return result
}
