package hw10programoptimization

import (
	"bufio"
	"fmt"
	jsoniter "github.com/json-iterator/go"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	reader := bufio.NewReader(r)
	u, err := getUsers(reader)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r *bufio.Reader) (users, error) {
	var user User
	result := make([]User, 0)
	for {
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			return []User{}, err
		}
		if err := jsoniter.Unmarshal([]byte(line), &user); err != nil {
			return []User{}, err
		}
		result = append(result, user)
		if err == io.EOF {
			break
		}
	}
	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, user := range u {
		if strings.HasSuffix(user.Email, domain) {
			fullDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[fullDomain] = result[fullDomain] + 1
		}
	}
	return result, nil
}
