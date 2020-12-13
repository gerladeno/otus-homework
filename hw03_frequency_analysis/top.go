package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

func Top10(input string) []string {
	lowerCleaned := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(input), "\n", " "), "\t", " ")
	reStr := `(([^A-Za-zА-Яа-я0-9]-[^A-Za-zА-Яа-я0-9])|[,. !?;:—'"()@+<>\[\]{}\\|/*&#$^%~_=])+`
	re := regexp.MustCompile(reStr)

	words := re.Split(lowerCleaned, -1)

	wcDict := make(map[string]int)
	for _, word := range words {
		if _, ok := wcDict[word]; ok {
			wcDict[word]++
		} else if word != "" {
			wcDict[word] = 1
		}
	}

	if len(wcDict) == 0 {
		return nil
	}

	type wordCount struct {
		word  string
		count int
	}

	var wc = make([]wordCount, 0)
	for word, count := range wcDict {
		wc = append(wc, wordCount{word, count})
	}
	sort.Slice(wc, func(i, j int) bool {
		return wc[i].count > wc[j].count
	})

	var result []string
	for i := 0; i < 10; i++ {
		result = append(result, wc[i].word)
		if i+1 == len(wc) {
			break
		}
	}
	return result
}
