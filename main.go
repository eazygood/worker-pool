package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eazygood/worker-pool/internal"
)

const (
	// define worker counts if needed
	workerCount = 15
)

func getDomains() ([][]string, error) {
	f, err := os.Open("top-1m.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	execGetHttpFn := func(ctx context.Context, arg interface{}) (interface{}, error) {
		domain := fmt.Sprintf("http://%v", arg.(string))
		start := time.Now()
		resp, err := http.Get(domain)
		end := time.Since(start)

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		return internal.HttpResponseMetrics{
			Label:    "GET",
			Domain:   domain,
			Duration: end,
			SizeInKb: float64(len(body) / (1 << 10)),
		}, nil
	}

	domains, err := getDomains()

	if err != nil {
		log.Fatal(err)
	}

	jobCount := len(domains)
	jobs := make([]internal.Job, jobCount)

	for i := 0; i < jobCount; i++ {
		jobs[i] = internal.Job{
			ID:       i,
			ExecFunc: execGetHttpFn,
			Args:     domains[i][1],
		}
	}

	wp := internal.New(workerCount)
	ctx, cancel := context.WithCancel(context.TODO())

	defer cancel()
	go wp.GenerateJobBulk(jobs)
	go wp.Run(ctx)

	for {
		select {
		case result, ok := <-wp.Results():
			if !ok {
				continue
			}

			if result.Err != nil {
				fmt.Printf("error orrured: %v", result.Err.Error())
				continue
			}

			fmt.Println(result.Value.(internal.HttpResponseMetrics))

		case <-wp.Done:
			return
		}
	}
}
