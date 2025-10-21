package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func makeData() []string {
	var dbData = []string{}

	for i := range 10 {
		dbData = append(dbData, "id"+strconv.Itoa(i))
	}

	return dbData
}

func main() {
	dbData := makeData()
	t0 := time.Now()

	for _, i := range dbData {
		wg.Add(1)
		go dbCall(i)
	}
	wg.Wait()
	fmt.Printf("\nTotal execution time: %v", time.Since(t0))
}

func dbCall(i string) {
	var delay float32 = rand.Float32() * 2000
	time.Sleep(time.Duration(delay) * time.Millisecond)
	fmt.Println("The result from the database is", i)
	wg.Done()
}
