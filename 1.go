/* https://www.codeeval.com/open_challenges/1/
normal fizzbuzz question but the factors and limit are definied in a text file
like so...

3 5 10
2 7 15
5 6 30

Where 3 and 5 are factors to look for and 10 is the upper limit to count to
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func FizzBuzz(multipleA int, multipleB int, limit int) {
	// counts up from 1 to limit, printing 'F' if limit is divisible by multipleA
	// 'B' if divisible by multipleB
	if limit < 0 {
		panic("limit must be greater than 0!")
	}
	// a place to hold output
	var out []string
	for i := 1; i <= limit; i++ {
		if i%multipleA == 0 && i%multipleB == 0 {
			out = append(out, "FB")
		} else if i%multipleA == 0 {
			out = append(out, "F")
		} else if i%multipleB == 0 {
			out = append(out, "B")
		} else {
			out = append(out, fmt.Sprint(i))
		}
	}
	fmt.Printf("%s\n", strings.Join(out, " "))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
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
		chunks := strings.Split(scanner.Text(), " ")

		first, err := strconv.Atoi(chunks[0])
		check(err)
		second, err := strconv.Atoi(chunks[1])
		check(err)
		limit, err := strconv.Atoi(chunks[2])
		check(err)

		FizzBuzz(first, second, limit)
	}
}
