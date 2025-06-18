# ServiceInternalApplicationCreateRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**client_id** | **string** | Custom client ID | [default to undefined]
**client_secret** | **string** | Custom client secret | [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**name** | **string** |  | [default to undefined]
**redirect_uris** | **Array&lt;string&gt;** |  | [default to undefined]
**scopes** | **Array&lt;string&gt;** |  | [optional] [default to undefined]
**trusted** | **boolean** |  | [optional] [default to undefined]
**website** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { ServiceInternalApplicationCreateRequest } from './api';

const instance: ServiceInternalApplicationCreateRequest = {
    client_id,
    client_secret,
    description,
    name,
    redirect_uris,
    scopes,
    trusted,
    website,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
