# API doc

This is API documentation for List Controller. This is generated by `httpdoc`. Don't edit by hand.

## Table of contents

- [[200] GET /api/members](#200-get-apimembers)


## [200] GET /api/members

get user list

### Request



Headers

| Name  | Value  | Description |
| ----- | :----- | :--------- |
| Accept-Encoding | gzip |  |
| Authorization | admin | auth token |
| User-Agent | Go-http-client/1.1 |  |







### Response

Headers

| Name  | Value  | Description |
| ----- | :----- | :--------- |
| Content-Type | application/json |  |





Response example

<details>
<summary>Click to expand code.</summary>

```javascript
[
    {
        "id": 1,
        "name": "hoge"
    },
    {
        "id": 2,
        "name": "foo"
    },
    {
        "id": 3,
        "name": "bar"
    }
]
```

</details>



