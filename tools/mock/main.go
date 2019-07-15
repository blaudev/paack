package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	numRegisters := flag.String("r", "", "Number of registers")
	filepath := flag.String("f", "", "Customer filepath")
	flag.Parse()

	if *numRegisters == "" {
		log.Fatalln("Number of registers is required:\n    importer -r 1000000")
	}

	num, err := strconv.Atoi(*numRegisters)
	if err != nil {
		log.Fatalln(err)
	}

	if *filepath == "" {
		log.Fatalln("Customers filepath is required:\n    importer -f filepath.csv")
	}

	f, err := os.Create(*filepath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	for i := 1; i <= num; i++ {
		buf.WriteString(fmt.Sprintf("\"%d\",\"Paack SPV Investments, SL.\",\"info@paack.co\"\n", i))
	}

	_, err = f.WriteString(buf.String())
	if err != nil {
		log.Fatalln(err)
	}
}
