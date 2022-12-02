# Menace-Go

A basic Golang, web based implementation of menace

## Run

```bash
go run .
```

## Run

```bash
go run build
```

Tutorial by captain obvious

docker build --tag dockerized-menace .
docker tag dockerized-menace europe-west2-docker.pkg.dev/key-autumn-362512/quickstart-docker-repo/dockerized-menace
docker push europe-west2-docker.pkg.dev/key-autumn-362512/quickstart-docker-repo/dockerized-menace

kubectl delete -f ./deployment.yml
kubectl apply -f ./deployment.yml

NODEIP
NODEPORT

NODEIP : kubectl get nodes --output wide
NODEPORT : kubectl get service np-service --output yaml
gcloud compute firewall-rules create test-32147-port --allow tcp:32147
35.246.118.144:32147
