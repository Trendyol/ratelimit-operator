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

## 1.1 request_headers definition example
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
## 1.2 request_headers definition specific value example

you can rate limit header specific value.

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


## 2.1 header_value_match definition specific **prefix_match** example

This configuration apply to rate limit service method and url path match actions. You can configure by route based policy.
If specified, header match will be performed based on the prefix of the header value. Note: empty prefix is not allowed, please use present_match instead.

Examples:

The prefix abcd matches the value abcdxyz, but not for abcxyz.

```yaml
    - dimensions:
        - header_value_match:
            descriptor_value: similar-category
            headers:
              - name: ':path'
                prefix_match: /product-recommendation/similar-category
      request_per_unit: 1
      unit: minute
```

## 2.2 header_value_match definition specific **suffix_match** example

If specified, header match will be performed based on the suffix of the header value. Note: empty suffix is not allowed, please use present_match instead.

Examples:

The suffix abcd matches the value xyzabcd, but not for xyzbcd.
```yaml
    - dimensions:
        - header_value_match:
            descriptor_value: similar-category
            headers:
              - name: ':path'
                suffix_match: /product-recommendation/similar-category
      request_per_unit: 1
      unit: minute
```

## 2.3 header_value_match definition specific **contains_match** example

If specified, header match will be performed based on whether the header value contains the given value or not. Note: empty contains match is not allowed, please use present_match instead.

Examples:

The value abcd matches the value xyzabcdpqr, but not for xyzbcdpqr.


```yaml
    - dimensions:
        - header_value_match:
            descriptor_value: similar-category
            headers:
              - name: ':path'
                contains_match: /product-recommendation/similar-category
      request_per_unit: 1
      unit: minute
```

## 2.4 header_value_match definition specific **safe_regex_match** example

If specified, this regex string is a regular expression rule which implies the entire request header value must match the regex. The rule will not match if only a subsequence of the request header value matches the regex.

Specifies how the header match will be performed to route the request.


```yaml
    - dimensions:
        - header_value_match:
            descriptor_value: similar-category
            headers:
              - name: ':path'
                safe_regex_match:
                  google_re2: {}
                  regex: "/product-recommendation/similar-category/\d+/id"
      request_per_unit: 1
      unit: minute
```

## 2.5 header_value_match definition specific **GET** request example

This configuration apply to rate limit service only GET requests.

```yaml
    - dimensions:
        - header_value_match:
            descriptor_value: similar-category
            headers:
              - name: ':method'
                exact_match: "GET"
      request_per_unit: 1
      unit: minute
```

## 2.6 header_value_match definition specific **host header** request example

This configuration apply to rate limit service only specific host header requests.

```yaml
    - dimensions:
        - header_value_match:
            descriptor_value: similar-category
            headers:
              - name: ':authority'
                exact_match: "discovery-sellerstore-follow-service.earth.trendyol.com"
      request_per_unit: 1
      unit: minute
```

## 2.7 header_value_match and request_headers definition

This configuration apply to rate limit service header_value_match and request_header.

Ratelimit service matches on a full “descriptor”, not on individual “descriptor entries”. In order to match a “descriptor” with multiple “descriptor entries”, a nested “descriptor configuration” must be used. In this case, nested “descriptor configurations” are joined by a logical AND.

```yaml
    - dimensions:
        - header_value_match:
            descriptor_value: get_limit
            headers:
              - name: ':path'
                prefix_match: /api
              - name: ':method'
                prefix_match: GET
      request_per_unit: 1
      unit: minute
```

You can limit starting with /api endpoint and GET requests.


### Property Definition

| Command | Description |
| --- | --- |
| `header_value_match` | Rate limit on the existence of request headers. |
| `request_headers` | Rate limit on request headers.|
| `header_name` | (string, REQUIRED) The header name to be queried from the request headers. The header’s value is used to populate the value of the descriptor entry for the descriptor_key.|
| `descriptor_key` | (string, REQUIRED) The key to use in the descriptor entry.|
| `descriptor_value` | (string, REQUIRED) The value to use in the descriptor entry.|
| `headers` |  Specifies a set of headers that the rate limit action should match on. The action will check the request’s headers against all the specified headers in the config. A match will happen if all the headers in the config are present in the request with the same values (or based on presence if the value field is not in the config).|
| `name` | (string, REQUIRED) Specifies the name of the header in the request.|
| `exact_match` | Enable or disable policy.|
| `safe_regex_match` | (type.matcher.v3.RegexMatcher) If specified, this regex string is a regular expression rule which implies the entire request header value must match the regex. The rule will not match if only a subsequence of the request header value matches the regex.|
| `present_match` | (bool) If specified as true, header match will be performed based on whether the header is in the request. If specified as false, header match will be performed based on whether the header is absent.|
| `prefix_match` | (string) If specified, header match will be performed based on the prefix of the header value. Note: empty prefix is not allowed, please use present_match instead.|
| `suffix_match` | (string) If specified, header match will be performed based on the suffix of the header value. Note: empty suffix is not allowed, please use present_match instead..|
| `contains_match` | (string) If specified, header match will be performed based on whether the header value contains the given value or not. Note: empty contains match is not allowed, please use present_match instead.|
| `disabled` | Enable or disable policy.|


For more details [https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#envoy-v3-api-msg-config-route-v3-ratelimit-action-requestheaders]
