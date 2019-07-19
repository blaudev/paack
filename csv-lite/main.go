package main

import (
        "database/sql"
        "encoding/csv"
        "flag"
        "io"
        "log"
        "net/rpc"
        "os"
        "strconv"

        _ "github.com/lib/pq"
)

const (
        postgresConnectionString = "host=localhost port=5432 dbname=postgres user=postgres password=123456 sslmode=disable"
        cmsHost                  = "localhost:5002"
        cmsProcessFuncName       = "Service.Process"
)

// Customer is the domain model
type Customer struct {
        ID     int
        Name   string
        Email  string
        Status string
}

// Request represents the request to the RPC server
type Request struct {
        Customer Customer
}

// Response represents the response from the RPC server
type Response struct {
        Ok bool
}

func main() {
        filepath := flag.String("f", "", "options of customers csv")
        flag.Parse()

        if *filepath == "" {
                flag.PrintDefaults()
                os.Exit(2)
        }

        db, err := sql.Open("postgres", postgresConnectionString)
        if err != nil {
                log.Fatalln(err)
        }
        defer db.Close()

        f, err := os.Open(*filepath)
        if err != nil {
                log.Fatalln(err)
        }
        defer f.Close()

        cl, err := rpc.DialHTTP("tcp", cmsHost)
        if err != nil {
                log.Fatalln(err)
        }

        rd := csv.NewReader(f)
        for {
                data, err := rd.Read()
                if err == io.EOF {
                        // end of file/end of process
                        break
                }

                if err != nil {
                        log.Fatalln(err)
                }

                c := Customer{}

                id, err := strconv.Atoi(data[0])
                if err != nil {
                        log.Fatalln(err)
                }

                c.ID = id
                c.Name = data[1]
                c.Email = data[2]

                _, err = db.Exec(`INSERT INTO customers (ID, Name, Email, Status) VALUES ($1, $2, $3, $4);`, c.ID, c.Name, c.Email, c.Status)
                if err != nil {
                        log.Fatalln(err)
                }

                req := &Request{
                        c,
                }

                resp := &Response{}

                err = cl.Call(cmsProcessFuncName, req, resp)
                if err != nil {
                        log.Fatalln(err)
                }
        }
}
