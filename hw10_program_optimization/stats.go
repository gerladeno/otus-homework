package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson/jlexer"
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
	return dict, nil
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	dict := make(DomainStat)
	reader := bufio.NewReader(r)
	var (
		fullDomain string
		user       User
		line       []byte
		err        error
	)
	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		user = User{}
		user.UnmarshalEasyJSON(&jlexer.Lexer{Data: line})
		if strings.HasSuffix(user.Email, domain) && strings.Contains(user.Email, "@") {
			fullDomain = strings.ToLower(strings.Split(user.Email, "@")[1])
			dict[fullDomain]++
		}
	}
	return dict, nil
}
