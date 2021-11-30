# vault-client.AuthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**authTokenCreatePost**](AuthApi.md#authTokenCreatePost) | **POST** /auth/token/create | The token create path is used to create new tokens.



## authTokenCreatePost

> authTokenCreatePost()

The token create path is used to create new tokens.

### Example

```javascript
import vault-client from 'hashi_corp_vault_api';

let apiInstance = new vault-client.AuthApi();
apiInstance.authTokenCreatePost((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

