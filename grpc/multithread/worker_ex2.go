package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
	"errors"
)

type job struct {
	name     string
	duration time.Duration
}

func listenForJobReceive()  {
	for j := range jobsChannel {
		doWork(1, j)
	}
}

func doWork(id int, j *job) {
	fmt.Printf("worker%d: started %s, working for %f seconds\n", id, j.name, j.duration.Seconds())
	time.Sleep(j.duration)
	fmt.Printf("worker%d: completed %s!\n", id, j.name)

	if j.duration < 4000000000 {
		go func() {
			jobsFinishedChannel <- j
		}()
	} else {
		go func() {
			err := errors.New("Not found")
			errorChannel <- err
		}()
	}
}

func requestHandler(jobs chan *job, w http.ResponseWriter, r *http.Request) (*job, error) {

	job := &job{"some job", 5000000000}
	go func() {
		fmt.Printf("added: %s %s\n", job.name, job.duration)
		jobs <- job
	}()

	w.WriteHeader(http.StatusCreated)

	select {
	case jobFinished := <-jobsFinishedChannel:
		return jobFinished, nil
	case error := <-errorChannel:
		return nil, error
	}
}

var jobsChannel chan *job
var jobsFinishedChannel chan *job
var errorChannel chan error

func main() {

	var (
		maxQueueSize = flag.Int("max_queue_size", 10000, "The size of job queue")
		port         = flag.String("port", "8080", "The server port")
	)
	flag.Parse()

	jobsChannel = make(chan *job, *maxQueueSize)
	jobsFinishedChannel = make(chan *job, *maxQueueSize)
	errorChannel = make(chan error, *maxQueueSize)

	go listenForJobReceive()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		res, err := requestHandler(jobsChannel, w, r)
		if err != nil {
			fmt.Println(err)
		}
		if res != nil {
			fmt.Printf("worker: totallly completed %s!\n", res.name)
		}
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
