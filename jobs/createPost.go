package jobs

import (
	"fmt"
	"sync"
	"time"
)

type createPost struct {
	Name string
}

var Time = time.Now()

func NewCreatePostJob() *createPost {
	return &createPost{
		Name: "createPost",
	}
}

func (j *createPost) Run() error {
	fmt.Println("Job is Running", time.Since(Time))
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Wait()
	return nil
}
