# Rate limits


## Overview

Rate limiting is a means of protecting backend services against unwanted traffic. This can be useful for a variety of different scenarios:

- Protecting against denial-of-service (DoS) attacks by malicious actors
- Protecting against DoS incidents due to bugs in client applications/services
- Login, Sms, Mail Service control

Envoy supports two forms of HTTP rate limiting: local and global.

In local rate limiting, rate limits are enforced by each Envoy instance, without any communication with other Envoys or any external service.

In global rate limiting, an external rate limit service (RLS) is queried by each Envoy via gRPC for rate limit decisions.


## Architecture

![ratelimit](docs/images/ratelimit.jpeg)


## Local Rate Limit
See [Local Rate Limit](docs/local-ratelimit.md)


## Global Rate Limit
See [Global Rate Limit](docs/global-ratelimit.md)


## Metrics && Stats
See [Stats](docs/global-ratelimit.md)
