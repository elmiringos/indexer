#!/bin/bash

set -e

echo "🧹 Deleting services..."
kubectl delete -f k8s/services/explorer.yaml --ignore-not-found
kubectl delete -f k8s/services/producer.yaml --ignore-not-found
kubectl delete -f k8s/services/core.yaml --ignore-not-found

echo "🧨 Deleteing deployments..."
kubectl delete -f k8s/deployments/explorer.yaml --ignore-not-found
kubectl delete -f k8s/deployments/producer.yaml --ignore-not-found
kubectl delete -f k8s/deployments/core.yaml --ignore-not-found
kubectl delete -f k8s/deployments/rabbitmq.yaml --ignore-not-found
kubectl delete -f k8s/deployments/redis.yaml --ignore-not-found
kubectl delete -f k8s/deployments/postgres.yaml --ignore-not-found

echo "📦 Deleting PVC..."
kubectl delete -f k8s/pvc/postgres-pvc.yaml --ignore-not-found
kubectl delete -f k8s/pvc/postgres-pv.yaml --ignore-not-found
kubectl delete -f k8s/storageclass/manual.yaml --ignore-not-found

echo "✅ All resources were deleted!"
