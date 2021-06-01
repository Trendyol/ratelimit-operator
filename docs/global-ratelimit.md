# Global Rate Limit

Rate limits are configured in two places: first is the envoy “rate_limits filters” and the second is envoy’s “ratelimit service” configuration. Envoy’s filters contain “actions”, which result in a “descriptor”. This “descriptor” is sent to the ratelimit service, which uses them to make a decision on a specific limit.

## 1. Example Configuration

```yaml
apiVersion: trendyol.com/v1beta1
kind: GlobalRateLimit
metadata:
  name: favorite-api
  namespace: browsing-team
spec:
  domain: discovery-reco-favorite-service
  rate:
  - dimensions:
    - request_headers:
        descriptor_key: agentName
        header_name: x-agentname
    request_per_unit: 1000
    unit: minute
  - dimensions:
    - header_value_match:
        descriptor_value: healthcheck
        headers:
        - name: :path
          prefix_match: /_monitoring/health
    request_per_unit: 10
    unit: minute
  workload: favorite-api
```
**Domain** : Must be same with service discovery name with all DC
**Workload** : Must be select kubernetes labels app name

## 1.1 Request_headers defination example
Let’s assume we have a single **“rate_limits rate”** with two **“dimensions”**.

The first **“dimensions”** block reads as follows: match any request with header name x-agentname key. For example if you want to rate limit any of header key you need to declare that. This block reads for each x-agentname header key so you can make 1000 request in minute.

```yaml
  - dimensions:
    - request_headers:
        descriptor_key: clientip
        header_name: x-client-ip
    request_per_unit: 1000
    unit: minute
```
If you want to rate limit for another header name you need to declare different “dimensions”.
## 1.2 Request_headers defination specific value example

If you wan to rate limit header specific value. You need to add value property.

```yaml
  - dimensions:
    - request_headers:
        descriptor_key: agentname
        header_name: x-agentname
        value: zeus-social
    request_per_unit: 1000
    unit: minute
```
This configuration apply to rate limit service x-agentname header key and when you have header value zeus-social you can make  1000 request in minute.
