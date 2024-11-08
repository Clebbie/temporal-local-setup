package main

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

func main() {
	c, _ := client.NewClient(client.Options{})

	triggerWorkflow(c)
	// signal(c, "asn", "Kill")
}

// triggerWorkflow initiates a workflow
func triggerWorkflow(c client.Client) {
	options := client.StartWorkflowOptions{
		// Name of the workflow
		ID: "MyFirstWorkflow",
		// THIS MUST MATCH THE TASK QUEUE YOUR WORKER IS LISTENING TO
		TaskQueue: "Parent-Task-Queue",
	}
	_, _ = c.ExecuteWorkflow(context.Background(), options, "ParentWorkflow", "Hello")
}

// signal This signals a workflow by ID with a string message
func signal(c client.Client, workflow string, channel string, message string) {
	fmt.Println("Signaling")
	_ = c.SignalWorkflow(context.Background(), workflow, "", channel, message)
}
