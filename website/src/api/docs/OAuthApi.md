# OAuthApi

All URIs are relative to *http://localhost:8080/api*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminOauthApplicationsGet**](#adminoauthapplicationsget) | **GET** /admin/oauth/applications | List OAuth Applications|
|[**adminOauthApplicationsIdDelete**](#adminoauthapplicationsiddelete) | **DELETE** /admin/oauth/applications/{id} | Delete OAuth Application|
|[**adminOauthApplicationsIdPut**](#adminoauthapplicationsidput) | **PUT** /admin/oauth/applications/{id} | Update OAuth Application|
|[**adminOauthApplicationsIdTogglePost**](#adminoauthapplicationsidtogglepost) | **POST** /admin/oauth/applications/{id}/toggle | Toggle OAuth Application Status|
|[**adminOauthApplicationsIdToggleTrustedPost**](#adminoauthapplicationsidtoggletrustedpost) | **POST** /admin/oauth/applications/{id}/toggle-trusted | Toggle OAuth Application Trusted Status|
|[**adminOauthApplicationsPost**](#adminoauthapplicationspost) | **POST** /admin/oauth/applications | Create OAuth Application|
|[**adminOauthInternalApplicationsBatchPost**](#adminoauthinternalapplicationsbatchpost) | **POST** /admin/oauth/internal/applications/batch | Batch Internal Create OAuth Applications|
|[**adminOauthInternalApplicationsPost**](#adminoauthinternalapplicationspost) | **POST** /admin/oauth/internal/applications | Internal Create OAuth Application|
|[**oauthAuthorizeGet**](#oauthauthorizeget) | **GET** /oauth/authorize | OAuth Authorization|
|[**oauthAuthorizePost**](#oauthauthorizepost) | **POST** /oauth/authorize | OAuth Authorization Decision|
|[**oauthTokenPost**](#oauthtokenpost) | **POST** /oauth/token | OAuth Token|
|[**oauthUserinfoGet**](#oauthuserinfoget) | **GET** /oauth/userinfo | OAuth User Info|

# **adminOauthApplicationsGet**
> Array<ServiceApplicationResponse> adminOauthApplicationsGet()

Get all OAuth applications (admin only)

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

const { status, data } = await apiInstance.adminOauthApplicationsGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<ServiceApplicationResponse>**

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
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOauthApplicationsIdDelete**
> adminOauthApplicationsIdDelete()

Delete an OAuth application (admin only)

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let id: number; //Application ID (default to undefined)

const { status, data } = await apiInstance.adminOauthApplicationsIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Application ID | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOauthApplicationsIdPut**
> ServiceApplicationResponse adminOauthApplicationsIdPut(request)

Update an OAuth application (admin only)

### Example

```typescript
import {
    OAuthApi,
    Configuration,
    ServiceApplicationCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let id: number; //Application ID (default to undefined)
let request: ServiceApplicationCreateRequest; //Application data

const { status, data } = await apiInstance.adminOauthApplicationsIdPut(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **ServiceApplicationCreateRequest**| Application data | |
| **id** | [**number**] | Application ID | defaults to undefined|


### Return type

**ServiceApplicationResponse**

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
|**403** | Forbidden |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOauthApplicationsIdTogglePost**
> adminOauthApplicationsIdTogglePost()

Activate or deactivate an OAuth application (admin only)

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let id: number; //Application ID (default to undefined)

const { status, data } = await apiInstance.adminOauthApplicationsIdTogglePost(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Application ID | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOauthApplicationsIdToggleTrustedPost**
> adminOauthApplicationsIdToggleTrustedPost()

Toggle the trusted status of an OAuth application (admin only)

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let id: number; //Application ID (default to undefined)

const { status, data } = await apiInstance.adminOauthApplicationsIdToggleTrustedPost(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Application ID | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | No Content |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOauthApplicationsPost**
> ServiceApplicationResponse adminOauthApplicationsPost(request)

Create a new OAuth application (admin only)

### Example

```typescript
import {
    OAuthApi,
    Configuration,
    ServiceApplicationCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let request: ServiceApplicationCreateRequest; //Application data

const { status, data } = await apiInstance.adminOauthApplicationsPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **ServiceApplicationCreateRequest**| Application data | |


### Return type

**ServiceApplicationResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOauthInternalApplicationsBatchPost**
> { [key: string]: any; } adminOauthInternalApplicationsBatchPost(request)

Create multiple OAuth applications with custom client_id and secret (internal API, requires internal token or admin session)

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let request: Array<ServiceInternalApplicationCreateRequest>; //Array of application data with custom credentials
let authorization: string; //Internal token (Bearer <token> or Internal <token>) (optional) (default to undefined)
let xInternalToken: string; //Internal token (optional) (default to undefined)
let internalToken: string; //Internal token (optional) (default to undefined)

const { status, data } = await apiInstance.adminOauthInternalApplicationsBatchPost(
    request,
    authorization,
    xInternalToken,
    internalToken
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **Array<ServiceInternalApplicationCreateRequest>**| Array of application data with custom credentials | |
| **authorization** | [**string**] | Internal token (Bearer &lt;token&gt; or Internal &lt;token&gt;) | (optional) defaults to undefined|
| **xInternalToken** | [**string**] | Internal token | (optional) defaults to undefined|
| **internalToken** | [**string**] | Internal token | (optional) defaults to undefined|


### Return type

**{ [key: string]: any; }**

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
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOauthInternalApplicationsPost**
> ServiceApplicationResponse adminOauthInternalApplicationsPost(request)

Create an OAuth application with custom client_id and secret (internal API, requires internal token or admin session)

### Example

```typescript
import {
    OAuthApi,
    Configuration,
    ServiceInternalApplicationCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let request: ServiceInternalApplicationCreateRequest; //Application data with custom credentials
let authorization: string; //Internal token (Bearer <token> or Internal <token>) (optional) (default to undefined)
let xInternalToken: string; //Internal token (optional) (default to undefined)
let internalToken: string; //Internal token (optional) (default to undefined)

const { status, data } = await apiInstance.adminOauthInternalApplicationsPost(
    request,
    authorization,
    xInternalToken,
    internalToken
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **ServiceInternalApplicationCreateRequest**| Application data with custom credentials | |
| **authorization** | [**string**] | Internal token (Bearer &lt;token&gt; or Internal &lt;token&gt;) | (optional) defaults to undefined|
| **xInternalToken** | [**string**] | Internal token | (optional) defaults to undefined|
| **internalToken** | [**string**] | Internal token | (optional) defaults to undefined|


### Return type

**ServiceApplicationResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **oauthAuthorizeGet**
> oauthAuthorizeGet()

Start OAuth authorization flow

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let responseType: string; //Response type (must be \'code\') (default to undefined)
let clientId: string; //OAuth client ID (default to undefined)
let redirectUri: string; //Redirect URI (default to undefined)
let scope: string; //Requested scopes (space-separated) (optional) (default to undefined)
let state: string; //State parameter for CSRF protection (optional) (default to undefined)
let codeChallenge: string; //PKCE code challenge (optional) (default to undefined)
let codeChallengeMethod: string; //PKCE code challenge method (optional) (default to undefined)

const { status, data } = await apiInstance.oauthAuthorizeGet(
    responseType,
    clientId,
    redirectUri,
    scope,
    state,
    codeChallenge,
    codeChallengeMethod
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **responseType** | [**string**] | Response type (must be \&#39;code\&#39;) | defaults to undefined|
| **clientId** | [**string**] | OAuth client ID | defaults to undefined|
| **redirectUri** | [**string**] | Redirect URI | defaults to undefined|
| **scope** | [**string**] | Requested scopes (space-separated) | (optional) defaults to undefined|
| **state** | [**string**] | State parameter for CSRF protection | (optional) defaults to undefined|
| **codeChallenge** | [**string**] | PKCE code challenge | (optional) defaults to undefined|
| **codeChallengeMethod** | [**string**] | PKCE code challenge method | (optional) defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**302** | Redirect to authorization page or back to client |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **oauthAuthorizePost**
> { [key: string]: string; } oauthAuthorizePost(request)

Handle user\'s authorization decision

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let request: { [key: string]: any; }; //Authorization decision

const { status, data } = await apiInstance.oauthAuthorizePost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **{ [key: string]: any; }**| Authorization decision | |


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

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **oauthTokenPost**
> ServiceTokenResponse oauthTokenPost()

Exchange authorization code or refresh token for access token

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

let grantType: string; //Grant type (authorization_code or refresh_token) (default to undefined)
let clientId: string; //OAuth client ID (default to undefined)
let code: string; //Authorization code (required for authorization_code grant) (optional) (default to undefined)
let redirectUri: string; //Redirect URI (required for authorization_code grant) (optional) (default to undefined)
let clientSecret: string; //OAuth client secret (optional) (default to undefined)
let codeVerifier: string; //PKCE code verifier (optional) (default to undefined)
let refreshToken: string; //Refresh token (required for refresh_token grant) (optional) (default to undefined)

const { status, data } = await apiInstance.oauthTokenPost(
    grantType,
    clientId,
    code,
    redirectUri,
    clientSecret,
    codeVerifier,
    refreshToken
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **grantType** | [**string**] | Grant type (authorization_code or refresh_token) | defaults to undefined|
| **clientId** | [**string**] | OAuth client ID | defaults to undefined|
| **code** | [**string**] | Authorization code (required for authorization_code grant) | (optional) defaults to undefined|
| **redirectUri** | [**string**] | Redirect URI (required for authorization_code grant) | (optional) defaults to undefined|
| **clientSecret** | [**string**] | OAuth client secret | (optional) defaults to undefined|
| **codeVerifier** | [**string**] | PKCE code verifier | (optional) defaults to undefined|
| **refreshToken** | [**string**] | Refresh token (required for refresh_token grant) | (optional) defaults to undefined|


### Return type

**ServiceTokenResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **oauthUserinfoGet**
> { [key: string]: any; } oauthUserinfoGet()

Get user information using access token

### Example

```typescript
import {
    OAuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OAuthApi(configuration);

const { status, data } = await apiInstance.oauthUserinfoGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**{ [key: string]: any; }**

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

