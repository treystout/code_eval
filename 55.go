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

	var re = regexp.MustCompile(`[a-z]+`)
	words := re.FindAllString(corpus, -1)

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
		needle := strings.ToLower(chunks[1])
		needle_words := strings.Split(needle, " ")
		needle_len := len(needle_words)

		total_freq := 0
		freq := make(map[string]int)
		for idx := 0; idx <= len(words)-needle_len; idx++ {
			this_gram := strings.Join(words[idx:idx+needle_len], " ")
			if this_gram == needle {
				follow := strings.Join(words[idx+needle_len:idx+ngrams], " ")
				total_freq++
				if _, found := freq[follow]; found {
					freq[follow]++
				} else {
					freq[follow] = 1
				}
			}
		}
		out_pairs := make([]string, 0)
		by_value := make(map[float64][]string)
		// go through again to get relative score of each word
		for k, v := range freq {
			score := float64(v) / float64(total_freq)
			by_value[score] = append(by_value[score], k)
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
		if len(out_pairs) > 0 {
			fmt.Printf("%s\n", strings.Join(out_pairs, ";"))
		} else {
			fmt.Printf("no matches for %s, %s\n", needle, ngrams)
		}
	}
}
