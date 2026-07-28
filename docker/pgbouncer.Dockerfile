FROM pgbouncer/pgbouncer:latest

# Copy configuration files
COPY configs/pgbouncer.ini /etc/pgbouncer/pgbouncer.ini
COPY configs/pgbouncer_userlist.txt /etc/pgbouncer/userlist.txt

# Set permissions
RUN chmod 640 /etc/pgbouncer/userlist.txt

# Create log directory
RUN mkdir -p /var/log/pgbouncer && \
    chown -R pgbouncer:pgbouncer /var/log/pgbouncer

# Expose pgBouncer port
EXPOSE 6432

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD pgbouncer --show-config || exit 1

# Run pgBouncer
CMD ["pgbouncer", "/etc/pgbouncer/pgbouncer.ini"]
