#!/bin/sh

# Backup script for PostgreSQL database
# This script creates a backup and uploads it to Cloudflare R2

set -e

# Configuration
DB_NAME=${POSTGRES_DB:-starehian_db}
DB_USER=${POSTGRES_USER:-starehian_user}
DB_PASSWORD=${POSTGRES_PASSWORD}
DB_HOST=${POSTGRES_HOST:-postgres}
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/starehian_backup_${DATE}.sql"
R2_BUCKET=${R2_BUCKET}
R2_ACCESS_KEY=${R2_ACCESS_KEY}
R2_SECRET_KEY=${R2_SECRET_KEY}
R2_ENDPOINT=${R2_ENDPOINT}

# Create backup directory
mkdir -p ${BACKUP_DIR}

# Create database backup
echo "Creating database backup..."
PGPASSWORD=${DB_PASSWORD} pg_dump -h ${DB_HOST} -U ${DB_USER} ${DB_NAME} > ${BACKUP_FILE}

# Compress backup
echo "Compressing backup..."
gzip ${BACKUP_FILE}
BACKUP_FILE="${BACKUP_FILE}.gz"

# Upload to R2 using AWS CLI (if available)
if command -v aws &> /dev/null; then
    echo "Uploading backup to R2..."
    AWS_ACCESS_KEY_ID=${R2_ACCESS_KEY} \
    AWS_SECRET_ACCESS_KEY=${R2_SECRET_KEY} \
    aws --endpoint-url ${R2_ENDPOINT} \
        s3 cp ${BACKUP_FILE} s3://${R2_BUCKET}/backups/
    
    echo "Backup uploaded successfully"
else
    echo "AWS CLI not found, skipping R2 upload"
    echo "Backup file: ${BACKUP_FILE}"
fi

# Clean up old backups (keep last 7 days)
echo "Cleaning up old backups..."
find ${BACKUP_DIR} -name "starehian_backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed successfully"
