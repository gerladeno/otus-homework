package hw10programoptimization

import (
	"bufio"
	"fmt"
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
	dict := make(DomainStat)
	err := countDomains(&dict, r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return dict, nil
}

func countDomains(dict *DomainStat, r io.Reader, domain string) error {
	reader := bufio.NewReader(r)
	var fullDomain string
	var user User
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err = user.UnmarshalJSON([]byte(line)); err != nil {
			return err
		}
		if strings.HasSuffix(user.Email, domain) {
			fullDomain = strings.ToLower(strings.Split(user.Email, "@")[1])
			(*dict)[fullDomain] = (*dict)[fullDomain] + 1
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

//func getUsers(r *bufio.Reader) (users, error) {
//	var user User
//	result := make([]User, 0)
//	for {
//		line, err := r.ReadString('\n')
//		if err != nil && err != io.EOF {
//			return []User{}, err
//		}
//		if err := jsoniter.Unmarshal([]byte(line), &user); err != nil {
//			return []User{}, err
//		}
//		result = append(result, user)
//		if err == io.EOF {
//			break
//		}
//	}
//	return result, nil
//}
