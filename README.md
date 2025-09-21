<div align="center">
    <img src='./.github/images/whoami-logo.png' height=100/>
</div>

# whoami

A Central Authentication Service written entirely in Go. The service is intended to be fully Open-Source and used as a separate service in your application. It's designed to be easily plug-and-play for microservices architectures and frontend applications that need authentication features instead of using services like Firebase or Supabase.

## 🚀 Features

- **Basic Username/Password Authentication** - Secure login with email/username and password
- **Email Verification** - Email-based account verification system
- **OAuth2 Integration** - Support for Google and GitHub OAuth providers
- **Password Reset Flows** - Secure password reset with HOTP-based OTP verification
- **Account Management** - Account activation/deactivation capabilities
- **Rate Limiting** - Per IP and per user rate limiting to prevent abuse
- **Security Monitoring** - Suspicious activity detection and logging
- **Password Security** - Password history tracking and strength requirements
- **HaveIBeenPwned Integration** - Check passwords against known breaches
- **Account Lockout** - Automatic account lockout after multiple failed attempts
- **Session Management** - Short-lived access tokens with longer refresh tokens
- **Audit Logging** - Comprehensive audit trail for all user actions
- **Device Management** - Track and manage user devices
- **Data Export** - GDPR-compliant data export functionality

## 🏗️ Architecture

### System Overview

```mermaid
graph TB
    subgraph "Frontend (React + Vite)"
        FE[React Frontend]
        FE --> |API Calls| API
    end

    subgraph "Backend (Go + Gin)"
        API[REST API]
        API --> |Rate Limiting| RL[Rate Limiter]
        API --> |Authentication| AUTH[Auth Middleware]
        API --> |Token Management| TM[Token Maker]
        API --> |Session Management| SM[Session Service]
    end

    subgraph "Services Layer"
        US[User Service]
        SS[Security Service]
        PS[Password Security Service]
        ES[Email Service]
        PRS[Password Reset Service]
        AS[Audit Service]
        OS[OAuth Service]
        DS[Device Service]
        DES[Data Export Service]
    end

    subgraph "Repositories Layer"
        UR[User Repository]
        SR[Security Repository]
        PR[Password Repository]
        ER[Email Repository]
        AR[Audit Repository]
        OR[OAuth Repository]
        DR[Device Repository]
        DER[Data Export Repository]
    end

    subgraph "External Services"
        DB[(PostgreSQL)]
        REDIS[(Redis)]
        SMTP[SMTP Server]
        GOOGLE[Google OAuth]
        GITHUB[GitHub OAuth]
        HIBP[HaveIBeenPwned API]
    end

    API --> US
    API --> SS
    API --> PS
    API --> ES
    API --> PRS
    API --> AS
    API --> OS
    API --> DS
    API --> DES

    US --> UR
    SS --> SR
    PS --> PR
    ES --> ER
    AS --> AR
    OS --> OR
    DS --> DR
    DES --> DER

    UR --> DB
    SR --> DB
    PR --> DB
    ER --> DB
    AR --> DB
    OR --> DB
    DR --> DB
    DER --> DB

    RL --> REDIS
    SM --> REDIS
    TM --> REDIS

    ES --> SMTP
    PRS --> SMTP
    OS --> GOOGLE
    OS --> GITHUB
    PS --> HIBP
```

### Database Schema

