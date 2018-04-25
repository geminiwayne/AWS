Simple Restful API
This is a simple restful API to request the database and give the response information and status.

#1.   Register an AWS account and get the role link the dynamodb database and APIgateway

#2.   In the code, set the table name of your database, and if you run in your local environment, you should get your IAM private key; if you deploy an lambda, just run it and summit these two function files: GetDevice and PostDevice by the build scripts.

-------------------------------------------------------------------------
1. run the build scrpit by terminal(if you want to deploy the Serverless, you can run the deploy.sh)(You also can run each function respectively by:GOOS=linux go build -o function.name)
2. summit the files to lambda
3. You can use Postman,Curl, APIgatway test to test this API. 
Post(with body): https://lwnw3y1crf.execute-api.us-east-1.amazonaws.com/dev/PostDevice
Get: https://lwnw3y1crf.execute-api.us-east-1.amazonaws.com/dev/{id+}