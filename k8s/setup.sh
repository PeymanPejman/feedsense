#!/bin/bash

# This script sets up the frontned deployment
# and service.

# Compute cwd
CWD="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"\

# Create the frontend deployment
echo "Creating deployment for frontend..."
kubectl apply -f $CWD/fs-fe.yaml

# Create frontend service
echo "Creating the frontend service..."
kubectl apply -f $CWD/service-fs-fe.yaml