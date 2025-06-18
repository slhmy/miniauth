# AdminApi

All URIs are relative to *http://localhost:8080/api*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminUsersGet**](#adminusersget) | **GET** /admin/users | List all users (Admin)|
|[**adminUsersIdDelete**](#adminusersiddelete) | **DELETE** /admin/users/{id} | Delete user (Admin)|
|[**adminUsersIdGet**](#adminusersidget) | **GET** /admin/users/{id} | Get user by ID (Admin)|
|[**adminUsersIdPut**](#adminusersidput) | **PUT** /admin/users/{id} | Update user (Admin)|
|[**adminUsersIdResetPasswordPost**](#adminusersidresetpasswordpost) | **POST** /admin/users/{id}/reset-password | Reset user password (Admin)|
|[**adminUsersIdRolePut**](#adminusersidroleput) | **PUT** /admin/users/{id}/role | Update user role (Admin)|
|[**adminUsersPost**](#adminuserspost) | **POST** /admin/users | Create user (Admin)|

# **adminUsersGet**
> HandlersAdminListUsersResponse adminUsersGet()

Get a paginated list of all users in the system

### Example

```typescript
import {
    AdminApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AdminApi(configuration);

let page: number; //Page number (default: 1) (optional) (default to undefined)
let size: number; //Page size (default: 10) (optional) (default to undefined)

const { status, data } = await apiInstance.adminUsersGet(
    page,
    size
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number (default: 1) | (optional) defaults to undefined|
| **size** | [**number**] | Page size (default: 10) | (optional) defaults to undefined|


### Return type

**HandlersAdminListUsersResponse**

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminUsersIdDelete**
> { [key: string]: string; } adminUsersIdDelete()

Permanently delete a user account and all related data (admin operation)

### Example

```typescript
import {
    AdminApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AdminApi(configuration);

let id: number; //User ID (default to undefined)

const { status, data } = await apiInstance.adminUsersIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | User ID | defaults to undefined|


### Return type

**{ [key: string]: string; }**

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminUsersIdGet**
> HandlersGetUserResponse adminUsersIdGet()

Get detailed information about a specific user

### Example

```typescript
import {
    AdminApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AdminApi(configuration);

let id: number; //User ID (default to undefined)

const { status, data } = await apiInstance.adminUsersIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | User ID | defaults to undefined|


### Return type

**HandlersGetUserResponse**

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminUsersIdPut**
> HandlersAdminUserResponse adminUsersIdPut(user)

Update user information (admin operation)

### Example

```typescript
import {
    AdminApi,
    Configuration,
    HandlersAdminUpdateUserRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AdminApi(configuration);

let id: number; //User ID (default to undefined)
let user: HandlersAdminUpdateUserRequest; //User update request

const { status, data } = await apiInstance.adminUsersIdPut(
    id,
    user
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **user** | **HandlersAdminUpdateUserRequest**| User update request | |
| **id** | [**number**] | User ID | defaults to undefined|


### Return type

**HandlersAdminUserResponse**

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminUsersIdResetPasswordPost**
> { [key: string]: string; } adminUsersIdResetPasswordPost(password)

Reset a user\'s password to a new value (admin operation)

### Example

```typescript
import {
    AdminApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AdminApi(configuration);

let id: number; //User ID (default to undefined)
let password: { [key: string]: string; }; //New password

const { status, data } = await apiInstance.adminUsersIdResetPasswordPost(
    id,
    password
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **password** | **{ [key: string]: string; }**| New password | |
| **id** | [**number**] | User ID | defaults to undefined|


### Return type

**{ [key: string]: string; }**

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminUsersIdRolePut**
> HandlersAdminUserResponse adminUsersIdRolePut(role)

Update a user\'s role (admin or user)

### Example

```typescript
import {
    AdminApi,
    Configuration,
    HandlersAdminUpdateUserRoleRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AdminApi(configuration);

let id: number; //User ID (default to undefined)
let role: HandlersAdminUpdateUserRoleRequest; //User role update request

const { status, data } = await apiInstance.adminUsersIdRolePut(
    id,
    role
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **role** | **HandlersAdminUpdateUserRoleRequest**| User role update request | |
| **id** | [**number**] | User ID | defaults to undefined|


### Return type

**HandlersAdminUserResponse**

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminUsersPost**
> HandlersAdminUserResponse adminUsersPost(user)

Create a new user account (admin operation)

### Example

```typescript
import {
    AdminApi,
    Configuration,
    HandlersAdminCreateUserRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AdminApi(configuration);

let user: HandlersAdminCreateUserRequest; //User creation request

const { status, data } = await apiInstance.adminUsersPost(
    user
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **user** | **HandlersAdminCreateUserRequest**| User creation request | |


### Return type

**HandlersAdminUserResponse**

### Authorization

[BasicAuth](../README.md#BasicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

