# Pattern 2 - Proxy chaining/ Envoy "sandwich"

This sample implements Envoy's [external authorization](https://www.envoyproxy.io/docs/envoy/latest/api-v2/config/filter/http/ext_authz/v2/ext_authz.proto) filter to demonstrate simple routing of requests to upstream targets

## Scenario

In this example, a client sends requests to an endpoints serviced by Envoy. The consumer passes a header (`x-backend-name`) with the name the name of the backend. Envoy calls the ext_authz service and dynamically routes the consumer to the appropriate backend.

![Routing Sample](../envoy-pattern2.png)