```mermaid
erDiagram
    users ||--o{ user_profiles : has
    users ||--o{ refresh_tokens : has
    users ||--o{ email_verifications : has
    users ||--o{ password_resets : has
    users ||--o{ password_history : has
    users ||--o{ login_attempts : has
    users ||--o{ account_lockouts : has
    users ||--o{ suspicious_activities : has
    users ||--o{ audit_logs : generates
    users ||--o{ user_devices : has
    users ||--o{ oauth_accounts : has
    users ||--o{ data_exports : requests

    users {
        bigint id PK
        varchar email UK
        varchar username UK
        varchar password_hash
        boolean email_verified
        boolean active
        varchar role
        jsonb privacy_settings
        timestamptz last_login_at
        timestamptz password_changed_at
        timestamptz created_at
        timestamptz updated_at
    }

    user_profiles {
        bigint id PK
        bigint user_id FK
        varchar first_name
        varchar last_name
        varchar phone
        varchar avatar_url
        text bio
        varchar timezone
        varchar locale
        timestamptz created_at
        timestamptz updated_at
    }

    refresh_tokens {
        bigint id PK
        bigint user_id FK
        varchar token_hash
        timestamptz expires_at
        timestamptz created_at
    }

    password_resets {
        bigint id PK
        bigint user_id FK
        varchar token_hash
        varchar hotp_secret
        integer counter
        timestamptz expires_at
        timestamptz used_at
        timestamptz created_at
    }

    oauth_accounts {
        bigint id PK
        bigint user_id FK
        varchar provider
        varchar provider_user_id
        varchar email
        varchar name
        varchar avatar_url
        jsonb provider_data
        timestamptz created_at
        timestamptz updated_at
    }

    audit_logs {
        bigint id PK
        bigint user_id FK
        varchar action
        varchar resource_type
        bigint resource_id
        inet ip_address
        text user_agent
        jsonb details
        timestamptz created_at
    }
```

## 🔄 Common Flows

### 1. User Registration Flow

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as API
    participant US as User Service
    participant ES as Email Service
    participant DB as Database
    participant SMTP as SMTP Server

    U->>F: Fill registration form
    F->>A: POST /api/v1/register
    A->>A: Rate limiting check
    A->>US: Create user
    US->>US: Validate password strength
    US->>US: Check HaveIBeenPwned
    US->>US: Hash password
    US->>DB: Save user (email_verified=false)
    US->>ES: Send verification email
    ES->>SMTP: Send email
    A->>F: 201 Created
    F->>U: Show success message
    SMTP->>U: Verification email
```

### 2. User Login Flow

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as API
    participant SS as Security Service
    participant US as User Service
    participant TM as Token Maker
    participant SM as Session Service
    participant DB as Database
    participant R as Redis

    U->>F: Enter credentials
    F->>A: POST /api/v1/login
    A->>A: Rate limiting check
    A->>SS: Check account lockout
    A->>US: Get user by email
    US->>DB: Query user
    A->>A: Verify password
    alt Invalid credentials
        A->>SS: Record failed login
        SS->>DB: Save login attempt
        A->>F: 401 Unauthorized
    else Valid credentials
        A->>SS: Record successful login
        SS->>DB: Save login attempt
        A->>TM: Generate access token
        A->>SM: Create session
        SM->>R: Store session
        A->>F: 200 OK + tokens
        F->>U: Redirect to dashboard
    end
```

### 3. Password Reset Flow

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as API
    participant PRS as Password Reset Service
    participant PS as Password Security Service
    participant ES as Email Service
    participant DB as Database
    participant SMTP as SMTP Server

    U->>F: Request password reset
    F->>A: POST /api/v1/password-reset/request
    A->>PRS: Request reset
    PRS->>DB: Get user by email
    PRS->>PRS: Generate reset token
    PRS->>PRS: Generate HOTP secret
    PRS->>DB: Save reset record
    PRS->>ES: Send reset email
    ES->>SMTP: Send email with token
    A->>F: 200 OK
    SMTP->>U: Reset email

    U->>F: Click reset link
    F->>A: POST /api/v1/password-reset/verify
    A->>PRS: Verify token
    PRS->>DB: Get reset record
    PRS->>PRS: Generate HOTP OTP
    PRS->>ES: Send OTP email
    ES->>SMTP: Send OTP
    SMTP->>U: OTP email

    U->>F: Enter OTP
    F->>A: POST /api/v1/password-reset/verify-otp
    A->>PRS: Verify OTP
    PRS->>PRS: Validate HOTP

    U->>F: Enter new password
    F->>A: POST /api/v1/password-reset/reset
    A->>PRS: Reset password
    PRS->>PS: Validate new password
    PRS->>PS: Update password
    PS->>DB: Update user password
    PRS->>DB: Mark reset as used
    A->>F: 200 OK
