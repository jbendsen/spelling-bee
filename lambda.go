package main

import (
	"context"
	"log"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

//aws transforms this to json
type LambdaResponse struct {
	Letters   string   `json:"letters"`
	Mandatory string   `json:"mandatory"`
	Words     []string `json:"words"`
}

type Params struct {
	Letters   string `json:"letters"`
	Mandatory string `json:"mandatory"`
}

func handleRequest(ctx context.Context, params Params) (LambdaResponse, error) {
	log.Printf("Letters: %s", params.Letters)
	log.Printf("Mandatory: %s", params.Mandatory)
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
	r := []rune(params.Mandatory)
	res, err := GetMatchingWords(params.Letters, r[0])
	if err != nil {
		return LambdaResponse{}, err
	}

	return LambdaResponse{
		Letters:   params.Letters,
		Mandatory: params.Mandatory,
		Words:     res}, nil
}

func main() {
	runtime.Start(handleRequest)
}
