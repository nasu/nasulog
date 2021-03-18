#!/bin/sh
eval $(minikube -p minikube docker-env)
docker build -t nasu/nasulog:v0.1 ../src/
kubectl apply -f deployment.yaml

### TODO: これ良くないよ
#minikube service -n nasulog dynamodb &
#minikube service -n nasulog elasticmq &
#minikube service -n nasulog api &
