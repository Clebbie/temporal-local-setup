package main
import(
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)


func main(){
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		fmt.Printf("unable to create Temporal client %v", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, "A2TQ", worker.Options{})
	w.RegisterActivity(RestCall)
	//w.RegisterActivity(app.ComposeGreeting)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		fmt.Println("unable to start Worker", err)
	}
}

func RestCall(ctx context.Context, request string) (string,error){
	return fmt.Sprintf("%s activity 2", request),nil
}