package main

import (
	"context"
	"log"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Params struct {
	Letters   string
	Mandatory string
}

func handleRequest(ctx context.Context, event Params) (string, error) {
	// event
	//eventJson, _ := json.MarshalIndent(event, "", "  ")
	//log.Printf("EVENT: %s", string(event))
	// environment variables
	log.Printf("Letters: %s", event.Letters)
	log.Printf("Mandatory: %s", event.Mandatory)
	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	log.Println("ALL ENV VARS:")
	for _, element := range os.Environ() {
		log.Println(element)
	}
	// request context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("REQUEST ID: %s", lc.AwsRequestID)
	// global variable
	log.Printf("FUNCTION NAME: %s", lambdacontext.FunctionName)
	// context method
	deadline, _ := ctx.Deadline()
	log.Printf("DEADLINE: %s", deadline)
	// AWS SDK call
	r := []rune(event.Mandatory)
	res := GetMatchingWordsResponse(event.Letters, r[0])
	return res, nil
}

func main() {
	runtime.Start(handleRequest)
}
