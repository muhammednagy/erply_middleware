# erply_middleware
Middleware service to handle authentication for erply requests.

## Description

**Server** has 1 endpoint:

* `POST` `/` - Receive parameters then authenticate it then send it to erply API then responds with erply API response 
Make sure you have request parameter otherwise it will reply with 422 unprocessable entity error

When sessionKey is not existent in redis cache then it queries the sessionKey and saves it to redis and makes 
it expire according to the sessionLength value in the response
This account access is just for experimental purposes only
```bash
REDIS_ADDRESS=localhost:6379
ERPLY_USERNAME=me@muhnagy.com
ERPLY_PASSWORD=b645636973dbf4cd985dB
ADDRESS=localhost:1232
ERPLY_CLIENT=506460
```
You can also add `REDIS_PASSWORD`

Variables could also be passed as cli parameters if you would like run with  ```-h```  to see a list of possible parameters

## Building
To build: ```make build```  
Running or testing requires the environment variables  
To run: ```make run```  
To test: ```make test```  
To run tests and display coverage percentage (70.2 % currently): ```make cover```   
To clean the binary file: ```make clean```