```

### 4. OAuth Login Flow

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as API
    participant OS as OAuth Service
    participant OP as OAuth Provider
    participant US as User Service
    participant TM as Token Maker
    participant DB as Database

    U->>F: Click OAuth login
    F->>A: GET /api/v1/oauth/login/google
    A->>OS: Generate OAuth state
    A->>OP: Get authorization URL
    A->>F: Return auth URL
    F->>U: Redirect to OAuth provider

    U->>OP: Authorize application
    OP->>F: Redirect with code
    F->>A: GET /api/v1/oauth/callback/google
    A->>OS: Validate state
    A->>OP: Exchange code for token
    OP->>A: Return user info
    A->>OS: Authenticate/create user
    OS->>DB: Check existing OAuth account
    alt New user
        OS->>US: Create new user
        US->>DB: Save user
        OS->>DB: Save OAuth account
    else Existing user
        OS->>DB: Get existing user
    end
    A->>TM: Generate tokens
    A->>F: Return tokens
    F->>U: Login successful
```

### 5. Session Management Flow

```mermaid
sequenceDiagram
    participant F as Frontend
    participant A as API
    participant AM as Auth Middleware
    participant TM as Token Maker
    participant SM as Session Service
    participant R as Redis

    F->>A: API request with token
    A->>AM: Validate token
    AM->>TM: Verify token signature
    AM->>SM: Check session
    SM->>R: Get session data
    alt Valid session
        SM->>AM: Return session
        AM->>A: Continue request
        A->>F: Return response
    else Invalid/expired session
        AM->>F: 401 Unauthorized
        F->>A: POST /api/v1/refresh
        A->>TM: Generate new tokens
        A->>SM: Update session
        SM->>R: Store new session
        A->>F: Return new tokens
    end
```

## 🛠️ Technology Stack

### Backend

- **Go 1.24.1** - Main programming language
- **Gin** - HTTP web framework
- **PostgreSQL** - Primary database
- **Redis** - Caching and session storage
- **PASETO** - Token generation and validation
- **HOTP** - One-time password generation
- **SQLC** - Type-safe SQL code generation
- **Golang-migrate** - Database migrations

### Frontend

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **TanStack Router** - File-based routing
- **TanStack Query** - Data fetching and caching
- **Zustand** - State management
- **Tailwind CSS** - Styling
- **Radix UI** - Component primitives
- **React Hook Form** - Form handling
- **Zod** - Schema validation

### Infrastructure

- **Docker & Docker Compose** - Containerization
- **Nginx** - Reverse proxy and static file serving
- **SSL/TLS** - HTTPS support

## 🚀 Quick Start

### Prerequisites

- Go 1.24.1+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 17+
- Redis 7+

### 1. Clone and Setup

```bash
git clone https://github.com/m1thrandir225/whoami.git
cd whoami
make setup
```

### 2. Configure Environment

Edit `deployment/.env` file with your configuration:

```bash
# Database
DB_SOURCE=postgres://whoami_user:secret@whoami-db:5432/whoami_db?ENABLE_TLS=false

# Redis
REDIS_URL=redis://whoami-redis:6379

# Email (SMTP)
SMTP_HOST=your-smtp-host
SMTP_PORT=587
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password

# OAuth Providers
GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret
GITHUB_OAUTH_CLIENT_ID=your-github-client-id
GITHUB_OAUTH_CLIENT_SECRET=your-github-client-secret

# Frontend
FRONTEND_URL=http://localhost:3000
VITE_BACKEND_URL=http://localhost:8080
```

### 3. Start Services

