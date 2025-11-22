# Notes Backend API

A Go backend API server built with gorilla/mux, featuring CORS support and hot-reload development.

## Setup

### Prerequisites

- Go 1.25.3 or higher
- Git

### Installation

1. Clone the repository and navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

## Dependencies

### Core Dependencies

- **gorilla/mux** (`v1.8.1`) - HTTP router and URL matcher
- **gorilla/handlers** (`v1.5.2`) - HTTP handlers (optional, for additional middleware)

### Development Dependencies

- **CompileDaemon** (`v1.4.0`) - Auto-reload server on file changes
- **fsnotify** (`v1.4.9`) - File system notifications
- **watcher** (`v1.0.7`) - File watching utilities

## Project Structure

```
backend/
├── controllers/          # Request handlers
│   ├── getNotes.go
│   └── hello.go
├── middleware/          # Middleware functions
│   └── middleware.go     # CORS middleware
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
└── README.md            # This file
```

## Setting Up Gorilla Mux

The router is initialized in `main.go`:

```go
import "github.com/gorilla/mux"

func main() {
    // Initialize router
    r := mux.NewRouter()
    
    // Define routes
    r.HandleFunc("/", controllers.Hello)
    r.HandleFunc("/notes", controllers.GetNotes).Methods("POST", "OPTIONS")
    
    // Start server
    http.ListenAndServe(":8080", r)
}
```

### Route Methods

You can specify HTTP methods for routes:
```go
r.HandleFunc("/notes", controllers.GetNotes).Methods("POST", "GET", "PUT", "DELETE")
```

### Route Parameters

Mux supports path variables:
```go
r.HandleFunc("/notes/{id}", controllers.GetNote).Methods("GET")
// Access with: mux.Vars(r)["id"]
```

## CORS Configuration

CORS (Cross-Origin Resource Sharing) is handled by custom middleware in `middleware/middleware.go`.

### Current CORS Setup

The middleware:
- Allows all origins (`*`) - **Change this in production!**
- Allows methods: GET, POST, PUT, DELETE, OPTIONS
- Allows headers: Content-Type, Authorization
- Handles preflight OPTIONS requests automatically

### Implementation

```go
func EnableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Applying CORS Middleware

In `main.go`, apply the middleware to all routes:

```go
import "github.com/aminofabian/notes/middleware"

func main() {
    r := mux.NewRouter()
    
    // Apply CORS middleware to all routes
    r.Use(middleware.EnableCORS)
    
    // ... routes
}
```

### Production CORS Configuration

For production, restrict allowed origins:

```go
// Allow specific origins
allowedOrigins := []string{
    "https://yourdomain.com",
    "https://www.yourdomain.com",
}

// In middleware, check origin:
origin := r.Header.Get("Origin")
for _, allowed := range allowedOrigins {
    if origin == allowed {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        break
    }
}
```

## Running the Server

### Standard Run

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Using CompileDaemon (Hot Reload)

CompileDaemon automatically rebuilds and restarts the server when files change.

#### Installation

```bash
go install github.com/githubnemo/CompileDaemon@latest
```

#### Usage

```bash
CompileDaemon -command="./notes"
```

Or with more options:

```bash
CompileDaemon \
  -directory=. \
  -command="go run main.go" \
  -pattern="(.+\\.go)$" \
  -exclude-dir=".git"
```

#### CompileDaemon Options

- `-directory`: Directory to watch (default: current directory)
- `-command`: Command to run (e.g., `go run main.go` or `./notes`)
- `-pattern`: Regex pattern for files to watch (default: `\.go$`)
- `-exclude-dir`: Directories to exclude (e.g., `.git`, `node_modules`)

#### Build and Run with CompileDaemon

If you want to build first, then run:

```bash
CompileDaemon \
  -build="go build -o notes main.go" \
  -command="./notes"
```

## Development Workflow

1. **Start with hot reload:**
   ```bash
   CompileDaemon -command="go run main.go"
   ```

2. **Make changes** to your `.go` files

3. **Server automatically restarts** when files are saved

4. **Test your API:**
   ```bash
   curl http://localhost:8080/
   curl -X POST http://localhost:8080/notes
   ```

## API Endpoints

- `GET /` - Hello endpoint
- `POST /notes` - Get/create notes (with CORS support)

## Troubleshooting

### CORS Errors

If you see CORS errors in the browser:
1. Ensure `r.Use(middleware.EnableCORS)` is applied in `main.go`
2. Check that OPTIONS method is included: `.Methods("POST", "OPTIONS")`
3. Verify CORS headers are being set (check browser Network tab)

### CompileDaemon Not Restarting

- Check that CompileDaemon is watching the correct directory
- Verify file patterns match your `.go` files
- Ensure you're saving files (not just closing them)

### Port Already in Use

If port 8080 is already in use:
```go
// Change in main.go
http.ListenAndServe(":3000", r)  // or any other port
```

## Additional Resources

- [Gorilla Mux Documentation](https://github.com/gorilla/mux)
- [CompileDaemon GitHub](https://github.com/githubnemo/CompileDaemon)
- [CORS MDN Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
