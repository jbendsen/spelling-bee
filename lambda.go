package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
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
	// request context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("REQUEST ID: %s", lc.AwsRequestID)
	// global variable
	log.Printf("FUNCTION NAME: %s", lambdacontext.FunctionName)
	r := []rune(params.Mandatory)
	//call function to get word list
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
	lambda.Start(handleRequest)
	//trimLocalFile()
}

//remove lines len>4
func trimLocalFile() {

	l, err := GetWordlistFromLocalFile()
	fmt.Println(err)

	f, _ := os.Create("./corncop_trimmed.txt")
	defer f.Close()

	for _, w := range l {
		if len(w) > 3 {
			n3, err := f.WriteString(w + "\n")
			fmt.Println(err)
			fmt.Printf("wrote %d bytes\n", n3)
		}
	}
	f.Sync()
}
