## Terminology

### Domain

A domain is a container for a set of rate limits. All domains known to the Ratelimit service must be globally unique. They serve as a way for different teams/projects to have rate limit configurations that don't conflict.

### Descriptor

A descriptor is a list of key/value pairs owned by a domain that the Ratelimit service uses to select the correct rate limit to use when limiting. Descriptors are case-sensitive. Examples of descriptors are:

- ("client_ip", "10.2.80.89")
- ("header_key", "test")