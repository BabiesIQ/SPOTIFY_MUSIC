# ---- Base image ----
FROM python:3.10-slim

# ---- Working directory ----
WORKDIR /app

# ---- Copy project files ----
COPY . /app

# ---- Install dependencies ----
RUN pip install --no-cache-dir -r requirements.txt

# ---- Environment variables (can be overridden by Heroku/Render) ----
ENV PYTHONUNBUFFERED=1
ENV PORT=8080

# ---- Default command ----
CMD ["python", "app.py"]
