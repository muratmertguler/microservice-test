docker run -p 9090:9090 -v $PWD/.deploy/prometheus:/etc/prometheus  prom/prometheus

docker run --net=host -v $PWD/.deploy/prometheus:/etc/prometheus  prom/prometheus

docker run -d -p 3000:3000 --name=grafana grafana/grafana-oss

docker run -d --net=host --name=grafana grafana/grafana-oss

Prometheus -> https://github.com/grafana/grafana/issues/46434
