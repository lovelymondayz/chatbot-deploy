# Chatbot Deploy — AI Chatbot Platform

A platform for deploying AI-powered chatbots with document-based knowledge retrieval.

## Quick Start

```bash
# Clone
git clone https://github.com/lovelymondayz/chatbot-deploy.git
cd chatbot-deploy

# Start all services
docker compose up -d --build

# Dashboard: http://localhost:8102
# API: http://localhost:8101
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        NGINX (80/443)                        │
│                   chatbot.arjism.com → :8102                │
├─────────────────────────────────────────────────────────────┤
│  React + Vite + TS + Tailwind  │  Python + FastAPI          │
│        (Dashboard :8102)       │       (Backend :8101)      │
├─────────────────────────────────────────────────────────────┤
│              ChromaDB (Vector Store)  │  OpenAI/9Router      │
└─────────────────────────────────────────────────────────────┘
```

## Features

- **Document Upload**: Upload PDFs, docs, and text files
- **RAG Chat**: Chat with AI using your documents as context
- **Vector Search**: Fast similarity search with ChromaDB
- **Bot Management**: Create and configure multiple bots
- **API Access**: RESTful API for integration
- **Dashboard**: Manage bots, documents, and view analytics

## API Endpoints

### Public
- `GET /api/health` — Health check
- `POST /api/chat` — Send message

### Authenticated
- `POST /api/documents` — Upload document
- `GET /api/documents` — List documents
- `DELETE /api/documents/:id` — Delete document
- `POST /api/bots` — Create bot
- `GET /api/bots` — List bots
- `PUT /api/bots/:id` — Update bot

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| OPENAI_API_KEY | - | OpenAI API key |
| API_KEY | - | Tenant API key |
| DATABASE_PATH | /app/data | Data directory |
| CHROMA_PATH | - | ChromaDB path |
| PORT | 8000 | Backend port |

## Development

```bash
# Backend only
cd backend
pip install -r requirements.txt
uvicorn src.api:app --reload

# Frontend only
cd frontend
npm install
npm run dev
```

## Deployment

1. Push to `main` → GitHub Action auto-deploys
2. Or manually: `ssh vps && cd /root/chatbot-deploy && ./update.sh`

## License

MIT
