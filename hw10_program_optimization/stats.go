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
	dict := make(DomainStat)
	err := countDomains(&dict, r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return dict, nil
}

func countDomains(dict *DomainStat, r io.Reader, domain string) error {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var fullDomain string
	var user User
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if err = user.UnmarshalJSON([]byte(line)); err != nil {
			return err
		}
		if strings.HasSuffix(user.Email, domain) {
			fullDomain = strings.ToLower(strings.Split(user.Email, "@")[1])
			(*dict)[fullDomain] = (*dict)[fullDomain] + 1
		}
	}
	return nil
}
