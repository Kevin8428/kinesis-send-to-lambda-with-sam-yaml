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

### RESOURCES

#### S3
create S3 bucket:
`aws s3 mb s3://kevin-kinesis-to-lambda`

#### Kinesis

`zip function.zip index.js`

```
Resources:
  KinesisStream:
    Type: AWS::Kinesis::Stream
    Properties:
      ShardCount: 1
```

#### Lambda

template:
```
Resources:
  LambdaFunctionHelloWorld:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs8.10
      Timeout: 10
      Tracing: Active
      Events:
        HelloWorldApi:
          Type: Api
          Properties:
            Path: /
            Method: GET
        Stream:
          Type: Kinesis
          Properties:
            Stream: !GetAtt KinesisStream.Arn
            BatchSize: 100
            StartingPosition: LATEST
```

generates:
```
Resources:
  KinesisStream:
    Properties:
      ShardCount: 1
    Type: AWS::Kinesis::Stream
  LambdaFunctionHelloWorld:
    Properties:
      CodeUri: s3://kevin-test-bucket-dispatch-deng/6e30d9d588e0d20eb98f99bf6fb9598f
      Events:
        HelloWorldApi:
          Properties:
            Method: GET
            Path: /
          Type: Api
        Stream:
          Properties:
            BatchSize: 100
            StartingPosition: LATEST
            Stream:
              Fn::GetAtt:
              - KinesisStream -- whatever the kinesis is named, here "KinesisStream"
              - Arn
          Type: Kinesis
      Handler: index.handler
      Runtime: nodejs8.10
      Timeout: 10
      Tracing: Active
    Type: AWS::Serverless::Function

```

#### Analytics

template:
```
MyKinesisAnalytics:
  Type: AWS::KinesisAnalytics::Application
  Properties:
    ApplicationCode: string -- SQL statements that: 1. read input 2. transform it 3. generate output. 
    ApplicationDescription: string -- description of app
    ApplicationName: string -- name of analytics app
    Inputs: -- eg: configure to receive from single stream
    - NamePrefix: "sampleAppPrevix"
          InputSchema:
            RecordColumns:
              - Name: "example"
                SqlType: "INTEGER"
                Mapping: "$.example"
            RecordFormat: -- record format on the streaming source
              RecordFormatType: "JSON"
              MappingParameters:
                JSONMappingParameters:
                  RecordRowPath: "$"
```

#### IAM Role

```
KinesisAnalyticsRole:
Type: AWS::IAM::Role
Properties:
    AssumeRolePolicyDocument: -- The trust policy that is associated with this role. Defines which entities can assume the role.
    Version: "2012-10-17"
    Statement:
        - Effect: Allow
        Principal:
            Service: kinesisanalytics.amazonaws.com
        Action: "sts:AssumeRole"
    Path: "/"
    Policies:
    - PolicyName: Open
        PolicyDocument:
        Version: "2012-10-17"
        Statement:
            - Effect: Allow
            Action: "*"
            Resource: "*"
```

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