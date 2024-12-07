package main

import (
	"fmt"
	"gator/internal/config"
	"log"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)

	if err = conf.SetUser("Tinothy"); err != nil {
		log.Fatal(err)
	}

	conf, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
}
