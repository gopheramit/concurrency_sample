package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Data struct {
	FieldString string `json:"fieldString"`
	Number      int    `json:"number"`
	Process     bool   `json:"process"`
}

func callAPI(number int, wg *sync.WaitGroup, resultChan chan<- int) {
	defer wg.Done()

	data := Data{
		FieldString: "test",
		Number:      number,
		Process:     true,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal data:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/process", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error calling API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	var responseData Data
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return
	}

	resultChan <- responseData.Number
}

func main() {
	var wg sync.WaitGroup
	resultChan := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go callAPI(i, &wg, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []int
	for num := range resultChan {
		results = append(results, num)
	}

	fmt.Println("Accumulated results:", results)
}
