# spelling-bee
This is a simple program for solving New York Times' Spelling Bee Game. It's a leisure project for exploring Golang. 

# How to use

To compile: go build *.go

## To deploy to AWS
1. You must configure AWS CLI with your credentials, see [aws docs](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).
1. `./1-create-bucket.sh`
This creates an AWS S3 bucket that's needed for storing the lambda function contained in the zip file.

1. `./2-deploy.sh`
This deploys the function to AWS Lamda.

1. `./3-invoke.sh`
This invokes the lambda function using AWS CLI.

## Next step
- You may want to inspect the lamda function in the [AWS Console](https://aws.amazon.com/console). 
- Using the Console you can add an API Gateway trigger to activate the function using http.
- When deployed invoke URL: [https://....execute-api.eu-west-1.amazonaws.com/default/function-go-function-DaGnyXFmx5pd?letters=abcdefg&mandatory=a](https://xxx.execute-api.eu-west-1.amazonaws.com/default/function-go-function-DaGnyXFmx5pd?letters=abcdefg&mandatory=a)



## Stuff
AWS API Gateway Integration Request mapping template:

`{
    "letters" : "$input.params('letters')",
    "mandatory" : "$input.params('mandatory')"
}`

