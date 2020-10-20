# AWS Helper

***Motivation:** Converting some generic aws scripts that I'm using all the time into a smooth and silky CLI*

## Usage
- `cp .example.env .env`
- `./aws-helper -h`

## Features
### CloudFront invalidation
+ **Use of aliases instead of long ids**
+ **Wait and notify on completion or error**
- Required IAM policies:
  - cloudfront:GetInvalidation
  - cloudfront:CreateInvalidation

<img src="https://i.ibb.co/X5Kk1zv/aws-helper-cloudfront-example.png" alt="aws-helper-cloudfront-example" border="0">
