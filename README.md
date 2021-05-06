# ratelimit-operator


```yaml
kind: LocalRateLimit
apiVersion: trendyol.com/v1beta1
metadata:
  name: localratelimit-sample2
  namespace: default
spec:
  workload: application
  token_bucket:
    fill_interval: 1s
    max_tokens: 3000
    tokens_per_fill: 1000


```