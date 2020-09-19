
# custom-plugin

This sample implements Envoy's (external authorization)[https://www.envoyproxy.io/docs/envoy/latest/api-v2/config/filter/http/ext_authz/v2/ext_authz.proto] filter

## Testing locally

```bash
./build-docker.sh
```

```bash
curl -v http://localhost:8080/httpbin
```

```bash
./clean-docker.sh
```

___

## Support

This is not an officially supported Google product
