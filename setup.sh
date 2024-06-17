#!/bin/bash

# This script creates a new user 'yampa', adds it to necessary groups,
# enables password-less sudo, and configures SSH key-based authentication.

# Ensure the script is running as root
if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root" 1>&2
    exit 1
fi

# Create a new user called 'yampa' with a home directory
useradd -m -s /bin/bash yampa

# Check if the user was created successfully
if [ $? -eq 0 ]; then
    echo "User 'yampa' has been created successfully."
else
    echo "Failed to create user 'yampa'."
    exit 2
fi

# Add 'yampa' user to the 'sudo' and 'docker' groups
usermod -aG sudo,docker yampa

# Ensure docker group exists or create it if necessary
if ! getent group docker >/dev/null; then
    groupadd docker
    usermod -aG docker yampa
fi

# Copy the authorized_keys from the root user to 'yampa'
mkdir -p /home/yampa/.ssh
cp /root/.ssh/authorized_keys /home/yampa/.ssh/authorized_keys
chown -R yampa:yampa /home/yampa/.ssh
chmod 700 /home/yampa/.ssh
chmod 600 /home/yampa/.ssh/authorized_keys

# Configure sudo to allow 'yampa' to run sudo commands without a password
echo "yampa ALL=(ALL) NOPASSWD: ALL" >/etc/sudoers.d/90-yampa-nopasswd

echo "Setup completed successfully."
