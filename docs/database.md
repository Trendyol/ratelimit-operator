## Redis

We use active-active geo-distributed redis topology. 

Redis key = domain + dimension value | service name + timestamp given a period
Redis value = hit counter

Example:

Key= discovery-reco-favorite-service_header_match_healthcheck_1621587240

Value = 3