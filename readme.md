# SAM YAML TUTORIAL

### Why use SAM?
Use when you want to define multiple interlinked pieces of your application.
Therefore, typically have multiple `Resources`
#### Remember: 
- AWS SAM defines a set of objects which can be included in a CloudFormation template to describe common components of serverless applications easily.
- SAM is based on AWS CloudFormation
- SAM is an extension of CloudFormation
- SAM is CloudFormation under the hood
- An AWS SAM template is a CloudFormation template
- Serverless == functions triggered by events

- Specify resources for Lambda application in a SAM template.
- SAM is converted to a CloudFormation template, which is then deployed
- SAM adds some resource types that are not present in CloudFormation
    - most importantly: AWS::Serverless::Function
        - builds concisely defined Lambda functions



#### Typical events: 
    - Upload to S3
    - SNS notification
    - API action

## BASIC SAM FILE:
```
AWSTemplateFormatVersion: '2010-09-09'                  // always default to these headers
Transform: AWS::Serverless-2016-10-31                   // always default to these headers - REQUIRED TO INCLUDE SAM WITHIN CLOUDFORMATION TEMPLATE (including 2016-10-31)

Resources:
    HelloWorldFunction:
        Type: AWS::Serverless::Function
        Properties:
           Handler: index.handler                       // file and function name
           Runtime: nodejs4.3                           // language spec
           CodeUri: s3://bucketName/codepackage.zip     // will pull code from this location
```

- The template above will create a Lambda, but is fundamentally incomplete since it needs somethign to trigger a Lambda function in order to run.
`Events` property defines triggers for Lambda function:

- The template above can be deployed liek any CloudFormation template: With the AWS CLI
- However, there's a better way: Using the SAM specific CLI

## DEPLOYING WITH SAM CLI:
1. `package` command
eg: `sam package --template-file template.yml --output-template-file package.yml --s3-bucket kevin-test-bucket-dispatch-deng`
`sam package` === `aws cloudformation package help`

- upload any artifacts that your lambda application requires to an AWS S bucket
- artifacts will be automatically retrieved from this S3 bucket
- artifacts include the actual lamba function
- produces output file `package.yml`
    - copy of `template.yml` but with the CodUri property of each lambda function int he template set to URL of the uplaoded package on the S3 bucket
- `package.yml` will be deployed to AWS
- can specify any bucket, and can use same bucket for multiple apps

2. `deploy` command
eg: `sam deploy --template-file package.yml --stack-name my-sam-application capabilities CAPABILITY_IAM`
`sam deploy` === `aws cloudformation deploy help`

- deploys template defined in package.yml to AWS
- creates new, or updates existing, CloudFormation stack
- `--capabilities CAPABILITY_IAM` authorizes stack to create IAM roles, which SAM applications do by default.

3. delete
- SAM cli doesn't provide command for deleting lamba apps
- can delete however by deleting the CloudFormation  stack:
`aws cloudformation delete-stack --stack-name my-sam-application`


## BUILD REAL EXAMPLE
1. `sam package --template-file template.yml --output-template-file package.yml --s3-bucket kevin-test-bucket-dispatch-deng`

2. `sam deploy --template-file package.yml --stack-name kevin-test-stack --capabilities CAPABILITY_IAM`

stack: `https://console.aws.amazon.com/cloudformation/home?region=us-east-1#/stacks/resources?stackId=arn%3Aaws%3Acloudformation%3Aus-east-1%3A176506442715%3Astack%2Fkevin-test-stack%2Fb40a4d50-87c2-11e9-a5f2-0e75601403f8`

- single function with no event source associated. Therefore, this will never be triggered by an event (or run).

3. Add `Events` to template.yml:
```
Events:
  HelloWorldApi:
  Type: Api
    Properties:
      Path: /
      Method: GET
```

- `Events` defines an event source for our Lambda function
- This will cause an API Gateway API to be created and associated with our Lambda application