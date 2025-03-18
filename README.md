dokcer run --net=host -v ./deploy/prometheus:etc/prometheus prom/prometheus
docker run -d -p 3000:3000 --name=grafana grafana/grafana-oss
