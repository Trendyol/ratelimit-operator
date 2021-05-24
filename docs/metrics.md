### With StatsD

Rate limit exporter translates StatsD metrics to Prometheus metrics via
configured mapping rules.

    +----------+                         +-------------------+                        +--------------+
    |  StatsD  |---(UDP/TCP repeater)--->|  statsd_exporter  |<---(scrape /metrics)---|  Prometheus  |
    +----------+                         +-------------------+                        +--------------+

## Total Rate limit request

This panel shows a time series with the total hits per rate limit configured. In this panel service owners can see trends over time.

![ratelimit](images/total_hits.png)

## Over Limit request

This panel shows the metrics that are over the configured limit. This panel allows service owners to have quantifiable data with which to go back to their services and assess call patterns, and do capacity planning for high load events.
![ratelimit](images/over_limits.png)


## Runtime Success Configuration 

This panel shows when the metric is hitting 80% of the limit configured.

![config](images/success-config.png)
