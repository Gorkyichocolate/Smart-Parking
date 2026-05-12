#!/bin/bash
# docker/start.sh

set -e

echo "🚀 Starting Smart Parking services with Docker Compose..."

cd "$(dirname "$0")"

# Create .env file if not exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cat > .env << 'EOF'
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=
EOF
fi

# Start services
docker-compose up -d

echo ""
echo "✅ Services started successfully!"
echo ""
echo "📍 Access URLs:"
echo "   Grafana: http://localhost:3000 (admin/admin)"
echo "   Prometheus: http://localhost:9090"
echo "   RabbitMQ: http://localhost:15672 (guest/guest)"
echo "   Payment Service Metrics: http://localhost:9091/metrics"
echo "   Notification Metrics: http://localhost:9092/metrics"
echo ""
echo "📊 To view logs: docker-compose logs -f"
echo "🛑 To stop: docker-compose down"