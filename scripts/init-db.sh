#!/bin/sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
while ! pg_isready -h postgres -p 5432 -U postgres; do
    sleep 1
done

# Check if the users table exists
echo "Checking if users table exists..."
if ! psql -h postgres -U postgres -d myorder -c "\dt users" | grep -q "users"; then
    echo "Creating users table..."
    psql -h postgres -U postgres -d myorder -f /docker-entrypoint-initdb.d/000001_create_users_table.up.sql
fi

echo "Database initialization complete!" 