# 🚀 Deploying SPOTIFY_MUSIC

This guide walks you through every deployment option available — from a quick local run to a production Docker setup.

---

## 📋 What You'll Need

Before deploying, have these ready:

| Requirement | Where to Get It |
|-------------|----------------|
| Telegram `API_ID` & `API_HASH` | [my.telegram.org](https://my.telegram.org) |
| Bot Token | [@BotFather](https://t.me/BotFather) on Telegram |
| MongoDB URI | [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) (free tier works) |
| Session String | Use a session generator bot |

---

## ⚙️ Environment Setup

**1. Clone the repo:**
```sh
git clone https://github.com/BabiesIQ/SPOTIFY_MUSIC.git
cd SPOTIFY_MUSIC
```

**2. Create your config file:**
```sh
cp sample.env .env
```

**3. Fill in your credentials:**

```env
API_ID=your_api_id
API_HASH=your_api_hash
TOKEN=your_bot_token
MONGO_URI=mongodb+srv://...
STRING1=your_session_string
OWNER_ID=your_telegram_id
LOGGER_ID=-100your_log_group_id
```

See `sample.env` for the full list of available options.

---

## 🐳 Docker (Recommended)

The easiest way to run SPOTIFY_MUSIC in production.

**Requirements:** [Docker](https://docs.docker.com/get-docker/) installed.

```sh
# 1. Build the image
docker build -t spotify-music-bot .

# 2. Run it
docker run -d \
  --name spotify-music-bot \
  --env-file .env \
  --restart unless-stopped \
  spotify-music-bot
```

**Useful commands:**

```sh
# View live logs
docker logs -f spotify-music-bot

# Stop the bot
docker stop spotify-music-bot

# Restart after a code update
docker stop spotify-music-bot && docker rm spotify-music-bot
git pull origin main
docker build -t spotify-music-bot .
docker run -d --name spotify-music-bot --env-file .env --restart unless-stopped spotify-music-bot
```

Or use Docker Compose (recommended for multi-service setups):

```sh
docker compose up -d --build
```

---

## 🔧 Manual — Linux / macOS

**Requirements:** Go 1.25+, FFmpeg

```sh
# Install on Debian/Ubuntu
sudo apt-get update && sudo apt-get install -y golang ffmpeg unzip

# Install on macOS (Homebrew)
brew install go ffmpeg
```

**Build & run:**

```sh
# Download required native libraries
go run setup_ntgcalls.go
go run github.com/BabiesIQ/gotdbot/scripts/tools@latest

# Install Go dependencies
go mod tidy

# Build binary
go build -o spotify-music .

# Run
./spotify-music
```

**Keep it running (systemd):**

Create `/etc/systemd/system/spotify-music.service`:

```ini
[Unit]
Description=SPOTIFY_MUSIC Telegram Bot
After=network.target

[Service]
Type=simple
WorkingDirectory=/path/to/SPOTIFY_MUSIC
EnvironmentFile=/path/to/SPOTIFY_MUSIC/.env
ExecStart=/path/to/SPOTIFY_MUSIC/spotify-music
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```sh
sudo systemctl daemon-reload
sudo systemctl enable spotify-music
sudo systemctl start spotify-music
```

---

## 🪟 Manual — Windows

**Requirements:** Go 1.25+, FFmpeg (added to PATH)

1. Install Go: [golang.org/doc/install](https://golang.org/doc/install)
2. Install FFmpeg: [ffmpeg.org/download.html](https://ffmpeg.org/download.html) → add to PATH
3. Clone the repo and set up `.env` as shown above

```powershell
# In PowerShell inside the project folder
go run setup_ntgcalls.go
go mod tidy
go build -o spotify-music.exe .
.\spotify-music.exe
```

---

## ☁️ Heroku

Click the button in the main README to deploy directly:

1. Fill in all required environment variables in the Heroku dashboard.
2. Click **Deploy**.
3. In the **Resources** tab: **enable** the `worker` dyno, **disable** the `web` dyno.

---

## ❓ Need Help?

Open an issue on [GitHub](https://github.com/BabiesIQ/SPOTIFY_MUSIC/issues) or join the support group linked in the README.
