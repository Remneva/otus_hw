//go:generate ffjson $GOFILE
package hw10_program_optimization //nolint:golint,stylecheck

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"sync/atomic"
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
type emails []string

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	e, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(e, domain)
}

func getUsers(r io.Reader) (emails, error) {

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	size := len(lines)

	wg := sync.WaitGroup{}
	wg.Add(size)
	user := &User{}

	result := make(emails, 0)
	rw := sync.RWMutex{}

	for _, line := range lines {
		line := []byte(line)

		go func(result *emails) {
			defer wg.Done()
			getEmails(&rw, *user, line, result)
		}(&result)
	}
	wg.Wait()
	return result, nil
}

func countDomains(e emails, domain string) (DomainStat, error) {
	result := make(DomainStat)
	wg := sync.WaitGroup{}

	for _, email := range e {
		wg.Add(1)
		go func(result *DomainStat) {
			defer wg.Done()

			matcher(email, *result, domain)
		}(&result)
		wg.Wait()
	}
	return result, nil
}

func matcher(email string, result DomainStat, domain string) DomainStat {
	matched := strings.Contains(email, domain)

	if matched {
		domain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
		num := result[domain]
		atomic.AddInt32(&num, 1)
		result[domain] = atomic.LoadInt32(&num)
	}
	return result
}

func getEmails(rw sync.Locker, user User, line []byte, result *emails) emails {

	if err := ffjson.UnmarshalFast(line, &user); err != nil {
		return *result
	}
	rw.Lock()
	*result = append(*result, user.Email)
	rw.Unlock()
	return *result
}
