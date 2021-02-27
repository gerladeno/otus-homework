// +build !bench

package hw10programoptimization

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Auxiliary struct to imitate a very long input without using much memory
type Streamer struct {
	Input    chan byte
	lastByte *byte
}

func (s *Streamer) Read(p []byte) (n int, err error) {
	if s.lastByte != nil {
		p[n] = *s.lastByte
		n++
	}
	if c := cap(p); c > 0 {
		for b := range s.Input {
			p[n] = b
			n++
			if n == c {
				break
			}
		}
	}
	if i, ok := <-s.Input; ok {
		s.lastByte = &i
		err = nil
		return
	}
	err = io.EOF
	return
}

func NewStreamer() *Streamer {
	return &Streamer{make(chan byte), nil}
}

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("very long input", func(t *testing.T) {
		dataSource := NewStreamer()
		count := 20000
		go func() {
			for i := 0; i < count; i++ {
				for _, b := range []byte(data) {
					dataSource.Input <- b
				}
				dataSource.Input <- '\n'
			}
			close(dataSource.Input)
		}()
		result, err := GetDomainStat(dataSource, "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": count}, result)
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString("{}"), "whatever")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
	
	t.Run("broken input", func(t *testing.T) {
		
	})
}

func TestStreamer(t *testing.T) {
	t.Run("simple check", func(t *testing.T) {
		s := "sample string\nwith two lines"
		dataSource := NewStreamer()
		count := 2
		go func() {
			for i := 0; i < count; i++ {
				for _, b := range []byte(s) {
					dataSource.Input <- b
				}
			}
			close(dataSource.Input)
		}()
		result, err := ioutil.ReadAll(dataSource)
		require.NoError(t, err)
		require.Equal(t, string(result), strings.Repeat(s, count))
	})

	t.Run("input longer than buffer size", func(t *testing.T) {
		s := "sample string\nwith two lines"
		dataSource := NewStreamer()
		count := 200
		go func() {
			for i := 0; i < count; i++ {
				for _, b := range []byte(s) {
					dataSource.Input <- b
				}
			}
			close(dataSource.Input)
		}()
		result, err := ioutil.ReadAll(dataSource)
		require.NoError(t, err)
		require.Equal(t, string(result), strings.Repeat(s, count))
	})
}
