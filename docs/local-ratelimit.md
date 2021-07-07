# Local Rate Limit

The HTTPProxy API supports defining local rate limit policies that can be applied to either individual entire virtual hosts. Local rate limit policies define a maximum number of requests per unit of time that an Envoy should proxy to the upstream service.
Requests beyond the defined limit will receive a 429 (Too Many Requests) response by default. Local rate limit policies program Envoy’s HTTP local rate limit filter.

It’s important to note that local rate limit policies apply per Envoy pod. For example, a local rate limit policy of 100 requests per second for a given route will result in each Envoy pod allowing up to 100 requests per second for that route.

## Configuration

A local rate limit policy requires a ``max_token, fill_interval and tokens_per_fill`` fields, defining the number of max token per unit of time that are allowed. Requests must be a positive integer. We need to specify which deployment can be applied this configuration. We use ``workload`` key that selects kubernetes deployment label app.


```yaml
kind: LocalRateLimit
apiVersion: trendyol.com/v1beta1
metadata:
  name: favorite-api
  namespace: browsing-team
spec:
  token_bucket:
    fill_interval: 10s
    max_tokens: 100
    tokens_per_fill: 100
  workload: favorite-api  //Kubernetes labels app name
  disabled: true | false //Optinal if you want to disable policy 
```

### Property Definition

| Command | Description |
| --- | --- |
| `max_tokens` | The maximum request for given a period. |
| `tokens_per_fill` | The number of requests added to the bucket during each fill interval.|
| `fill_interval` | The fill interval that requests are added to the bucket.|
| `disabled` | Enable or disable policy.|



## Testing

- ``Response header should return`` x-local-rate-limit: true 
- ``Status code should be`` 429 (Too Many Requests)
