# Old Starehian Society Platform

A comprehensive alumni community platform built with Go, featuring social networking, real-time chat, and admin management capabilities.

## Tech Stack

- **Backend**: Go (Fiber framework)
- **Database**: PostgreSQL 16
- **Cache/Session**: Redis (with AOF persistence)
- **Real-time Chat**: Centrifugo
- **Object Storage**: Cloudflare R2
- **Logging**: Loki + Promtail
- **Monitoring**: Grafana + Prometheus
- **Containerization**: Docker Compose

## Features

### Authentication & User Management
- Phone OTP authentication via Africa's Talking
- JWT + refresh token authentication
- Admin login with email/password
- Session management with Redis
- User role system (super_admin, moderator, support, member)
- User status management (active, pending, suspended, deleted)

### Profiles & Privacy
- Profile management with full CRUD
- Privacy controls (public, connections, private)
- Per-section visibility (profile, contact, career)
- Alumni directory search with filters
- Alumni roster management
- CSV import tool for bulk alumni data

### Social Features
- Connection requests (send, accept, reject)
- Connection management (view, remove)
- User blocking and unblocking
- Feed with posts and media
- Post creation with text and images
- Comments on posts
- Reactions (likes) on posts
- Post visibility controls

### Real-time Chat
- 1:1 direct messaging
- Group conversations
- Message history
- Read receipts
- Online presence
- Conversation member management
- WebSocket support via Centrifugo

### Notifications
- Real-time notifications
- Notification types (connection, post, comment, announcement)
- Unread count tracking
- Mark as read functionality

### Admin Portal
- User management (view, suspend, activate, delete)
- Role management
- Content moderation (reports)
- Audit log for all admin actions
- Analytics dashboard
  - User metrics (total, active, DAU, MAU)
  - Engagement metrics (posts, messages)
  - Cohort analytics by class year
  - Time-series data (signups, engagement)
  - Top content tracking
- Broadcast announcements
  - Target by all, class year, house, location
- Bulk operations
  - Bulk suspend/activate/verify/delete
  - Max 100 users per operation

### Infrastructure
- Rate limiting (per endpoint, per user/IP)
- Centralized authorization for privacy
- Health and readiness checks
- Comprehensive logging (Loki + Promtail)
- Monitoring (Grafana)
- Automated backups to R2
- WAL archiving for point-in-time recovery

## Project Structure

```
starehian-society-platform/
├── cmd/
│   ├── api/              # Main API server
│   └── import-roster/    # CSV import tool
├── internal/
│   ├── admin/            # Admin handlers and bulk operations
│   ├── auth/             # Authentication services (OTP, JWT, Africa's Talking)
│   ├── chat/             # Chat handlers and Centrifugo integration
│   ├── connections/      # Connection handlers
│   ├── middleware/       # Auth, authorization, rate limiting middleware
│   ├── models/           # Data models (User, Profile, Post, Chat, etc.)
│   ├── notifications/    # Notification handlers
│   ├── posts/            # Post handlers and media upload
│   ├── profiles/         # Profile handlers
│   ├── repository/       # Database repositories
│   └── services/         # Business logic services
├── migrations/           # Database migrations
├── pkg/                  # Shared packages
│   ├── config/           # Configuration
│   ├── database/         # Database connection
│   ├── logger/           # Logging
│   ├── ratelimit/        # Rate limiting
│   ├── redis/            # Redis client
│   └── storage/          # Cloudflare R2 storage
├── configs/              # Configuration files (Centrifugo, Loki, Promtail)
├── docker/               # Dockerfiles
├── scripts/              # Utility scripts (backup)
└── docker-compose.yml    # Docker orchestration
```

## Setup Instructions

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Cloudflare R2 account
- Africa's Talking account
- Firebase account (for FCM)

### Environment Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Edit `.env` with your actual values:
```env
POSTGRES_PASSWORD=your_secure_password
JWT_SECRET=your_jwt_secret
REFRESH_TOKEN_SECRET=your_refresh_token_secret
AFRICAS_TALKING_API_KEY=your_api_key
AFRICAS_TALKING_USERNAME=your_username
CLOUDFLARE_R2_ACCESS_KEY=your_access_key
CLOUDFLARE_R2_SECRET_KEY=your_secret_key
CLOUDFLARE_R2_BUCKET=your_bucket_name
CLOUDFLARE_R2_ENDPOINT=your_r2_endpoint
FCM_SERVER_KEY=your_fcm_server_key
CENTRIFUGO_SECRET=your_centrifugo_secret
CENTRIFUGO_API_KEY=your_centrifugo_api_key
CENTRIFUGO_ADMIN_PASSWORD=your_admin_password
CENTRIFUGO_ADMIN_SECRET=your_admin_secret
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=your_grafana_password
```

