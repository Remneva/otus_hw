package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int32

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	e := getUsers(r, domain)
	return countDomains(e, domain)
}

type users []User

func getUsers(r io.Reader, domain string) (result users) {
	scanner := bufio.NewScanner(r)
	lines := make([]string, 0)
	for scanner.Scan() {
		str := scanner.Text()
		if strings.Contains(str, domain) {
			lines = append(lines, scanner.Text())
		}
	}
	result = make(users, len(lines))
	for i, line := range lines {
		u := User{}
		if err := json.Unmarshal([]byte(line), &u); err != nil {
			return nil
		}
		result[i] = u
	}

	return result
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, user := range u {
		matched := strings.Contains(user.Email, domain)
		if matched {
			str := strings.ToLower(user.Email)
			domain := strings.SplitN(str, "@", 2)[1]
			num := result[domain]
			num++
			result[domain] = num
		}
	}
	return result, nil
}
