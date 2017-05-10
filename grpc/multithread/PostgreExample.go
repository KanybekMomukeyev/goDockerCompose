/*
        Original idea from
        http://www.acloudtree.com/how-to-shove-data-into-postgres-using-goroutinesgophers-and-golang/
*/
package main

import (
"log"
"time"
"fmt"
"database/sql"
_ "github.com/lib/pq"
)

func main() {
	//var ch chan int
	fmt.Printf("StartTime: %v\n", time.Now())
	var (
		sStmt string = "insert into test (gopher_id, created) values ($1, $2)"
		gophers int = 1
		entries int = 10
	)

	finishChan := make(chan int)

	for i := 0; i < gophers; i++ {
		go func(c chan int) {
			db, err := sql.Open("postgres", "host=localhost dbname=template1 sslmode=disable")
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			stmt, err := db.Prepare(sStmt)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			for j := 0; j < entries; j++ {
				res, err := stmt.Exec(j, time.Now())
				if err != nil || res == nil {
					log.Fatal(err)
				}
			}

			c <- 1
		}(finishChan)
	}

	finishedGophers := 0
	finishLoop := false
	for {
		if finishLoop {
			break
		}
		select {
		case n := <-finishChan:
			finishedGophers += n
			if finishedGophers == 10 {
				finishLoop = true
			}
		}
	}

	fmt.Printf("StopTime: %v\n", time.Now())
}