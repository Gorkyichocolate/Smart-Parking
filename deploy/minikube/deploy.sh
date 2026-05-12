#!/bin/bash
# deploy/minikube/deploy.sh

set -e

echo "🚀 Starting Minikube..."
minikube start --cpus=4 --memory=8192 --driver=docker

echo "📦 Enabling addons..."
minikube addons enable metrics-server
minikube addons enable ingress
minikube addons enable dashboard

echo "🐳 Building Docker images in Minikube..."
eval $(minikube docker-env)

cd ../../services/payment-service
docker build -t payment-service:latest .

cd ../notification-service
docker build -t notification-service:latest .

cd ../../deploy/k8s

echo "📋 Applying Kubernetes manifests..."
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f deployments.yaml
kubectl apply -f services.yaml

echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=payment-service -n smart-parking --timeout=300s || true
kubectl wait --for=condition=ready pod -l app=notification-service -n smart-parking --timeout=300s || true

echo "🔌 Setting up port forwarding..."
kubectl port-forward -n smart-parking service/payment-service 50052:50052 &
kubectl port-forward -n smart-parking service/payment-service 9091:9090 &
kubectl port-forward -n smart-parking service/notification-service 9092:9090 &

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📍 Access URLs:"
echo "   Payment Service gRPC: localhost:50052"
echo "   Payment Metrics: http://localhost:9091/metrics"
echo "   Notification Metrics: http://localhost:9092/metrics"
echo "   Grafana: http://localhost:3000 (admin/admin)"
echo "   Prometheus: http://localhost:9090"
echo "   RabbitMQ: http://localhost:15672 (guest/guest)"
echo ""
echo "📊 To access Grafana dashboard:"
echo "   kubectl port-forward -n smart-parking service/grafana 3000:3000"