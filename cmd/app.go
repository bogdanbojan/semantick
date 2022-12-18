package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/bogdanbojan/semantick/pkg/epub"
)

func main() {
	bk, err := epub.Open("book_test.epub")
	if err != nil {
		log.Fatal(err)
	}
	defer bk.Close()

	// log.Printf("files: %+v", bk.Files())
	s, _ := json.MarshalIndent(bk.Opf.Manifest[2].ID, "", "\t")
	log.Println(string(s))

	r, err := bk.Open("9780698163744_EPUB-13.xhtml")
	if err != nil {
		log.Fatal(err)
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, r)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
}
