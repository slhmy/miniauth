# AuthApi

All URIs are relative to *http://localhost:8080/api*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**loginPost**](#loginpost) | **POST** /login | User login|
|[**logoutPost**](#logoutpost) | **POST** /logout | User logout|
|[**meChangePasswordPut**](#mechangepasswordput) | **PUT** /me/change-password | Change user password|
|[**meGet**](#meget) | **GET** /me | Get current user|
|[**meProfilePut**](#meprofileput) | **PUT** /me/profile | Update user profile|

# **loginPost**
> HandlersLoginResponse loginPost(credentials)

Authenticate user with email and password and create a session

### Example

```typescript
import {
    AuthApi,
    Configuration,
    HandlersLoginRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let credentials: HandlersLoginRequest; //Login credentials

const { status, data } = await apiInstance.loginPost(
    credentials
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **credentials** | **HandlersLoginRequest**| Login credentials | |


### Return type

**HandlersLoginResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **logoutPost**
> { [key: string]: string; } logoutPost()

Log out the current user and destroy the session

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.logoutPost();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**{ [key: string]: string; }**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **meChangePasswordPut**
> { [key: string]: string; } meChangePasswordPut(request)

Change the current user\'s password

### Example

```typescript
import {
    AuthApi,
    Configuration,
    HandlersChangePasswordRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: HandlersChangePasswordRequest; //Change password request

const { status, data } = await apiInstance.meChangePasswordPut(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **HandlersChangePasswordRequest**| Change password request | |


### Return type

**{ [key: string]: string; }**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **meGet**
> HandlersGetUserResponse meGet()

Get the currently authenticated user\'s information from session

### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.meGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**HandlersGetUserResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **meProfilePut**
> HandlersGetUserResponse meProfilePut(request)

Update the current user\'s profile information

### Example

```typescript
import {
    AuthApi,
    Configuration,
    HandlersUpdateProfileRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: HandlersUpdateProfileRequest; //Update profile request

const { status, data } = await apiInstance.meProfilePut(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **HandlersUpdateProfileRequest**| Update profile request | |


### Return type

**HandlersGetUserResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

