# HomeCooking 🍳

A modern, open-source recipe storage web app built for families and cooking enthusiasts.

## Features

- 📚 **Recipe Management**: Store, organize, and browse recipes with markdown support
- 🏷️ **Categorization & Tagging**: Organize recipes by categories and custom tags
- 🍽️ **Recipe Groups**: Group recipes that are meant to be served together (e.g., "Biscuits & Gravy")
- 🖼️ **Image Support**: Featured images and inline body images with automatic optimization
- 🤖 **AI Integration**: Optional AI features for recipe extraction from images and text enhancement
- 🔗 **Sharing**: Share recipes via links, user invites, and cross-instance sharing
- 👥 **Multi-user**: User accounts with role-based permissions (admin, editor, user)
- 📱 **Mobile-First**: Responsive design that works perfectly on all devices
- 🎨 **Beautiful UI**: Clean, modern interface built with Tailwind CSS v4

## Tech Stack

### Backend
- Go 1.21+ with `net/http` standard library router
- `sqlc` for type-safe database queries
- Database: PostgreSQL, SQLite, or MySQL
- JWT authentication with bcrypt
- Local filesystem storage with image optimization

### Frontend
- Astro 5.x with Tailwind CSS v4
- Route-based code splitting
- Mobile-first responsive design
- Native Markdown support

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL, SQLite, or MySQL

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/homecooking.git
   cd homecooking
   ```

2. **Backend Setup**
   ```bash
   cd backend
   go mod download
   go build -o bin/server ./cmd/server
   ```

3. **Frontend Setup**
   ```bash
   cd ../frontend
   npm install
   ```

4. **Configuration**
   
   Create a `.env` file in the `backend` directory:
   ```bash
   # Server
   SERVER_PORT=8080
   SERVER_ENV=development
   SERVER_BASE_URL=http://localhost:8080

   # Database (PostgreSQL example)
   DATABASE_TYPE=postgres
   DATABASE_HOST=localhost
   DATABASE_PORT=5432
   DATABASE_NAME=homecooking
   DATABASE_USER=postgres
   DATABASE_PASSWORD=yourpassword

   # Auth
   JWT_SECRET=your-jwt-secret-here
   REFRESH_SECRET=your-refresh-secret-here
   TOKEN_EXPIRY_HOURS=24

   # AI (Optional - disabled by default)
   AI_ENABLED=false
   AI_PROVIDER=openai
   AI_API_KEY=
   AI_MODEL=gpt-4o

   # Storage
   STORAGE_TYPE=local
   STORAGE_LOCAL_PATH=./uploads
   ```

5. **Database Setup**
   
   For PostgreSQL:
   ```bash
   createdb homecooking
   psql -d homecooking -f backend/internal/db/migrations/001_init.up.sql
   ```

   For SQLite (simpler option):
   ```bash
   DATABASE_TYPE=sqlite
   DATABASE_PATH=./data.db
   # Tables will be created automatically on first run
   ```

6. **Run the Server**
   ```bash
   cd backend
   ./bin/server
   ```

7. **Run the Frontend** (in a new terminal)
   ```bash
   cd frontend
   npm run dev
   ```

8. **Access the Application**
   - Frontend: http://localhost:4321
   - Backend API: http://localhost:8080

## Development

### Backend Development

```bash
cd backend
go run ./cmd/server
```

### Frontend Development

```bash
cd frontend
npm run dev
```

### Running sqlc

After modifying SQL queries, regenerate the Go code:

```bash
cd backend
export PATH="$HOME/go/bin:$PATH"
sqlc generate
```

## Project Structure

```
homecooking/
├── backend/
│   ├── cmd/server/           # Main application entry
│   ├── internal/
│   │   ├── config/          # Configuration management
│   │   ├── db/              # Database layer (sqlc generated)
│   │   ├── models/          # Domain models
│   │   ├── repository/      # Data access layer
│   │   ├── services/        # Business logic
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # HTTP middleware
│   │   └── ai/             # AI integration
│   ├── storage/             # File storage
│   └── static/uploads/     # Uploaded images
├── frontend/
│   ├── src/
│   │   ├── components/      # Astro components
│   │   ├── layouts/         # Page layouts
│   │   ├── pages/          # Route pages
│   │   └── lib/            # Utilities
│   └── public/uploads/     # Frontend static assets
└── docs/                  # Documentation
```

## API Documentation

See [API.md](docs/API.md) for complete API documentation.

## Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md) for contribution guidelines.

## License

MIT License - see LICENSE file for details

## Acknowledgments

Built with ❤️ for families who love to cook together.
