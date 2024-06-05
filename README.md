# ki-call

## Prerequisite
Go 1.8 or above <br>
Python3 <br>
Make <br>

## How to start service
1. Install golang library `go mod tidy` command
2. Add proto file in `/proto` directory 
3. Generate pb golang file with command . <br> ```make gen-kitex obj={your package}``` <br> For example if your package's name is `example` . So the command is <br> `make gen-kitex obj=example`

4. Generate client with command <br> `make gen-client`
5. Build the service <br> `make build`
6. Service will start on localhost at port 9700

## List of endpoints
|endpoint|method|functionality|
|---|---|---|
|`/ls-svc`|**GET**| Get all availabe services
|`/ls-func`|**GET**| Get all available methods on a service 
|`/requests`|**GET**| Get format request from a method in a service
|`/responses`|**GET**| Get format response from a method in a service
|`/ki-call`|**POST**| Call kitex function

##  Curl Example

### get services
```curl --location 'localhost:9700/ls-svc'```

example response
```json
{
    "header": {
        "error_code": "",
        "status_code": 200
    },
    "data": {
        "list_service": [
            "YourPackage",
            "Example"
        ]
    }
}
```

### get methods
```curl --location 'localhost:9700/ls-func?service=YourPackage'```

example response
```json
{
    "header": {
        "error_code": "",
        "status_code": 200
    },
    "data": {
        "list_function": [
            "CreateYourModel",
            "UpdateYourModel"
        ]
    }
}
```

### get requests
```curl --location 'localhost:9700/requests?method=CreateYourModel&no_empty=false&service=YourPackage```

example response
```json
{
    "header": {
        "error_code": "",
        "status_code": 200
    },
    "data": {
        "method": "CreateYourModel",
        "request": {
            "param1": "",
            "param2": "",
            "param3": 0,
            "param4": []
        },
        "service": "YourPackage"
    }
}
```
### get responses
```curl --location 'localhost:9700/responses?method=CreateYourModel&no_empty=true&service=YourPackage```

example response
```json
{
    "header": {
        "error_code": "",
        "status_code": 200
    },
    "data": {
        "response": {
            "resfield1": 0,
            "resfeild2": "",
            "resfield3": []
        },
    }
}
```

### Ki-Call
```
curl --location 'localhost:9700/ki-call' \
--header 'Content-Type: application/json' \
--data '{
    "host" : "127.0.0.1:8888",
    "method": "CreateYourModel",
    "request": {
        "param1": "asd",
        "param2": 5,
        "param3": [""],
        "param 4": 100000
    },
    "service": "YourPackage"
}'
```

example respone
```json
{
    "header": {
        "error_code": "",
        "status_code": 200
    },
    "data": {
        "response": {
            "resfield1": 17,
            "resfeild2": "arfaghif",
            "resfield3": ["a","r","f","a"]
        },
    }
}
```