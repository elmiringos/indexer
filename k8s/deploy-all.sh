#!/bin/bash

set -e

echo "🔧 Creating PersistentVolumeClaim..."
kubectl apply -f k8s/storageclass/manual.yaml
kubectl apply -f k8s/pvc/postgres-pv.yaml
kubectl apply -f k8s/pvc/postgres-pvc.yaml

echo "🐘 Deploy PostgreSQL..."
kubectl apply -f k8s/deployments/postgres.yaml

echo "📦 Deploy RabbitMQ и Redis..."
kubectl apply -f k8s/deployments/rabbitmq.yaml
kubectl apply -f k8s/deployments/redis.yaml

echo "🌐 Creating services for PostgreSQL, RabbitMQ, Redis..."
kubectl apply -f k8s/services/postgres.yaml
kubectl apply -f k8s/services/rabbitmq.yaml
kubectl apply -f k8s/services/redis.yaml

echo "⏳ Waiting for db, cache and broker ..."
kubectl wait --for=condition=available --timeout=300s deployment/postgres
kubectl wait --for=condition=available --timeout=300s deployment/rabbitmq
kubectl wait --for=condition=available --timeout=300s deployment/redis

echo "🚀 Deploy core, producer and explorer..."
kubectl apply -f k8s/deployments/core.yaml
kubectl apply -f k8s/deployments/producer.yaml
kubectl apply -f k8s/deployments/explorer.yaml

echo "🌐 Creating services for core, producer, and explorer..."
kubectl apply -f k8s/services/core.yaml
kubectl apply -f k8s/services/producer.yaml
kubectl apply -f k8s/services/explorer.yaml

echo "✅ Deploy is finished!"
