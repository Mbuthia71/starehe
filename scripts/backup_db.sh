#!/bin/bash

# PostgreSQL Backup Script for Starehe Society Platform
# This script performs daily database backups with retention policy

set -e

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-starehe_prod}"
DB_USER="${DB_USER:-postgres}"
BACKUP_DIR="${BACKUP_DIR:-/var/backups/postgresql}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/${DB_NAME}_${TIMESTAMP}.sql.gz"
LOG_FILE="${BACKUP_DIR}/backup_${TIMESTAMP}.log"

# Create backup directory if it doesn't exist
mkdir -p "${BACKUP_DIR}"

# Log function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "${LOG_FILE}"
}

log "Starting database backup for ${DB_NAME}"

# Check if pg_dump is available
if ! command -v pg_dump &> /dev/null; then
    log "ERROR: pg_dump not found. Please install PostgreSQL client tools."
    exit 1
fi

# Perform backup
log "Running pg_dump..."
if PGPASSWORD="${DB_PASSWORD}" pg_dump -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -F p --verbose 2>&1 | tee -a "${LOG_FILE}" | gzip > "${BACKUP_FILE}"; then
    log "Backup completed successfully: ${BACKUP_FILE}"
    
    # Get backup size
    BACKUP_SIZE=$(du -h "${BACKUP_FILE}" | cut -f1)
    log "Backup size: ${BACKUP_SIZE}"
    
    # Verify backup file exists and is not empty
    if [ -s "${BACKUP_FILE}" ]; then
        log "Backup file verified successfully"
    else
        log "ERROR: Backup file is empty or corrupted"
        exit 1
    fi
else
    log "ERROR: Backup failed"
    exit 1
fi

# Clean up old backups
log "Cleaning up backups older than ${RETENTION_DAYS} days..."
find "${BACKUP_DIR}" -name "${DB_NAME}_*.sql.gz" -type f -mtime +${RETENTION_DAYS} -delete
find "${BACKUP_DIR}" -name "backup_*.log" -type f -mtime +${RETENTION_DAYS} -delete

# Count remaining backups
BACKUP_COUNT=$(find "${BACKUP_DIR}" -name "${DB_NAME}_*.sql.gz" -type f | wc -l)
log "Remaining backups: ${BACKUP_COUNT}"

log "Backup process completed successfully"

# Optional: Upload to S3/R2 if configured
if [ -n "${S3_BUCKET}" ] && [ -n "${AWS_ACCESS_KEY_ID}" ]; then
    log "Uploading backup to S3..."
    if command -v aws &> /dev/null; then
        aws s3 cp "${BACKUP_FILE}" "s3://${S3_BUCKET}/database-backups/" --storage-class STANDARD_IA
        log "S3 upload completed"
    else
        log "WARNING: AWS CLI not found, skipping S3 upload"
    fi
fi

exit 0
