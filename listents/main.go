package main

import (
	"os"
	"fmt"
	"github.com/20kdc/go-vkv"
	"io/ioutil"
)

// go run ./listents $P2/sdk_content/maps/preview.vmf

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("listents expects one argument: target file\n")
	} else {
		inf, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("failed open! %s\n", err)
			return
		}
		defer inf.Close()
		data, err := ioutil.ReadAll(inf)
		if err != nil {
			fmt.Printf("failed read! %s\n", err)
			return
		}
		tkns, err := kvkv.InTokenize(string(data), true, "")
		if err != nil {
			fmt.Printf("failed tokenize! %s\n", err)
			return
		}
		obj, err := kvkv.InParse(tkns)
		if err != nil {
			fmt.Printf("failed parse! %s\n", err)
			return
		}
		for _, v := range obj {
			if v.Key == "entity" {
				ent, err := v.ValueObject()
				if err != nil {
					fmt.Printf("value of entity issue: %s\n", err)
					return
				}
				// ignore these errors
				tne, _ := ent.FindString("targetname")
				cne, _ := ent.FindString("classname")
				fmt.Printf("entity %s : %s\n", tne, cne)
			}
		}
	}
}
