/* https://www.codeeval.com/open_challenges/31/
Find the rightmost occurrence of a character in a string per line of
file given as the first arg

"Hello World,r" would be 8
"Hello CodeEval,E" would be 10

*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

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
		chunks := strings.Split(scanner.Text(), ",")
		fmt.Printf("%d\n", strings.LastIndex(chunks[0], chunks[1]))
	}
}