### Running with Docker Compose

1. Start all services:
```bash
docker-compose up -d
```

2. Check service status:
```bash
docker-compose ps
```

3. View logs:
```bash
docker-compose logs -f api
```

4. Stop services:
```bash
docker-compose down
```

### Database Migrations

The database schema is automatically applied when the PostgreSQL container starts. The migration files are located in the `migrations/` directory.

To manually run migrations:
```bash
# Connect to the PostgreSQL container
docker-compose exec postgres psql -U starehian_user -d starehian_db

# Run the migration manually
\i /docker-entrypoint-initdb.d/000001_init_schema.up.sql
```

### Import Alumni Roster

Use the CSV import tool to bulk import alumni data:

1. Create a CSV file with the following columns:
```csv
full_name,class_year,house,phone,email
John Doe,2010,Nairobi,+254712345678,john@example.com
Jane Smith,2015,Mombasa,+254798765432,jane@example.com
```

2. Run the import tool:
```bash
go run cmd/import-roster/main.go roster.csv "postgres://starehian_user:password@localhost:5432/starehian_db?sslmode=disable"
```

### Creating the First Admin User

1. Start the API server
2. Create a regular user via the signup endpoint
3. Manually update the user's role in the database:
```sql
UPDATE users SET role = 'super_admin' WHERE phone = '+254XXXXXXXXX';
```
4. Set a password for admin login:
```bash
# Use the SetPassword endpoint or update directly in database
```

## API Endpoints

### Authentication (Public)
- `POST /api/auth/request-otp` - Request OTP for phone number
- `POST /api/auth/signup` - Sign up with OTP
- `POST /api/auth/login` - Login with OTP
- `POST /api/auth/admin/login` - Admin login with email/password
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout

### Profiles (Protected)
- `GET /api/profiles/me` - Get own profile
- `PUT /api/profiles/me` - Update own profile
- `GET /api/profiles/:id` - Get user profile (with privacy checks)
- `POST /api/profiles/search` - Search alumni directory

### Connections (Protected)
- `POST /api/connections/` - Send connection request
- `POST /api/connections/:id/accept` - Accept connection request
- `POST /api/connections/:id/reject` - Reject connection request
- `DELETE /api/connections/:id` - Remove connection
- `POST /api/connections/:id/block` - Block user
- `DELETE /api/connections/:id/block` - Unblock user
- `GET /api/connections/` - Get connections
- `GET /api/connections/pending` - Get pending requests
- `GET /api/connections/sent` - Get sent requests
- `GET /api/connections/blocked` - Get blocked users

### Posts (Protected)
- `POST /api/posts/` - Create post
- `GET /api/posts/feed` - Get feed
- `GET /api/posts/:id` - Get post
- `PUT /api/posts/:id` - Update post
- `DELETE /api/posts/:id` - Delete post
- `GET /api/posts/user/:id` - Get user's posts
- `POST /api/posts/:id/comments` - Create comment
- `GET /api/posts/:id/comments` - Get comments
- `DELETE /api/posts/:id/comments/:commentId` - Delete comment
- `POST /api/posts/:id/reactions` - Create reaction
- `GET /api/posts/:id/reactions` - Get reactions
- `DELETE /api/posts/:id/reactions` - Delete reaction
- `GET /api/posts/:id/reaction` - Get user's reaction

### Media Upload (Protected)
- `POST /api/upload/media` - Upload single media file
- `POST /api/upload/media/multiple` - Upload multiple media files

### Chat (Protected)
- `POST /api/chat/direct/:id` - Create direct conversation
- `POST /api/chat/group` - Create group conversation
- `POST /api/chat/:id/messages` - Send message
- `GET /api/chat/:id/messages` - Get messages
- `GET /api/chat/` - Get conversations
- `POST /api/chat/:id/read` - Mark as read
- `POST /api/chat/:id/members` - Add member
- `DELETE /api/chat/:id/members/:memberId` - Remove member
- `DELETE /api/chat/:id` - Leave conversation
- `GET /api/chat/token` - Get WebSocket connection token
- `GET /api/chat/channel-token` - Get channel subscription token

