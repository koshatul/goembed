package main

import (
	"fmt"
	"log"

	"github.com/dave/jennifer/jen"
)

var data = map[string][]byte{
	"test1.dat": []byte{01, 02},
	"test2.dat": []byte{03, 04},
}

func main() {
	if v, ok := data["test1.dat"]; ok {
		log.Printf("Data: %v", v)
	}

	f := jen.NewFile("main")
	f.HeaderComment("Code generated - DO NOT EDIT.")
	f.Line()
	f.Type().Id("byteSlice").Index().Byte()
	f.Line()

	d := jen.Dict{}
	for k, v := range data {
		u := []jen.Code{}
		for _, b := range v {
			u = append(u, jen.Lit(int(b)))
		}
		d[jen.Lit(k)] = jen.Index().Byte().Values(u...)
	}
	f.Var().Id("a").Op("=").Map(jen.String()).Id("byteSlice").Values(
		d,
	)

	fmt.Printf("%#v", f)
}
