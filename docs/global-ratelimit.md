# Global Rate Limit

Rate limits are configured in two places: first is the envoy “rate_limits filters” and the second is envoy’s “ratelimit service” configuration. Envoy’s filters contain “actions”, which result in a “descriptor”. This “descriptor” is sent to the ratelimit service, which uses them to make a decision on a specific limit.



## Configuration

