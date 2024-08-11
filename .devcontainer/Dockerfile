########################################################################################

# Set the Base Image from a build-arg
ARG BASE_IMAGE=ubuntu:22.04
FROM ${BASE_IMAGE}

# This system may have been minimized by removing packages and content that are
# not required on a system that users do not log into.
# To restore this content, including manpages, you can run the 'unminimize'
# command. You will still need to ensure the 'man-db' package is installed.
# This process is slow, so we do a separate step so that it can be cached.
RUN apt-get update && apt-get install --no-install-recommends -y \
    man-db \
    && rm -rf /var/lib/apt/lists/*
RUN yes | unminimize

# Use bash shell
SHELL ["/bin/bash", "-c"]

# Make sure terminal does not wait at a prompt. Set timezone.
ENV TZ=America/New_York DEBIAN_FRONTEND=noninteractive

# Install system packages required at runtime
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    rm -rf /var/lib/apt/lists/*

#Add repo link
RUN add-apt-repository ppa:deadsnakes/ppa

# Add the latest git package
RUN add-apt-repository ppa:git-core/ppa

# Install system packages
RUN apt-get update && apt-get install --no-install-recommends -y \
    sudo \
    aptitude \
    zsh \
    fish \
    tmux \
    git \
    vim \
    neovim \
    nano \
    emacs \
    atop \
    htop \
    ctop \
    nvtop \
    tree \
    wget \
    fzf \
    rcm \
    xclip \
    tzdata \
    bat \
    pigz \
    ncdu \
    libicu-dev \
    zip \
    unzip \
    graphviz \
    iputils-ping \
    curl \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install RedPanda CLI
RUN curl -LO https://github.com/redpanda-data/redpanda/releases/latest/download/rpk-linux-amd64.zip && \
    mkdir -p /opt/rpk/bin && \
    unzip rpk-linux-amd64.zip -d /opt/rpk/bin

ENV PATH="/opt/rpk/bin:$PATH"

# Install Clickhouse
RUN GNUPGHOME=$(mktemp -d) && \
    GNUPGHOME="$GNUPGHOME" gpg --no-default-keyring --keyring /usr/share/keyrings/clickhouse-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8919F6BD2B48D754 && \
    rm -rf "$GNUPGHOME" && \
    chmod +r /usr/share/keyrings/clickhouse-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/clickhouse-keyring.gpg] https://packages.clickhouse.com/deb stable main" | sudo tee /etc/apt/sources.list.d/clickhouse.list && \
    apt-get update && \
    apt-get install -y clickhouse-server clickhouse-client

# Read user input
ARG USERNAME=yampa
ARG USER_UID=1003
ARG USER_GID=1002

# Create the user with sudo permissions
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME && \
    usermod -aG sudo $USERNAME && \
    mkdir /home/$USERNAME/workspace && \
    chown -R $USERNAME:$USERNAME /home/$USERNAME/workspace

# Default to non-root user with sudo permissions
USER $USERNAME

# Create a folder to store command history so that it can be persisted through container restarts
RUN sudo mkdir /commandhistory \
    && sudo touch /commandhistory/.bash_history \
    && sudo touch /commandhistory/.zsh_history \
    && sudo chown -R $USERNAME:$USERNAME /commandhistory

########################################################################################