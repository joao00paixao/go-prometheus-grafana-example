# go-prometheus-grafana-example
 Repository for Go API + Worker storing RPS in Grafana

# Installation

```
docker run -d -p 3000:3000 grafana/grafana 

docker pull prom/prometheus
docker run -d \
    -p 9090:9090 \
    -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus


```

Use the internal host not localhost otherwise there will be connection issues.