### Notifications (Protected)
- `GET /api/notifications/` - Get notifications
- `GET /api/notifications/unread-count` - Get unread count
- `POST /api/notifications/:id/read` - Mark as read
- `POST /api/notifications/read-all` - Mark all as read
- `DELETE /api/notifications/:id` - Delete notification

### Admin (Admin Only)
- `GET /api/admin/users` - List users
- `GET /api/admin/users/:id` - Get user details
- `POST /api/admin/users/:id/suspend` - Suspend user
- `POST /api/admin/users/:id/activate` - Activate user
- `PUT /api/admin/users/:id/role` - Update user role
- `POST /api/admin/reports` - Create report
- `GET /api/admin/reports` - List reports
- `PUT /api/admin/reports/:id` - Update report status
- `GET /api/admin/audit` - View audit log
- `POST /api/admin/roster` - Add roster entry
- `GET /api/admin/roster` - List alumni roster
- `POST /api/admin/users/:id/verify` - Verify user against roster
- `GET /api/admin/analytics/dashboard` - Get dashboard metrics
- `GET /api/admin/analytics/cohorts` - Get cohort analytics
- `GET /api/admin/analytics/signups` - Get signup trends
- `GET /api/admin/analytics/engagement` - Get engagement trends
- `GET /api/admin/analytics/top-content` - Get top content
- `POST /api/admin/broadcasts` - Send broadcast announcement
- `POST /api/admin/bulk/suspend` - Bulk suspend users
- `POST /api/admin/bulk/activate` - Bulk activate users
- `POST /api/admin/bulk/verify` - Bulk verify users
- `POST /api/admin/bulk/delete` - Bulk delete users

### Health & Monitoring
- `GET /health` - Health check (database, Redis status)
- `GET /ready` - Readiness check

## Security Features

- **Rate Limiting**: Multi-tier rate limiting (auth, general, upload, admin endpoints)
- **Privacy Controls**: Centralized authorization function for profile/post visibility
- **Audit Logging**: All admin actions are logged immutably
- **Session Management**: Refresh tokens stored in Redis with expiration
- **Password Hashing**: Bcrypt for admin passwords
- **CORS**: Configured for frontend integration
- **Block System**: Users can block others to prevent interactions
- **Connection Verification**: Chat requires mutual connection

## Monitoring

- **Grafana**: Available at `http://localhost:3001` (admin/changeme)
- **Loki**: Centralized log aggregation
- **Promtail**: Docker log shipping
- **Uptime Kuma**: Can be added for uptime monitoring

## Backup Strategy

Automated backups are configured via the backup script:
- Nightly pg_dump to Cloudflare R2
- WAL archiving for point-in-time recovery
- 7-day retention policy

To manually trigger a backup:
```bash
docker-compose run backup
```

## Development

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Run the API server:
```bash
go run cmd/api/main.go
```

### Adding New Features

1. Create models in `internal/models/`
2. Create repository in `internal/repository/`
3. Create service in `internal/services/`
4. Create handlers in `internal/[module]/handlers.go`
5. Add routes in `cmd/api/main.go`

## Testing

To test the API:

1. Request OTP:
```bash
curl -X POST http://localhost:3000/api/auth/request-otp \
  -H "Content-Type: application/json" \
  -d '{"phone": "+254712345678"}'
```

2. Signup (use the OTP received):
```bash
curl -X POST http://localhost:3000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"phone": "+254712345678", "full_name": "John Doe", "otp": "123456"}'
```

3. Access protected endpoint:
```bash
curl -X GET http://localhost:3000/api/profiles/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Troubleshooting

### Database Connection Issues
- Check if PostgreSQL is running: `docker-compose ps postgres`
- Verify connection string in `.env`
- Check PostgreSQL logs: `docker-compose logs postgres`

### Redis Connection Issues
- Check if Redis is running: `docker-compose ps redis`
- Verify Redis URL in `.env`
- Check Redis logs: `docker-compose logs redis`

### OTP Not Sending
- Verify Africa's Talking credentials
- Check rate limits (max 3 OTP/hour per phone)
- Check API logs: `docker-compose logs api`

## License

Proprietary - Old Starehian Society

## Support

For issues and questions, contact the development team.
