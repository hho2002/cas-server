package validators

import (
    "fmt"
    "errors"
    "sort"
    "net/http"
    
    "github.com/matthewvalimaki/cas-server/types"
)

// ValidateProxyGrantingURL validates proxy pgtURL
func ValidateProxyGrantingURL(config *types.Config, serviceID string, pgtURL string) *types.CasError {
    // make sure we have the pgtURL specified as service in the first place
    proxyServices, ok := config.FlatServiceIDList[pgtURL]
    if !ok {
        return &types.CasError{Error: errors.New("`pgtUrl` was not found from configuration"), CasErrorCode: types.CAS_ERROR_CODE_UNAUTHORIZED_SERVICE_PROXY}
    }
    
    // 
    services, ok := config.FlatServiceIDList[serviceID]
    if !ok {
        return &types.CasError{Error: errors.New("`service` was not found from configuration"), CasErrorCode: types.CAS_ERROR_CODE_INVALID_SERVICE}
    }

    match := false
    
    for _, service := range services {
        for _, proxyService := range proxyServices {    
            // make sure the pgtURL is found with the service
            sort.Strings(config.Services[service].ProxyServices)
            
            i := sort.SearchStrings(config.Services[service].ProxyServices, proxyService)
            if i < len(config.Services[service].ProxyServices) && config.Services[service].ProxyServices[i] == proxyService {
                match = true
                break
            }
        }
        
        if match {
            break
        }
    }

    if match {
        return nil
    }

    return &types.CasError{Error: fmt.Errorf("Proxy service `%s` is not allowed to get proxy ticket with service `%s`", pgtURL, serviceID), CasErrorCode: types.CAS_ERROR_CODE_UNAUTHORIZED_SERVICE_PROXY}
}

// ValidateProxyURLEndpoint reaches out to the proxy URL
func ValidateProxyURLEndpoint(pgtURL string) *types.CasError {
    _, err := http.Get(pgtURL)
    if err != nil {
        return &types.CasError{Error: fmt.Errorf("Proxy service `%s` validation failed with error: `%s`", pgtURL, err.Error()), CasErrorCode: types.CAS_ERROR_CODE_INVALID_PROXY_CALLBACK}
    }
    
    return nil
}