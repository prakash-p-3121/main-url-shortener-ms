#!/bin/bash


IMAGE_NAME="prakashp92/main-url-shortener-ms"
IMAGE_TAG="latest"

# Build the image
sudo docker build -t "$IMAGE_NAME:$IMAGE_TAG" .

echo "Successfully built image: $IMAGE_NAME:$IMAGE_TAG"

echo "Pushing image to registry..."
sudo docker push "$IMAGE_NAME:$IMAGE_TAG"

