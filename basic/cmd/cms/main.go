package main

import (
        "log"

        basic "github.com/danisimba/paack/basic/internal"
)

func main() {
        err := basic.CmsService()
        if err != nil {
                log.Fatalln(err)
        }
}
