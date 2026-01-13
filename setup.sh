#!/bin/bash

# HomeCooking Quick Start Script

set -e

echo "🍳 HomeCooking Quick Start"
echo "=========================="
echo ""

# Check for required tools
echo "Checking dependencies..."

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or higher."
    exit 1
fi

echo "✅ All dependencies found"
echo ""

# Backend setup
echo "📦 Setting up backend..."
cd backend

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    echo "⚠️  Please edit backend/.env with your database settings"
fi

# Build backend
echo "Building backend..."
go build -o bin/server ./cmd/server
echo "✅ Backend built"
echo ""

# Frontend setup
echo "📦 Setting up frontend..."
cd ../frontend

# Install Node dependencies
echo "Installing Node dependencies..."
npm install

# Build frontend
echo "Building frontend..."
npm run build
echo "✅ Frontend built"
echo ""

echo "✨ Setup complete!"
echo ""
echo "To start the application:"
echo ""
echo "1. Configure your database in backend/.env"
echo "   For SQLite (easiest):"
echo "   DATABASE_TYPE=sqlite"
echo "   DATABASE_PATH=./data.db"
echo ""
echo "2. Start the backend:"
echo "   cd backend && ./bin/server"
echo ""
echo "3. Start the frontend (in a new terminal):"
echo "   cd frontend && npm run dev"
echo ""
echo "4. Open your browser:"
echo "   Frontend: http://localhost:4321"
echo "   Backend API: http://localhost:8080"
echo ""
echo "📚 For more information, see README.md"
