package main

import (
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		fmt.Printf("unable to create Temporal client %v", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, "Parent-Task-Queue", worker.Options{})
	w.RegisterWorkflow(ParentWorkflow)
	//w.RegisterActivity(app.ComposeGreeting)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		fmt.Println("unable to start Worker", err)
	}
}
func ParentWorkflow(ctx workflow.Context, inputString string) (string, error) {
	return inputString + " World!", nil
}
