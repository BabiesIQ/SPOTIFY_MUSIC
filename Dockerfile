# ---- Base image ----
FROM python:3.10-slim

# ---- Set working directory ----
WORKDIR /app

# ---- Copy all project files ----
COPY . /app

# ---- Install system dependencies ----
# git = for GitHub packages
# curl = to install Node.js
# ffmpeg = required for pytgcalls audio streaming
RUN apt-get update && apt-get install -y \
    git curl ffmpeg \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && node -v && npm -v \
    && pip install --upgrade pip \
    && pip install --no-cache-dir -r requirements.txt \
    && rm -rf /var/lib/apt/lists/*

# ---- Default command ----
CMD ["bash", "start"]
