package main

import (
	"bytes"
	"io"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

// pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"

func main() {
	// err := exportImages("test-book2.pdf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	Annotations("test-book2.pdf")

}

func Annotations(name string) {
	pdfcpu.DecryptFile(name, "cmd", nil)
}

// func exportImages(name string) error {
// 	i, s, err := pdfcpu.ListAnnotationsFile(name, nil, nil)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(i)
// 	fmt.Println(s)

// 	return nil
// }
