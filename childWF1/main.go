package main

import (
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"time"
)

func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		fmt.Printf("unable to create Temporal client %v", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, "C1TQ", worker.Options{})
	w.RegisterWorkflow(ChildWorkflowOne)
	// w.RegisterActivity(app.ComposeGreeting)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		fmt.Println("unable to start Worker", err)
	}
}

func ChildWorkflowOne(ctx workflow.Context, id string) (string, error) {
	var signalObj string
	receivedKillSignal := false
	c := workflow.GetSignalChannel(ctx, "ack")
	cc := workflow.GetSignalChannel(ctx, "kill")
	s := workflow.NewSelector(ctx)
	s.AddReceive(cc, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signalObj)
		if signalObj == "Kill" {
			receivedKillSignal = true
		}
	})
	s.AddReceive(c, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signalObj)
		fmt.Println("I have received a signal! ", signalObj)
	})
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{TaskQueue: "A2TQ", StartToCloseTimeout: time.Hour})

	for !receivedKillSignal {
		s.AddFuture(workflow.NewTimer(ctx, 5*time.Second), func(f workflow.Future) {
			workflow.ExecuteActivity(ctx, "RestCall", "CWF1 calling")
		})
		s.Select(ctx)
	}

	return "Done!", nil
}
