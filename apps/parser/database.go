package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

const (
	databaseConcurrency      = 100
	databaseRecordsPerCycle  = 10000
	postgresConnectionString = "host=localhost port=5432 dbname=postgres user=postgres password=123456 sslmode=disable"
	sqlCmd                   = `INSERT INTO customers (ID, Name, Email, Status) VALUES ($1, $2, $3, 'to_add')
        ON CONFLICT (ID) DO UPDATE
        SET Name = EXCLUDED.Name, Email = EXCLUDED.Email, Status = 'to_update'
        RETURNING Status;`
)

type database struct {
	sync.Mutex
	postgres *sql.DB
}

type databaseResp struct {
	err error
	cs  []Customer
}

func newDatabase() *database {
	return &database{}
}

func (db *database) open() error {
	pg, err := sql.Open("postgres", postgresConnectionString)
	if err != nil {
		return err
	}

	db.postgres = pg

	return nil
}

func (db *database) saveCustomers(cs []Customer) ([]Customer, error) {
	size := len(cs) / databaseRecordsPerCycle
	if len(cs)%size != 0 {
		size++
	}

	jobs := make(chan []Customer, size)
	resp := make(chan databaseResp, size)

	for i := 0; i < databaseConcurrency; i++ {
		go db.worker(jobs, resp)
	}

	for i := 0; i < size; i++ {
		cst := cs[i*databaseRecordsPerCycle : (i+1)*databaseRecordsPerCycle]
		if i+1 == size {
			cst = cs[i*databaseRecordsPerCycle:]
		}

		jobs <- cst
	}

	csr := make([]Customer, 0, databaseRecordsPerCycle)
	for i := 0; i < size; i++ {
		r := <-resp
		if r.err != nil {
			return nil, r.err
		}

		csr = append(csr, r.cs...)
	}

	return csr, nil
}

func (db *database) worker(inCh <-chan []Customer, outCh chan<- databaseResp) {
	for cs := range inCh {
		outCh <- db.job(cs)
	}
}

func (db *database) job(cs []Customer) databaseResp {
	log.Printf("Sending %d customers to database\n", len(cs))

	tx, err := db.postgres.Begin()
	if err != nil {
		return databaseResp{err: err}
	}

	stmt, err := tx.Prepare(sqlCmd)
	if err != nil {
		return databaseResp{err: err}
	}
	defer stmt.Close()

	csr := make([]Customer, 0, len(cs))
	for _, c := range cs {
		err = stmt.QueryRow(c.ID, c.Name, c.Email).Scan(&c.Status)
		if err != nil {
			return databaseResp{err: err}
		}

		csr = append(csr, c)
	}

	err = tx.Commit()
	if err != nil {
		return databaseResp{err: err}
	}

	return databaseResp{cs: csr}
}

func (db *database) close() error {
	return db.postgres.Close()
}
