FROM redis:latest

# Set working directory
WORKDIR /usr/local/etc/redis

# Copy sentinel configuration
COPY sentinel.conf .

# Set correct permissions so Redis Sentinel can write updates
RUN chmod 777 /usr/local/etc/redis && chmod 666 /usr/local/etc/redis/sentinel.conf
