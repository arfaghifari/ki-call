# ki-call

## Prerequisite
Go 1.8 or above <br>
Python3 <br>
Make <br>

## How to start service
1. Add proto file in `/proto` directory 
2. Generate pb golang file with command . <br> ```make gen-kitex obj={your package}``` <br> For example if your package's name is `example` . So the command is <br> `make gen-kitex obj=example`

3. Generate client with command <br> `make gen-client`
4. Build the service <br> `make build`
5. Service will start on localhost at port 9700

## List of endpoints
|endpoint|method|functionality|
|---|---|---|
|`/ls-svc`|**GET**| Get all availabe services
|`/ls-func`|**GET**| Get all available methods on a service 
|`/ls-requests`|**GET**| Get format request from a method in a service
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
        "param3": [],
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
            "resfield1": 0,
            "resfeild2": "",
            "resfield3": []
        },
    }
}
```