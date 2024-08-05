FROM arm64v8/alpine:latest

# Install required packages for downloading files
RUN apk add --no-cache curl

# Install Cloud SQL Proxy
RUN curl -o /usr/local/bin/cloud_sql_proxy \
    -L https://dl.google.com/cloudsql/cloud_sql_proxy.linux.arm64 \
    && chmod +x /usr/local/bin/cloud_sql_proxy

# Expose the port used by Cloud SQL Proxy
EXPOSE 3306

# Start Cloud SQL Proxy
CMD cloud_sql_proxy -instances=${DB_HOST_NAME}=tcp:0.0.0.0:3306 -credential_file=/secrets/cloudsql/serviceAccountKey.json

