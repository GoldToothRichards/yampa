########################################################################################

# Set the Base Image from a build-arg
ARG BASE_IMAGE=ubuntu:22.04
FROM ${BASE_IMAGE}

# Set the working directory
WORKDIR /yampa_streams

ARG POETRY_VERSION=1.8.2
ARG POETRY_HOME="/opt/poetry"

# Use bash shell
SHELL ["/bin/bash", "-c"]

# Make sure terminal does not wait at a prompt
ENV DEBIAN_FRONTEND=noninteractive

# Install system packages required at runtime
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    rm -rf /var/lib/apt/lists/*

#Add repo link
RUN add-apt-repository ppa:deadsnakes/ppa

# Install system packages required at runtime
RUN apt-get update && apt-get install --no-install-recommends -y \
    curl \
    build-essential \
    python3.10 \
    python3.10-distutils \
    python3.10-venv \
    python3-pip \
    python3-setuptools \
    && rm -rf /var/lib/apt/lists/*

# Install Poetry - respects $POETRY_VERSION & $POETRY_HOME
ENV POETRY_HOME=${POETRY_HOME} \
    POETRY_VERSION=${POETRY_VERSION}
RUN curl -sSL https://install.python-poetry.org/ | python3

# Add Poetry to the path
ENV PATH="$POETRY_HOME/bin:$PATH"
ENV PATH="/usr/bin:$PATH"

# Make sure terminal does not wait at a prompt. Set timezone.
ENV TZ=America/New_York DEBIAN_FRONTEND=noninteractive

# This system may have been minimized by removing packages and content that are
# not required on a system that users do not log into.
# To restore this content, including manpages, you can run the 'unminimize'
# command. You will still need to ensure the 'man-db' package is installed.
# This process is slow, so we do a separate step so that it can be cached.
RUN apt-get update && apt-get install --no-install-recommends -y \
    man-db \
    && rm -rf /var/lib/apt/lists/*
RUN yes | unminimize

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

# Install Miniconda
RUN mkdir -p ~/miniconda3 && \
    wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh -O ~/miniconda3/miniconda.sh && \
    bash ~/miniconda3/miniconda.sh -b -u -p /opt/miniconda3 && \
    rm -rf ~/miniconda3/miniconda.sh && \
    /opt/miniconda3/bin/conda init bash && \
    /opt/miniconda3/bin/conda init zsh

ENV PATH="/opt/miniconda3/bin:$PATH"

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

# Give the user ownership of the working directory
RUN chown -R $USERNAME:$USERNAME /yampa_streams

# Default to non-root user with sudo permissions
USER $USERNAME

# Create a folder to store command history so that it can be persisted through container restarts
RUN sudo mkdir /commandhistory \
    && sudo touch /commandhistory/.bash_history \
    && sudo touch /commandhistory/.zsh_history \
    && sudo chown -R $USERNAME:$USERNAME /commandhistory

########################################################################################