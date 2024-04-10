#!/bin/bash

# from https://joshblais.com/posts/How-I-deploy-apps

# Check the number of arguments
if [ "$#" -ne 3 ]; then
echo "Usage: $0 <image_name> <port_mapping> <domain>"
exit 1
fi

# Assign arguments to variables
IMAGE_NAME=$1
PORT_MAPPING=$2
DOMAIN=$3
EMAIL=""
TAG="latest"  # Adjust tagging strategy as needed
DOCKER_REGISTRY=""  # Docker registry URL
SERVER_HOST=""
SSH_USER=""  # SSH user on the remote server
SSH_KEY_PATH=""  # SSH private key path

# Extract host port from PORT_MAPPING
HOST_PORT=$(echo $PORT_MAPPING | cut -d':' -f1)

# Build the Docker image with Buildx
docker buildx build --platform linux/amd64 -t $DOCKER_REGISTRY/"$IMAGE_NAME":$TAG --load .

# Push the Docker image
docker push $DOCKER_REGISTRY/$IMAGE_NAME:$TAG

# SSH into server to pull the image, restart the container, and configure NGINX and Certbot
ssh -i $SSH_KEY_PATH $SSH_USER@$SERVER_HOST << EOF
# Pull the latest Docker image
docker pull $DOCKER_REGISTRY/$IMAGE_NAME:$TAG

# Stop and remove the existing container if it exists
docker stop $IMAGE_NAME || true
docker rm $IMAGE_NAME || true

# Run the new container in the background with the specified port mapping
docker run -d --name $IMAGE_NAME -p $PORT_MAPPING $DOCKER_REGISTRY/$IMAGE_NAME:$TAG

# Check if NGINX config exists, if not, create it
NGINX_CONFIG="/etc/nginx/sites-available/$DOMAIN.conf"
NGINX_ENABLED="/etc/nginx/sites-enabled/$DOMAIN.conf"

if [ ! -f "\$NGINX_CONFIG" ]; then
    sudo bash -c "cat > \$NGINX_CONFIG" << 'ENDOFFILE'
server {
    listen 80;
    server_name $DOMAIN;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name $DOMAIN;

    ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;

    location / {
        proxy_pass http://localhost:$HOST_PORT;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
ENDOFFILE

    sudo ln -s "\$NGINX_CONFIG" "\$NGINX_ENABLED"
fi

# Reload NGINX to apply configuration
sudo nginx -t && sudo systemctl reload nginx

# After NGINX is reloaded:
sudo certbot certonly --webroot -w /var/www/certbot -d $DOMAIN --email $EMAIL --agree-tos --non-interactive --deploy-hook "sudo systemctl reload nginx"

# Reload NGINX to use new SSL certificate
sudo systemctl reload nginx
EOF

echo "Deployment complete."
