package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
)

type DomainStat map[string]int32

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	e, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(e, domain)
}

func getUsers(r io.Reader, domain string) ([]string, error) {

	scanner := bufio.NewScanner(r)
	lines := make([]string, 0, cap(scanner.Bytes()))
	for scanner.Scan() {
		str := scanner.Text()
		if strings.Contains(str, domain) {
			lines = append(lines, scanner.Text())
		}
	}
	return lines, nil

}

func countDomains(lines []string, domain string) (DomainStat, error) {
	result := make(DomainStat)
	wg := sync.WaitGroup{}

	for _, line := range lines {
		wg.Add(1)
		go func(result *DomainStat) {
			defer wg.Done()
			counter(line, *result, domain)
		}(&result)
		wg.Wait()
	}
	return result, nil
}

func counter(line string, result DomainStat, domain string) DomainStat {
	str := strings.ToLower(strings.SplitN(line, "@", 2)[1])
	email := strings.SplitN(str, "\"", 2)[0]

	matched := strings.Contains(email, domain)

	if matched {

		domain := strings.ToLower(email)
		num := result[domain]
		atomic.AddInt32(&num, 1)
		result[domain] = atomic.LoadInt32(&num)
	}
	return result
}
