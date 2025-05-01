# Je Suis Ici - ActivityPub Social Network Platform

Je Suis Ici is a decentralized social network platform based on the ActivityPub protocol, allowing users to post check-in at locations and interact with other ActivityPub-compatible platforms.

## Features

- Decentralized social networking using ActivityPub protocol
- Location-based check-in functionality
- Media upload support
- JWT authentication
- Distributed tracing (using Jaeger)
- Federation with other ActivityPub platforms

## Technical Architecture

### Backend Stack

- Language: Go 1.23.0
- Database: CockroachDB (PostgreSQL compatible)
- Object Storage: MinIO
- Distributed Tracing: Jaeger
- API Router: Chi Router
- Authentication: JWT

### Project Structure

```
je-suis-ici-activitypub/
├── cmd/                    # Main application entrypoints
│   └── server/            # Server startup code
├── internal/              # Internal packages
│   ├── activitypub/      # ActivityPub protocol implementation
│   ├── api/              # API handlers and middleware
│   ├── config/           # Configuration management
│   ├── db/              # Database related code
│   ├── services/        # Business logic services
│   ├── storage/         # File storage services
│   └── tracing/         # Distributed tracing
```

### Database Schema

#### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    actor_id VARCHAR(255) NOT NULL UNIQUE,
    private_key TEXT,
    public_key TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### Checkins Table
```sql
CREATE TABLE checkins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    location_name VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    activity_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### Media Table
```sql
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    checkin_id UUID REFERENCES checkins(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INT NOT NULL,
    width INT,
    height INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### Activities Table
```sql
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    activity_id VARCHAR(255) NOT NULL UNIQUE,
    actor VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    object_id VARCHAR(255),
    object_type VARCHAR(50),
    target VARCHAR(255),
    raw_content JSONB NOT NULL,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### Followers Table
```sql
CREATE TABLE followers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    follower_actor_id VARCHAR(255) NOT NULL,
    follower_inbox VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, follower_actor_id)
);
```

### Core Components

1. **ActivityPub Implementation**
   - Actor Management
   - Activity Processing
   - Federated Social Interactions

2. **User System**
   - Registration/Login
   - Profile Management
   - ActivityPub Actor Association

3. **Check-in System**
   - Create Check-ins
   - Media Upload
   - Timeline Display

## API Endpoints

### Authentication API
- `POST /auth/register` - User Registration
- `POST /auth/login` - User Login

### User API
- `PUT /api/users/{id}` - Update User Profile
- `DELETE /api/users/{id}` - Delete User

### Check-in API
- `POST /api/media` - Upload Media
- `POST /api/checkins` - Create New Check-in
- `GET /api/checkins` - Get User Check-ins
- `GET /api/checkins/{id}` - Get Specific Check-in

### ActivityPub API
- `GET /.well-known/webfinger` - WebFinger Service
- `GET /.well-known/nodeinfo` - NodeInfo Service
- `GET /api/users/{username}/activitypub-info` - Get User ActivityPub Info
- `POST /api/users/{sender_username}/send-checkin` - Send Check-in to User
- `GET /api/users/{username}/inbox` - Get User Inbox

## Environment Variables

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=26260
DB_USER=root
DB_PASSWORD=
DB_NAME=checkin
DB_SSL_MODE=disable

# MinIO Configuration
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadminpassword
MINIO_BUCKET=checkin-media
MINIO_USE_SSL=false

# JWT Configuration
JWT_SECRET=your-jwt-secret

# Jaeger Configuration
JAEGER_URL=http://jaeger:14268/api/traces
JAEGER_SERVICE_NAME=je-suis-ici
JAEGER_ENVIRONMENT=development
JAEGER_ENABLE=true
```

## Development Setup

1. Copy the environment variables template:
```bash
cp .env.example .env
```

2. Configure the required environment variables

3. Start the required services:
- CockroachDB
- MinIO
- Jaeger (optional)

4. Run database migrations:
```bash
go run cmd/server/main.go migrate
```

5. Start the server:
```bash
go run cmd/server/main.go
```
