#!/bin/bash
# Quick Docker deployment script for ShadowNet

set -e

echo "🐳 ShadowNet Docker Deployment"
echo "=============================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    echo "Please install Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed"
    echo "Please install Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"
echo ""

# Parse command
case "$1" in
    test)
        echo "🧪 Starting test environment (3 nodes + control plane + dashboard)"
        echo ""
        docker-compose up -d --build
        echo ""
        echo "✅ Services started!"
        echo ""
        echo "📊 Dashboard: http://localhost:3000"
        echo "🔌 Control Plane API: http://localhost:8080"
        echo ""
        echo "View logs: docker-compose logs -f"
        echo "Check status: docker-compose ps"
        ;;
    
    prod)
        echo "🚀 Starting production environment (control plane + dashboard only)"
        echo ""
        docker-compose -f docker-compose.prod.yml up -d --build
        echo ""
        echo "✅ Services started!"
        echo ""
        echo "📊 Dashboard: http://localhost:3000"
        echo "🔌 Control Plane API: http://localhost:8080"
        echo ""
        echo "⚠️  Remember to update NEXT_PUBLIC_CONTROLPLANE_URL in docker-compose.prod.yml"
        echo "    with your server's public IP or domain"
        ;;
    
    stop)
        echo "🛑 Stopping all services..."
        docker-compose down
        echo "✅ Services stopped"
        ;;
    
    logs)
        echo "📋 Showing logs (Ctrl+C to exit)..."
        docker-compose logs -f
        ;;
    
    status)
        echo "📊 Service Status:"
        echo ""
        docker-compose ps
        ;;
    
    clean)
        echo "🧹 Cleaning up (removing containers, networks, and volumes)..."
        docker-compose down -v
        echo "✅ Cleanup complete"
        ;;
    
    restart)
        echo "🔄 Restarting services..."
        docker-compose restart
        echo "✅ Services restarted"
        ;;
    
    *)
        echo "Usage: $0 {test|prod|stop|logs|status|clean|restart}"
        echo ""
        echo "Commands:"
        echo "  test     - Start test environment with 3 nodes (single machine)"
        echo "  prod     - Start production environment (control plane + dashboard)"
        echo "  stop     - Stop all services"
        echo "  logs     - View logs from all services"
        echo "  status   - Check status of all services"
        echo "  clean    - Stop and remove all containers, networks, and volumes"
        echo "  restart  - Restart all services"
        echo ""
        echo "Examples:"
        echo "  $0 test      # Test with 3 nodes on your machine"
        echo "  $0 logs      # Watch logs"
        echo "  $0 stop      # Stop everything"
        exit 1
        ;;
esac
