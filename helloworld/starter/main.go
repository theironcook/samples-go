package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/helloworld"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	wg := sync.WaitGroup{}

	for i := 0; i < 10_000; i++ {
		wg.Add(1)
		go func() {
			workflowOptions := client.StartWorkflowOptions{
				ID:        fmt.Sprintf("hello_world_workflowID-%d", i),
				TaskQueue: "hello-world",
			}

			we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworld.Workflow, fmt.Sprintf("Temporal-%d", i))
			if err != nil {
				log.Fatalln("Unable to execute workflow", err)
			}

			log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

			// Synchronously wait for the workflow completion.
			var result string
			err = we.Get(context.Background(), &result)
			if err != nil {
				log.Fatalln("Unable get workflow result", err)
			}
			log.Println("Workflow result:", result)
			wg.Done()
		}()
	}

	log.Println("waiting for workflows to execute")
	wg.Wait()
	log.Println("all done.  check ./worker/log/perf_metrics for throughput results")
}
