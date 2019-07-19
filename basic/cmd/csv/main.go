package main

import (
        "flag"
        "log"
        "os"

        basic "github.com/danisimba/paack/basic/internal"
)

func main() {
        filepath := flag.String("f", "", "filepath of customers csv")
        flag.Parse()

        if *filepath == "" {
                flag.PrintDefaults()
                os.Exit(2)
        }

        err := basic.ParseCSV(*filepath)
        if err != nil {
                log.Fatalln(err)
        }
}