```bash
# Start all services with Docker Compose
make docker-up

# Apply database migrations
make migrate-up-docker

# View logs
make docker-logs
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080 or https://localhost:8443
- **Health Check**: http://localhost:8080/health

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint           | Description          | Rate Limit   |
| ------ | ------------------ | -------------------- | ------------ |
| POST   | `/api/v1/register` | User registration    | Registration |
| POST   | `/api/v1/login`    | User login           | Auth         |
| POST   | `/api/v1/refresh`  | Refresh access token | Auth         |
| POST   | `/api/v1/logout`   | User logout          | Default      |

### Password Reset Endpoints

| Method | Endpoint                            | Description            | Rate Limit     |
| ------ | ----------------------------------- | ---------------------- | -------------- |
| POST   | `/api/v1/password-reset/request`    | Request password reset | Password Reset |
| POST   | `/api/v1/password-reset/verify`     | Verify reset token     | Password Reset |
| POST   | `/api/v1/password-reset/verify-otp` | Verify OTP             | Password Reset |
| POST   | `/api/v1/password-reset/reset`      | Reset password         | Password Reset |

### OAuth Endpoints

| Method | Endpoint                           | Description          | Rate Limit |
| ------ | ---------------------------------- | -------------------- | ---------- |
| GET    | `/api/v1/oauth/login/:provider`    | Initiate OAuth login | Default    |
| GET    | `/api/v1/oauth/callback/:provider` | OAuth callback       | Default    |
| POST   | `/api/v1/oauth/exchange`           | Exchange temp token  | Default    |

### Protected Endpoints

| Method | Endpoint                       | Description               | Rate Limit |
| ------ | ------------------------------ | ------------------------- | ---------- |
| GET    | `/api/v1/me`                   | Get current user          | Default    |
| PUT    | `/api/v1/user/:id`             | Update user               | Default    |
| POST   | `/api/v1/user/update-password` | Update password           | Default    |
| GET    | `/api/v1/sessions`             | Get user sessions         | Default    |
| DELETE | `/api/v1/sessions/:token`      | Revoke session            | Default    |
| GET    | `/api/v1/security/activities`  | Get suspicious activities | Default    |
| GET    | `/api/v1/audit/recent`         | Get recent audit logs     | Default    |
| GET    | `/api/v1/devices`              | Get user devices          | Default    |
| POST   | `/api/v1/exports`              | Request data export       | Default    |

## 🔧 Development

### Local Development Setup

```bash
# Install dependencies
go mod tidy
cd frontend && pnpm install

# Start database and Redis
make docker-up

# Run migrations
make migrate-up

# Start backend
make server

# Start frontend (in another terminal)
cd frontend && pnpm dev
```

### Available Make Commands

```bash
make help                    # Show all available commands
make setup                   # Complete setup with env generation
make build                   # Build Go binary
make test                    # Run tests
make lint                    # Run linter
make docker-up               # Start services with Docker
make docker-down             # Stop services
make migrate-up-docker       # Apply migrations
make sqlc                    # Generate SQL code
```

### Database Migrations

```bash
# Create new migration
make migrate-create name=add_new_table

# Apply migrations
make migrate-up-docker

# Rollback migrations
make migrate-down-docker steps=1

# Check migration status
make migrate-status-docker
```

## 🔒 Security Features

### Rate Limiting

- **Registration**: 5 requests per hour per IP
- **Authentication**: 10 requests per hour per IP
- **Password Reset**: 3 requests per hour per IP
- **Default**: 100 requests per hour per user

### Password Security

- Minimum 8 characters
- Must contain uppercase, lowercase, number, and special character
- Checked against HaveIBeenPwned database
- Password history tracking (prevents reuse of last 5 passwords)

### Account Security

- Account lockout after 5 failed login attempts
- Suspicious activity detection and logging
- Device tracking and management
- Comprehensive audit logging

### Token Security

- PASETO tokens for stateless authentication
- Short-lived access tokens (15 minutes)
- Longer refresh tokens (7 days)
- Token blacklisting for secure logout

## 🚀 Deployment

### Production Deployment

```bash
# Build production images
make docker-build

# Set production environment
export ENVIRONMENT=production

# Start production services
make docker-up
```

### Environment Variables

See `deployment/.env` for complete configuration options including:

- Database configuration
- Redis configuration
- Email/SMTP settings
- OAuth provider credentials
- Security settings
- Frontend configuration

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

<div align="center">
    <p>Built with ❤️ by Sebastijan Zindl.</p>
</div>
