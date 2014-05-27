/* https://www.codeeval.com/open_challenges/55/
Type-ahead prediction based on n-gram frequency

*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"sort"

	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type FrequencyMap map[string]*NGramFreq

type NGramFreq struct {
	Word string
	Next map[string]int
}

func MakeMap(words *[]string, ngrams int) *FrequencyMap {
	// make our frequency map return value
	out := make(FrequencyMap)

	// how many words we want after each match
	follow_words := ngrams - 1
	for i := 0; i <= len(*words)-ngrams; i++ {
		word := (*words)[i]
		next := make([]string, follow_words)
		for ng := 0; ng < follow_words; ng++ {
			next[ng] = (*words)[i+ng+1]
		}

		//fmt.Printf("%s -> %s\n", word, next)

		if _, found := out[word]; found {
			for _, n := range next {
				if _, found := out[word].Next[n]; found {
					out[word].Next[n] += 1
				} else {
					out[word].Next[n] = 1
				}
			}
		} else {
			out[word] = &NGramFreq{
				Word: word,
				Next: make(map[string]int),
			}
			for _, n := range next {
				out[word].Next[n] = 1
			}
		}
	}
	//fmt.Printf("%s", Map)
	/*for _, w := range out {
		fmt.Printf("%s -> %s\n", w.Word, w.Next)
	}*/
	return &out
}

func (Map *FrequencyMap) Predict(needle string) string {
	// to predict the next word after needle, we find the max frequency from the Map
	// and give probability of each ngram based on the frequency from 0 to 1
	total_freq := 0
	out_pairs := make([]string, 0)
	by_value := make(map[float64][]string)
	if freq, found := (*Map)[needle]; found {
		// go through once to get the total freq of all words
		for _, v := range freq.Next {
			total_freq += v
		}
		// go through again to get relative score of each word
		for k, v := range freq.Next {
			score := float64(v) / float64(total_freq)
			by_value[score] = append(by_value[score], k)
		}
	} else {
		return ""
	}
	// now the sorting... first by highest score
	scores := make([]float64, len(by_value))
	i := 0
	for k, _ := range by_value {
		scores[i] = k
		i++
	}
	sort.Float64s(scores)                   // they're in ascending order now
	for i := len(scores) - 1; i >= 0; i-- { // go backwards to get descending order
		// now sort the matching words for this score in alpha order
		score := scores[i]
		sort.Strings(by_value[score])
		for _, s := range by_value[score] {
			out_pairs = append(out_pairs, fmt.Sprintf("%s,%.3f", s, score))
		}
	}
	return strings.Join(out_pairs, ";")
}

func main() {
	corpus :=
		`Mary had a little lamb its fleece was white as snow;
And everywhere that Mary went, the lamb was sure to go.  It followed
her to school one day, which was against the rule; It made the children
laugh and play, to see a lamb at school.  And so the teacher turned it
out, but still it lingered near, And waited patiently about till Mary
did appear.  "Why does the lamb love Mary so?" the eager children cry;
"Why, Mary loves the lamb, you know" the teacher did reply."`
	// lowercase everything
	corpus = strings.ToLower(corpus)
	//fmt.Printf("%s\n", corpus)

	var word = regexp.MustCompile(`[a-zA-Z]+`)
	matches := word.FindAllString(corpus, -1)

	//fmt.Printf("%s\n", matches)

	var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	CommandLine.Parse(os.Args[1:])

	if CommandLine.NArg() < 1 {
		panic("this program must be passed a path as the first argument to read")
	}

	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		chunks := strings.Split(scanner.Text(), ",")
		ngrams, err := strconv.Atoi(chunks[0])
		check(err)
		needle := chunks[1]

		Map := MakeMap(&matches, ngrams)
		fmt.Printf("%s\n", Map.Predict(needle))
	}
}
