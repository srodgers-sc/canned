# Canned

`canned` is a REST API mocking utility allowing canned responses to be returned when a request is made to a known
endpoint. `canned` will listen on the specified port for all HTTP methods [GET, POST, PUT, PATCH, HEAD, OPTIONS,
DELETE, CONNECT, TRACE], and if a response is found for a given request endpoint, it will be returned.  Additional
features include:

* Regex endpoint match - optionally define a regex pattern to match a collection of endpoints.  Where a regex is defined
this will take precedence over the defined endpoint.
* Timeout - optionally define a time in seconds to wait before the response is sent.

Responses can be preloaded from file on application start or uploaded via the following endpoints:

* `/canned/upload` - suitable for single and multiple response upload at runtime.
* `/canned/upload/file` - suitable for mass response upload e.g. pre-loading responses to facilitate a testing cycle.

## Executing the utility

Execute without responses preloaded:
```
$ canned -port=:5555
```
Execute with responses preloaded:
```
$ canned -port=:5555 -responses=responses.json
```
## Response Format
```
{
  "responses":[
    {
      "endpoint": "<endpoint-to-respond-to>",
      "regex":"<optional-endpoint-pattern-match>",
      "method": "<http-method>",
      "code": "<http-code>",
      "headers": {
        "Content-Type": "application/json",
        "<other-key": "<other-value"
      },
      "body": "<body-in-required-format>",
      "timeout":"<optional-timeout-before-responding>"
    }
  ]
}
```

## Examples
```
{
  "responses":[
    {
      "endpoint":"/oauth2/token",
      "method":"POST",
      "code":"200",
      "headers":{
        "Content-Type":"application/json"
      },
      "body":"{\"access_token\":\"token\",\"expires_in\":3600,\"token_type\":\"type\"}"
    },
    {
      "endpoint":"/",
      "method":"GET",
      "code":"200",
      "headers":{
        "Content-Type":"application/xml"
      },
      "body":"<?xml version=\"1.0\" encoding=\"UTF-8\"?><AssumeRoleWithClientGrantsResponse xmlns=\"https://sts.amazonaws.com/doc/2011-06-15/\"><AssumeRoleWithClientGrantsResult><Credentials><AccessKeyId>FNBJHYPJYYC1UEQTPZJ2</AccessKeyId><SecretAccessKey>YJcjm1c3Crj9C6Chyk7Eg7kxnX0FdpKMq7nBkhoR</SecretAccessKey><Expiration>2019-07-02T18:31:51Z</Expiration><SessionToken>eyJhbGciOiJIUzUxMik2ips2Lg</SessionToken></Credentials></AssumeRoleWithClientGrantsResult></AssumeRoleWithClientGrantsResponse>"
    },
    {
      "endpoint":"/bin5/root/driver/search/df.docx",
      "code":"200",
      "method": "HEAD",
      "headers":{
        "Content-Type":"application/text",
        "X-Amz-Meta-Department":"22",
        "X-Amz-Meta-Firstname":"Ben",
        "X-Amz-Meta-Last-Name":"Vasquez",
        "X-Amz-Meta-Uuid":"df99234b-7ff3-4a14-9917-d4c4a7b02876",
        "X-Amz-Request-Id":"15ADD8EA152284DA",
        "Last-Modified": "Wed, 21 Oct 2015 07:28:00 GMT"
      },
      "body":""
    },
    {
      "endpoint":"/bin1/",
      "regex":"/bin1/",
      "method": "PUT",
      "code":"200",
      "headers":{
        "Content-Type":"application/text",
        "Last-Modified": "Wed, 21 Oct 2015 07:28:00 GMT"
      },
      "body":""
    },
    {
      "endpoint":"/oauth2/jwks",
      "method": "GET",
      "code":"200",
      "headers":{
        "Content-Type":"application/json"
      },
      "body":"{\"keys\":[{\"alg\":\"RS256\",\"e\":\"AQAB\",\"kid\":\"ZGMyZjRlM2U2OWNjMTExMmU3ZGRmNzk5NjNhZTBhMmNlYmE0YTZhNw\",\"kty\":\"RSA\",\"n\":\"xoQ4-zVEStdycv7FtFIBMYGEm_wMAyDndL04E-D2hMW0hvfRDGhSYgs-qQ4e5LZHJ2J74ZJgAonu_wO9kj4YVbrl5GSBcKHHZELza9sCtVdwNMvO0bCfcX1WG3A5qI0d0xXUm2AWpeTETyWZ8xKzVr_oRnlM8wotq6Q-1jM8SS8_o6xjXfDFP9izHDSVRa1BmOn9efyXCDTufna-HDZrtrktWdLVT74lfXXEFtmqttE-lGqmdoNoyw0pOcJQhW_gi5RyuLsdQ8CwgC2F2n0W4fK_x8OkJ5vB0hqwTMHDD-yxeZJruO6Ke3o5BRAatEm_cVEutCq_dpsyDCeRvirPPw\",\"use\":\"sig\"}]}",
      "timeout":"10"
    },
  ]
}
```
1. The first example is of a response to a request for an identity token.  The expected request will be for a `POST`
method to endpoint `/oauth2/token`.  The response issued will be a `200` with the content type header and body
in JSON format.

2. The second example serves all `GET` requests to endpoint `/` and will return a `200` with content type header and
body in XML format.

3. This example is for a `HEAD` request to the endpoint `/bin5/root/driver/search/df.docx`, returning a zero length
body and numerous headers.

4. This example shows the use of regex.  The response is for all `PUT` requests to eny endpoint matching the pattern
`/bin1/`.  So all endpoints beginning with `/bin1/` for method `PUT` will use this response.

5. The last example shows the use of the timeout field.  In this example, all `GET` requests to endpoint `oauth2/jwks`
will wait `10` seconds before sending the defined JSON response.

## Load a Response

To define a response for a given endpoint e.g. `/v1/messages`:
```
curl --header "Content-Type: application/json" --request POST --data '{"endpoint": "/v1/messages","code": "200", "method": "POST", headers": {"Content-Type": "application/json"},"body": "{\"id\":\"Y2lzY29zcGFS0xMWUI1YmUz\",\"Email\":\"person@email.com\",\"text\":\"this is a test\",\"created\":\"2019-02-25T08:15:35.029Z\"}"}' 127.0.0.1:5555/canned/upload
```
## Load Responses from File
To upload a file containing many responses:
```
curl -X POST --header 'Content-Type: multipart/form-data' --header 'Accept: application/json' -F responses=@"response.json" 127.0.0.1:5555/canned/upload/file
```
---
**NOTE**

Adding a response which has already been defined will overwrite the existing response.
___

## Make a Request
To see the response from a request to endpoint `/v1/messages`:
```
curl --request GET -i 127.0.0.1:5555/v1/messages
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 01 Jul 2019 20:44:00 GMT
Content-Length: 120

{"id":"Y2lzY29zcGFS0xMWUI1YmUz","Email":"person@email.com","text":"this is a test","created":"2019-02-25T08:15:35.029Z"}
```

