package hw10programoptimization

import (
	"fmt"
	"io"
	"io/ioutil"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(&u, domain)
}

type users [100000]User
//type users []User

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	var user User
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if err = user.UnmarshalJSON([]byte(line)); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u *users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var fullDomain string
	for _, user := range *u {
		if strings.HasSuffix(user.Email, domain) {
			fullDomain = strings.ToLower(strings.Split(user.Email, "@")[1])
			result[fullDomain] = result[fullDomain] + 1
		}
	}
	return result, nil
}
