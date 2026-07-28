#!/bin/bash

# PostgreSQL Restore Script for Starehe Society Platform
# This script restores a database from a backup file

set -e

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-starehe_prod}"
DB_USER="${DB_USER:-postgres}"
BACKUP_DIR="${BACKUP_DIR:-/var/backups/postgresql}"
BACKUP_FILE="${1}"

# Check if backup file is provided
if [ -z "${BACKUP_FILE}" ]; then
    echo "Usage: $0 <backup_file.sql.gz>"
    echo "Available backups in ${BACKUP_DIR}:"
    ls -lh "${BACKUP_DIR}"/*.sql.gz 2>/dev/null || echo "No backups found"
    exit 1
fi

# Check if backup file exists
if [ ! -f "${BACKUP_FILE}" ]; then
    echo "ERROR: Backup file not found: ${BACKUP_FILE}"
    exit 1
fi

# Log function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

log "Starting database restore from ${BACKUP_FILE}"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    log "ERROR: psql not found. Please install PostgreSQL client tools."
    exit 1
fi

# Confirm restore
read -p "This will DROP and RECREATE the database '${DB_NAME}'. Are you sure? (yes/no): " CONFIRM
if [ "${CONFIRM}" != "yes" ]; then
    log "Restore cancelled by user"
    exit 0
fi

# Perform restore
log "Decompressing and restoring backup..."
if gunzip -c "${BACKUP_FILE}" | PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d postgres; then
    log "Restore completed successfully"
else
    log "ERROR: Restore failed"
    exit 1
fi

# Verify restore
log "Verifying restore..."
if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c "SELECT COUNT(*) FROM users;" > /dev/null 2>&1; then
    log "Restore verification successful"
else
    log "WARNING: Restore verification failed"
fi

log "Restore process completed"

exit 0
