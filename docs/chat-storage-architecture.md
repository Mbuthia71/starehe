# Chat Storage Architecture for 6,000 Alumni

## Overview
This document outlines the storage architecture for supporting real-time chat among ~6,000 Starehe alumni.

## Current Infrastructure
- **Database**: PostgreSQL (primary data storage)
- **Cache**: Redis (sessions, rate limiting)
- **Backend**: Go Fiber API
- **Frontend**: React with TanStack Router

## Storage Requirements Analysis

### Message Volume Estimation
- **Users**: 6,000
- **Average messages per user per day**: 20 (conservative estimate)
- **Total daily messages**: 120,000
- **Total monthly messages**: 3.6 million
- **Message size**: ~500 bytes average (text + metadata)
- **Daily storage**: ~60 MB
- **Monthly storage**: ~1.8 GB
- **Yearly storage**: ~21.6 GB

### Peak Load Considerations
- **Peak hours**: 8 AM - 10 PM EAT
- **Peak message rate**: 2-3x average = 240,000-360,000 messages/day during peak
- **Concurrent active users**: ~1,500 (25% of total)
- **Messages per second during peak**: ~10-15 messages/second

## Proposed Architecture

### 1. Hybrid Storage Strategy

#### PostgreSQL (Primary Storage)
- **Purpose**: Persistent message storage, search, analytics
- **Table**: `messages` (already exists)
- **Retention**: 1 year full retention, older messages archived
- **Indexing**: 
  - `conversation_id` (for retrieval)
  - `created_at` (for time-based queries)
  - `sender_id` (for user history)
  - Full-text search on message content

#### Redis (Hot Cache & Real-time)
- **Purpose**: Recent messages cache, presence tracking, pub/sub
- **TTL**: 7 days for message cache
- **Data structures**:
  - `conversation:{id}:messages` (List, last 100 messages)
  - `user:{id}:presence` (Hash, online status, last seen)
  - `conversation:{id}:typing` (Set, typing indicators)
  - `pub/sub channels` for real-time delivery

### 2. Message Lifecycle

```
User sends message
    ↓
Store in Redis (immediate, TTL 7 days)
    ↓
Publish to conversation channel (real-time delivery)
    ↓
Async write to PostgreSQL (persistent storage)
    ↓
Update conversation metadata (last_message_at, message_count)
```

### 3. Storage Optimization

#### Message Compression
- Compress messages > 1KB using gzip
- Estimated savings: 30-40% storage reduction

#### Message Archival
- **Hot data**: Last 30 days in PostgreSQL + Redis
- **Warm data**: 30-365 days in PostgreSQL only
- **Cold data**: >365 days in compressed archive (S3/R2)
- **Cleanup**: Automated job to archive old messages

#### Conversation Partitioning
- Consider table partitioning by `created_at` (monthly partitions)
- Improves query performance and maintenance
- Easier to drop old partitions

### 4. Scaling Strategy

#### Horizontal Scaling
- **Read replicas**: 2 PostgreSQL read replicas for message retrieval
- **Redis cluster**: If message volume exceeds single Redis capacity
- **Connection pooling**: Optimize database connections

#### Caching Layers
- **L1**: Redis (messages, presence)
- **L2**: Application-level cache (conversation metadata)
- **CDN**: For static assets (media in messages)

### 5. Real-Time Delivery

#### WebSocket/Server-Sent Events
- Use existing WebSocket infrastructure
- Fallback to SSE if WebSocket unavailable
- Reconnection strategy with exponential backoff

#### Message Queue (Optional)
- For high-volume scenarios, consider message queue (RabbitMQ/Redis Streams)
- Decouple message sending from storage
- Handle backpressure during peak loads

## Implementation Plan

### Phase 1: Immediate (Current Setup)
- Continue using PostgreSQL for message storage
- Implement Redis caching for recent messages
- Add presence tracking in Redis
- Optimize existing message queries

### Phase 2: Optimization (1-2 weeks)
- Implement message compression
- Add message archival job
- Set up read replicas
- Implement conversation partitioning

### Phase 3: Scaling (As needed)
- Evaluate Redis cluster requirement
- Consider message queue for high-volume scenarios
- Implement advanced caching strategies

## Monitoring & Alerts

### Key Metrics
- **Message rate**: Messages per second/minute/hour
- **Storage growth**: Database size trends
- **Cache hit rate**: Redis cache effectiveness
- **Latency**: Message delivery time
- **Error rate**: Failed message deliveries

### Alerts
- High message rate (>50 messages/second sustained)
- Database storage approaching limit (>80%)
- Redis memory usage (>80%)
- High message delivery latency (>2s p95)

## Cost Estimation

### Storage Costs (Year 1)
- **PostgreSQL**: ~25 GB (including other data)
- **Redis**: ~5 GB (hot cache)
- **Archive (R2)**: ~15 GB (archived messages)
- **Total**: ~45 GB

### Bandwidth Costs
- **WebSocket traffic**: Minimal (text messages)
- **Media uploads**: Separate from chat storage
- **Estimated**: < 10 GB/month

## Recommendations

1. **Start with hybrid approach**: PostgreSQL + Redis caching
2. **Implement monitoring early**: Track message volume and patterns
3. **Archive proactively**: Don't let database grow indefinitely
4. **Optimize queries**: Ensure efficient message retrieval
5. **Plan for growth**: Architecture should scale to 10,000+ users

## Current Chat Schema Review

The existing `messages` table structure:
```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(id),
    sender_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    media_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Recommended additions**:
- `is_archived` BOOLEAN (for archival tracking)
- `compressed_content` BYTEA (for compressed messages)
- `message_type` VARCHAR (text, image, video, etc.)
- `reply_to_id` UUID (for threaded replies)
- `read_receipts` JSONB (for read status tracking)

## Conclusion

The proposed hybrid storage architecture using PostgreSQL for persistence and Redis for hot caching will efficiently support 6,000 alumni chatting. The system is designed to scale as the user base grows while maintaining performance and cost-effectiveness.
