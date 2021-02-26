package hw10programoptimization

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson/jlexer"
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
	dict, err := countDomains(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return *dict, nil
}

func countDomains(r io.Reader, domain string) (*DomainStat, error) {
	dict := make(DomainStat)
	reader := bufio.NewReader(r)
	var fullDomain string
	var user User
	var line []byte
	var err error
	for {
		line, _, err = reader.ReadLine()
		if err != nil && err != io.EOF {
			return nil, err
		}
		user.UnmarshalEasyJSON(&jlexer.Lexer{Data: line})
		if strings.HasSuffix(user.Email, domain) {
			fullDomain = strings.ToLower(strings.Split(user.Email, "@")[1])
			dict[fullDomain] ++
		}
		if err == io.EOF {
			break
		}
	}
	return &dict, nil
}

