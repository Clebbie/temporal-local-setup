package main

import (
	"context"
	"go.temporal.io/sdk/client"
)

func main() {
	c, _ := client.NewClient(client.Options{})

	createFPO(c)
	//signal(c,"ack","this is an ack")
	//signal(c, "asn", "this is an asn")
	//signal(c, "asn", "Kill")

}

func createFPO(c client.Client) {
	options := client.StartWorkflowOptions{
		ID:        "PworkflowCC",
		TaskQueue: "PTQ",
	}
	c.ExecuteWorkflow(context.Background(), options, "ParentWorkflow", "123456")
}

func signal(c client.Client, workflow string, message string) {
	switch workflow {
	case "ack":
		c.SignalWorkflow(context.Background(), "C1WF", "", workflow, message)
	case "asn":
		c.SignalWorkflow(context.Background(), "C2WF", "", workflow, message)
	}
}
