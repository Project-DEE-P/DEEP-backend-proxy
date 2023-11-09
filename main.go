package main

import (
	"DEEP-backend-hmux/laundry"
	"flag"
	"fmt"
	"log"
	"strconv"
)

var port string

func init() {
	flagPort := flag.Int("p", 8080, "Enter the port")
	flag.Parse()
	port = ":" + strconv.Itoa(*flagPort)
}

func main() {
	fmt.Println("Listen proxy server...  ", port)
	log.Fatal(laundry.New().Run(port))
}
