package basic

import (
        "database/sql"
        "encoding/csv"
        "io"
        "net/rpc"
        "os"
        "strconv"

        _ "github.com/lib/pq"
)

const (
        postgresConnectionString = "host=localhost port=5432 dbname=postgres user=postgres password=123456 sslmode=disable"
        integratorHost           = "localhost:5002"
)

// ParseCSV parses the csv file
func ParseCSV(filepath string) error {
        db, err := sql.Open("postgres", postgresConnectionString)
        if err != nil {
                return err
        }
        defer db.Close()

        f, err := os.Open(filepath)
        if err != nil {
                return err
        }
        defer f.Close()

        cl, err := rpc.DialHTTP("tcp", integratorHost)
        if err != nil {
                return err
        }

        rd := csv.NewReader(f)
        for {
                data, err := rd.Read()
                if err == io.EOF {
                        // end of file/end of process
                        return nil
                }

                if err != nil {
                        return err
                }

                c, err := parseCustomer(data)
                if err != nil {
                        return err
                }

                err = saveToDatabase(db, c)
                if err != nil {
                        return err
                }

                resp, err := sendToCMS(cl, c)
                if err != nil {
                        return err
                }

                if !resp.Ok {
                        return err
                }
        }
}

func parseCustomer(d []string) (Customer, error) {
        c := Customer{}

        id, err := strconv.Atoi(d[0])
        if err != nil {
                return c, err
        }

        c.ID = id
        c.Name = d[1]
        c.Email = d[2]

        return c, nil
}

func saveToDatabase(db *sql.DB, c Customer) error {
        _, err := db.Exec(`INSERT INTO customers (ID, Name, Email, Status) VALUES ($1, $2, $3, $4);`, c.ID, c.Name, c.Email, c.Status)
        return err
}

func sendToCMS(cl *rpc.Client, c Customer) (*Response, error) {
        req := &Request{
                c,
        }

        resp := &Response{}

        err := cl.Call(integratorProcessFuncName, req, resp)
        if err != nil {
                return nil, err
        }

        return resp, nil
}
