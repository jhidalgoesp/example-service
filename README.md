# example-services

This is an example of a Rest API application using aws serverless, it took around 14 hours to develop:

```bash
.
├── Makefile                    <-- Make to automate tasks
├── README.md                   <-- This instructions file
├── cmd                         <-- Source code all the lambda functions
│   ├── main.go                 <-- Lambda function code
│   └── handler.go              <-- HttpHandler
├── pkg                         <-- Library packages usable by external applications
├── internal                    <-- Non shared internal code 
├── tests                       <-- Shared test mocks and assertions
└── template.yaml
```

## Example

https://wwdg2g2q44.execute-api.us-east-1.amazonaws.com/dev/api/v1/tweets?username={username}&count={count} [GET]

https://wwdg2g2q44.execute-api.us-east-1.amazonaws.com/dev/api/v1/users?id={userId} [GET, PUT]

https://wwdg2g2q44.execute-api.us-east-1.amazonaws.com/dev/api/v1/metrics [GET]

## Live Demo:

https://wwdg2g2q44.execute-api.us-east-1.amazonaws.com/dev/api/v1/tweets?username=cristiano&count=5 [GET]

https://wwdg2g2q44.execute-api.us-east-1.amazonaws.com/dev/api/v1/users?id=123 [GET, PUT]

https://wwdg2g2q44.execute-api.us-east-1.amazonaws.com/dev/api/v1/metrics [GET]

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* Optionally install godoc and staticcheck to generate documentation and lint the codebase
## Setup process

### Installing dependencies & building the target 

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

The `sam build` command is wrapped inside of the `Makefile`. To execute this simply run
 
```shell
make build
```

### Setup Parameter store

This application uses AWS System Managers Parameter Store to securely store api keys, the following parameters should be created manually:

```shell
/test-services/consumerKey          -> Twitter Consumer Key
/test-services/consumerSecret       -> Twitter Consumer Secret
/test-services/token                -> Twitter Token
/test-services/tokenSecret          -> Twitter Token Secret
```

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://127.0.0.1:3000/api/v1/tweets`

**SAM CLI** is used to emulate both Lambda and API Gateway locally and uses our `template.yaml` to understand how to bootstrap this environment (runtime, where the source code is, etc.) - The following excerpt is what the CLI will read in order to initialize an API and its routes:

```yaml
...
Events:
    GetTweetsFunction:
        Type: Api # More info about API Event Source: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
        Properties:
            Path: /api/v1/tweets
            Method: get
```

## Packaging and deployment

AWS Lambda Golang runtime requires a flat folder with the executable generated on build step. SAM will use `CodeUri` property to know where to look up for the application:

```yaml
...
    FirstFunction:
        Type: AWS::Serverless::Function
        Properties:
            CodeUri: cmd/lambdas/getTweets
            ...
```

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modifies IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

**This command will also create the dynamoDb tables needed for the application.**

You can find your API Gateway Endpoint URL in the output values displayed after deployment.

### Testing

We use `testing` package that is built-in in Golang and you can simply run the following command to run our tests:

```shell
make test
```
# Appendix

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang website: https://golang.org/doc/install

A quickstart way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```

## Generating GoDocs

To generate the project documentation run the following command:

```
make docs
```

## Database Design

```
Users Table
+--------+------+----------------+---------------+
| ID(PK) | Name | WorkExperience | TwitterHandle |
+--------+------+----------------+---------------+
|    123 | Jose | 5 years        | jhidalgoesp   |
+--------+------+----------------+---------------+

Atomic Counters Table
+---------------+---------+
|    ID(PK)     | Counter |
+---------------+---------+
| profileVisits |      20 |
+---------------+---------+

```

## Get Profile Request Concurrency

Atomic counters:

Used the UpdateItem operation to implement an atomic counter—a numeric attribute that is incremented, unconditionally, without interfering with other write requests. (All write requests are applied in the order in which they were received.) With an atomic counter, the updates are not idempotent. In other words, the numeric value increments each time you call UpdateItem.

Here i use an atomic counter to track the number of visitors to a website. In this case, the application would increment a numeric value, regardless of its current value. If an UpdateItem operation fails, the application could simply retry the operation. This would risk updating the counter twice, but we can probably tolerate a slight overcounting or undercounting of website visitors.

An atomic counter would not be appropriate where overcounting or undercounting can't be tolerated (for example, in a banking application). In this case, it is safer to use a conditional update instead of an atomic counter.


## Built With

* [SAM](https://github.com/aws/serverless-application-model) - The serverless application model
* [Go](https://golang.org/) - The Go programming language
* [Aws-Sdk-Go](https://github.com/aws/aws-sdk-go) - AWS SDK for go
* [Aws-Lambda-Go](https://github.com/aws/aws-lambda-go/) - AWS package for lambdas
* [Uber-go/zap](https://github.com/uber-go/zap/) - Structured logging
* [Oauth-v1](https://github.com/dghubble/oauth1/) - Go oauth1 library
* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)
