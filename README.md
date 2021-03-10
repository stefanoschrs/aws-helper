# AWS Helper

***Motivation:** Converting some generic aws scripts that I'm using all the time into a smooth and silky CLI*

## Usage
- `cp example.env .env`
- `./aws-helper -h`

## Features
### CloudFront invalidation
Use after updating the static files behind CloudFront.
+ **Use of aliases instead of long ids**
+ **Wait and notify on completion or error**
- Required IAM policies:
  - cloudfront:GetInvalidation
  - cloudfront:CreateInvalidation

<img src="https://i.imgur.com/0YuknmP.png" alt="aws-helper cloudfront invalidation example" border="0">

### ECS service deployment
Use after updating the Docker image on ECR to deploy the new service to the cluster. 
+ **Register a new Task revision**
+ **Update the cluster service to use the new Task**
- Required IAM policies:
  - TODO

<img src="https://i.imgur.com/mNvjhqh.png" alt="aws-helper ecs deployment example" border="0">
