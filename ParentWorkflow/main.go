package main

import (
	"fmt"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main(){
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		fmt.Printf("unable to create Temporal client %v", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, "PTQ", worker.Options{})
	w.RegisterWorkflow(ParentWorkflow)
	//w.RegisterActivity(app.ComposeGreeting)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		fmt.Println("unable to start Worker", err)
	}
}
func ParentWorkflow(ctx workflow.Context, id string) (string,error){
	var signalObj string
	isOrderFulfilled := false
	c := workflow.GetSignalChannel(ctx, "error")
	s := workflow.NewSelector(ctx)
	s.AddReceive(c, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signalObj)
		if signalObj == "Kill"{
			isOrderFulfilled = true
		}
	})
	logString := "Parent " + id
	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON, TaskQueue: "C1TQ",WorkflowID: "C1WF"})
	var cwfOne workflow.Execution
	cwfOneFuture :=workflow.ExecuteChildWorkflow(ctx, "ChildWorkflowOne", id)
	cwfOneFuture.GetChildWorkflowExecution().Get(ctx,cwfOne)
	var cwfTwo workflow.Execution
	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON, TaskQueue: "C2TQ", WorkflowID: "C2WF"})
	cwfTwoFuture := workflow.ExecuteChildWorkflow(ctx, "ChildWorkflowTwo", id)
	cwfTwoFuture.Get(ctx, cwfTwo)
	s.AddFuture(cwfTwoFuture, func(f workflow.Future) {
		var response string
		f.Get(ctx, &response)
		print(response)
		isOrderFulfilled = true
	})
	for !isOrderFulfilled{
		s.Select(ctx)
	}
	workflow.SignalExternalWorkflow(ctx, "C1WF","","kill","Kill").Get(ctx,nil)
	return logString,nil
}

