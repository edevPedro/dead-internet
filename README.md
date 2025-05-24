yeah, even this project was made by an ai.

# 🤖 Dead Internet Project

A Go-based social network that demonstrates the "Dead Internet Theory" - where AI automatically responds to every user comment, creating an eerie simulation of artificial online interactions.

## 🎯 Project Overview

This project explores the boundaries between human and artificial interactions in social media. Every time a user posts a comment, an AI bot automatically generates and posts a response, creating the unsettling feeling that you might be the only real person in a sea of artificial conversations.

## ✨ Features

- **Real-time AI Responses**: Every user comment triggers an automatic AI-generated response
- **Full Social Network API**: Complete CRUD operations for users, posts, and comments
- **Modern Web Interface**: Clean, responsive frontend with real-time updates
- **PostgreSQL Database**: Robust data persistence with proper relationships
- **CORS-enabled API**: Ready for frontend integration
- **Visual AI Indicators**: Clear badges distinguish AI-generated content

## 🏗️ Architecture

### Backend (Go)
- **Chi Router**: Fast HTTP router with middleware support
- **PostgreSQL**: Database with proper foreign key relationships
- **Store Pattern**: Clean separation of data access logic
- **Context-aware**: Proper request context handling throughout

### Database Schema
```sql
-- Users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Posts table  
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Comments table with AI support
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    is_ai_generated BOOLEAN DEFAULT FALSE,
    parent_comment_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### API Endpoints

#### Posts
- `POST /v1/posts` - Create a new post
- `GET /v1/posts` - Get all posts
- `GET /v1/posts/{id}` - Get specific post
- `PUT /v1/posts/{id}` - Update post
- `DELETE /v1/posts/{id}` - Delete post
- `GET /v1/posts/{id}/comments` - Get post comments

#### Users
- `POST /v1/users` - Create a new user
- `GET /v1/users/{id}` - Get user details
- `GET /v1/users/{id}/posts` - Get user's posts

#### Comments
- `POST /v1/comments` - Create comment (triggers AI response)

#### System
- `GET /v1/health` - Health check

## 🚀 Getting Started

### Prerequisites
- Go 1.23+
- PostgreSQL 15+
- Modern web browser

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/edevPedro/dead-internet.git
cd dead-internet
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Set up PostgreSQL**
```bash
# Install PostgreSQL (Ubuntu/Debian)
sudo apt update && sudo apt install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres psql
CREATE DATABASE deadinternet;
CREATE USER deadinternet WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE deadinternet TO deadinternet;
\q
```

4. **Deploy database schema**
```bash
psql -h localhost -U deadinternet -d deadinternet -f scripts/schema.sql
```

5. **Build and run the application**
```bash
# Build
go build -o bin/main cmd/api/*.go

# Set environment variables
export ADDR=":12000"
export DB_ADDR="postgres://deadinternet:password@localhost/deadinternet?sslmode=disable"

# Run the API server
./bin/main
```

6. **Start the web interface**
```bash
# In a new terminal
cd web
python3 -m http.server 12001 --bind 0.0.0.0
```

7. **Access the application**
- Web Interface: http://localhost:12001
- API: http://localhost:12000

## 🎮 Usage

### Creating Users
```bash
curl -X POST http://localhost:12000/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@example.com"}'
```

### Creating Posts
```bash
curl -X POST http://localhost:12000/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello World", "content": "My first post!", "user_id": 1}'
```

### Adding Comments (Triggers AI Response)
```bash
curl -X POST http://localhost:12000/v1/comments \
  -H "Content-Type: application/json" \
  -d '{"content": "This is strange...", "user_id": 1, "post_id": 1}'
```

## 🤖 AI Response System

The AI response system automatically generates contextual replies to user comments using a curated set of responses:

- **Immediate Response**: AI comments are generated asynchronously within seconds
- **Contextual Variety**: 8 different response templates for varied interactions
- **Clear Identification**: AI comments are marked with `is_ai_generated: true`
- **No User Association**: AI comments have `user_id: null` to distinguish them

### Sample AI Responses
- "That's an interesting perspective! I think there's more to consider here."
- "I completely agree with your point. This is exactly what I was thinking."
- "Have you considered the alternative viewpoint? It might be worth exploring."
- "This reminds me of a similar discussion I had recently. Great insights!"

## 🎨 Frontend Features

- **Real-time Updates**: Comments refresh automatically to show new AI responses
- **Visual Distinction**: AI comments have red borders and "AI" badges
- **Responsive Design**: Works on desktop and mobile devices
- **Form Validation**: Proper input validation and error handling
- **Success Feedback**: Clear notifications for user actions

## 🔧 Development

### Project Structure
```
dead-internet/
├── cmd/api/           # API handlers and main application
│   ├── main.go        # Application entry point
│   ├── api.go         # Router and server setup
│   ├── posts.go       # Post handlers
│   ├── users.go       # User handlers
│   ├── comments.go    # Comment handlers (with AI logic)
│   ├── health.go      # Health check
│   └── middleware.go  # CORS and other middleware
├── internal/store/    # Data access layer
│   ├── storage.go     # Storage interface
│   ├── posts.go       # Post store implementation
│   ├── users.go       # User store implementation
│   └── comments.go    # Comment store implementation
├── scripts/           # Database scripts
│   └── schema.sql     # Database schema
├── web/              # Frontend
│   └── index.html    # Single-page application
├── go.mod            # Go dependencies
└── README.md         # This file
```

### Adding New AI Responses
Edit `cmd/api/comments.go` and add new responses to the `aiResponses` slice:

```go
aiResponses := []string{
    "Your new AI response here...",
    // ... existing responses
}
```

## 🌐 Deployment

### Environment Variables
- `ADDR`: Server address (default: ":12000")
- `DB_ADDR`: PostgreSQL connection string

### Production Considerations
- Use environment-specific database credentials
- Enable HTTPS for production deployment
- Consider rate limiting for API endpoints
- Implement proper logging and monitoring
- Add authentication/authorization as needed

## 🧪 Testing

The project includes comprehensive testing of all major functionality:

### Manual Testing Completed
- ✅ User creation and retrieval
- ✅ Post creation and listing
- ✅ Comment creation with AI responses
- ✅ Database relationships and constraints
- ✅ Frontend integration
- ✅ Real-time AI response generation

### API Testing Examples
```bash
# Health check
curl http://localhost:12000/v1/health

# Get all posts with comments
curl http://localhost:12000/v1/posts

# Get specific post comments
curl http://localhost:12000/v1/posts/1/comments
```

## 🎭 The "Dead Internet" Experience

This project demonstrates several key aspects of the Dead Internet Theory:

1. **Immediate Responses**: Every comment gets an instant reply, suggesting non-human response times
2. **Generic Interactions**: AI responses are contextually appropriate but somewhat generic
3. **Endless Engagement**: The system never stops responding, creating artificial activity
4. **Subtle Uncanny Valley**: Responses feel almost human but with a subtle artificial quality

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is open source and available under the MIT License.

## 🔮 Future Enhancements

- **Advanced AI**: Integration with GPT or other LLMs for more sophisticated responses
- **User Authentication**: JWT-based authentication system
- **Real-time Updates**: WebSocket support for live comment updates
- **Content Moderation**: Automatic content filtering
- **Analytics Dashboard**: Metrics on human vs AI interactions
- **Mobile App**: React Native or Flutter mobile application
- **Sentiment Analysis**: AI responses based on comment sentiment
- **Bot Detection**: Tools to identify AI-generated content

---

**⚠️ Disclaimer**: This project is for educational and demonstration purposes, exploring the concept of artificial online interactions. It's not intended to deceive users about the nature of AI-generated content.
