package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"database/sql"
	"math/rand"
	_ "github.com/lib/pq"
)


func fetchData(db *sql.DB){
	//Using random ID for each select so as to avoid cache effects from DB 
	id := rand.Intn(9000) + 267340
	var name string
	if err := db.QueryRow("SELECT name FROM names WHERE id=$1", id).Scan(&name); err != nil {
		fmt.Println(err)
	}
}


func dummySelect(db *sql.DB){
	_, err := db.Exec("SELECT 1 FROM names")
	if err != nil{
		fmt.Println(err)
	}
}

func main(){
	var wg sync.WaitGroup
	var poolSize = 5  // --> change db pool size here
	var workerSize = 1000
	db, err := sql.Open("postgres", "postgres://theuser:thepassword@127.0.0.1/thedb")
	if err != nil{
		log.Fatal(err)
	}

	//maximum concurrent open connections to 100. 
	db.SetMaxOpenConns(poolSize)
	//max idle concurrent idle connections
	db.SetMaxIdleConns(poolSize)
	defer db.Close()


	// we are doing a dummy select here so that we can initialize the DB pool
	// before the workers start off.
	for i := 0; i < poolSize ; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			dummySelect(db)
		}()
	}
	
	wg.Wait()


	// the benchmark tests starts from here by launching num workerSize workers concurrently 
	start := time.Now()

	for i := 0; i < workerSize ; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			fetchData(db)
		}()
	}
	
	wg.Wait() // wait for all the workers to finish up

	end := time.Now()
	delta := end.Sub(start).Nanoseconds()

	fmt.Printf("Worker Size: %d \n Pool Size: %d \n Time Taken: %f ms\n", workerSize, poolSize, float64(delta)/1000000)

}

