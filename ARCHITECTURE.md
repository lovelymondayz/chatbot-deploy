# Chatbot Deploy — Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Cloudflare Edge                          │
│                 chatbot.arjism.com (HTTPS)                      │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Cloudflare Tunnel (cf-tunnel)                │
│              http://192.168.88.101:8101 (plain HTTP)            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Nginx Reverse Proxy                      │
│                    :8101 → :8000 (backend)                      │
│                    :8102 → :80 (dashboard)                      │
└────────────────────────────┬────────────────────────────────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                              ▼
┌──────────────────────┐        ┌──────────────────────┐
│   Python + FastAPI   │        │  React Dashboard     │
│   :8000 (internal)   │        │  :80 (internal)      │
│                      │        │                      │
│  - LangChain         │        │  - Tailwind CSS      │
│  - ChromaDB          │        │  - Bot Management    │
│  - OpenAI            │        │  - Analytics         │
│  - Vector Search     │        │  - Chat Testing      │
└──────────┬───────────┘        └──────────────────────┘
           │
           ▼
┌──────────────────────┐
│   ChromaDB + Local   │
│   /app/data/         │
│   (Vector Store)     │
└──────────────────────┘
```

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Language | Python | 3.11+ |
| Web Framework | FastAPI | 0.115+ |
| Vector DB | ChromaDB | - |
| LLM Integration | OpenAI / 9Router | - |
| Document Processing | LangChain | - |
| Frontend | React + Vite + TypeScript | Vite 5, React 18 |
| Styling | Tailwind CSS | v3 |
| Deployment | Docker Compose | v3.8 |
| Reverse Proxy | Nginx | - |
| Tunnel | Cloudflare Tunnel | - |

## Key Design Decisions

### 1. RAG-Based Chatbot
- Retrieval-Augmented Generation for accurate responses
- Vector store for document embeddings
- Context-aware conversations

### 2. ChromaDB Vector Store
- Local vector database for fast similarity search
- Persistent storage for embeddings
- Easy migration to cloud vector DBs

### 3. Multi-Tenant Architecture
- API key per tenant
- Isolated document stores per tenant
- Rate limiting per API key

### 4. Dashboard for Management
- Bot configuration UI
- Chat testing interface
- Usage analytics

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check |
| POST | `/api/chat` | Send message |

### Authenticated (API Key)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/documents` | Upload document |
| GET | `/api/documents` | List documents |
| DELETE | `/api/documents/:id` | Delete document |
| POST | `/api/bots` | Create bot |
| GET | `/api/bots` | List bots |
| PUT | `/api/bots/:id` | Update bot |

## Ports

| Service | External | Internal |
|---------|----------|----------|
| Backend | `:8101` | `:8000` |
| Dashboard | `:8102` | `:80` |
