package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/velosypedno/information-coding-systems/arithmetic-coding/internal"
)

func main() {
	f, err := os.Create("tst.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	a := 'a'
	writer := bufio.NewWriter(f)
	writer.WriteRune(a)
	writer.Flush()

	f, err = os.Open("tst.txt")
	if err != nil {
		log.Fatal(err)
	}
	st := internal.NewInputStream(f)
	for {
		b, err := st.ReadBit()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d", b)
	}
}
