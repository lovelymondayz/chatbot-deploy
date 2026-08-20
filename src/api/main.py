
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional
import sqlite3, os, uuid, json, hashlib
from datetime import datetime
from pathlib import Path

app = FastAPI(title="Chatbot Deploy API", version="1.0.0")
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])

DATABASE_PATH = os.getenv("DATABASE_PATH", "/app/data/chatbot.db")
VECTOR_DB_PATH = os.getenv("VECTOR_DB_PATH", "/app/data/vectordb")
Path(DATABASE_PATH).parent.mkdir(parents=True, exist_ok=True)
Path(VECTOR_DB_PATH).mkdir(parents=True, exist_ok=True)

def init_db():
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()
    c.execute("CREATE TABLE IF NOT EXISTS clients (id TEXT PRIMARY KEY, name TEXT, status TEXT DEFAULT 'active', created_at TEXT)")
    c.execute("CREATE TABLE IF NOT EXISTS knowledge_docs (id TEXT PRIMARY KEY, client_id TEXT, content TEXT, doc_type TEXT, created_at TEXT)")
    c.execute("CREATE TABLE IF NOT EXISTS conversations (id TEXT PRIMARY KEY, client_id TEXT, user_message TEXT, bot_response TEXT, created_at TEXT)")
    conn.commit(); conn.close()

def get_db():
    conn = sqlite3.connect(DATABASE_PATH); conn.row_factory = sqlite3.Row; return conn

class ClientReq(BaseModel):
    name: str

class TrainReq(BaseModel):
    client_id: str
    content: str
    doc_type: str = "text"

class ChatReq(BaseModel):
    client_id: str
    message: str

@app.get("/health")
async def health(): return {"status": "healthy"}

@app.get("/")
async def root(): return {"service": "Chatbot Deploy API", "version": "1.0.0"}

@app.post("/clients")
async def add_client(req: ClientReq):
    cid = str(uuid.uuid4())
    conn = get_db()
    c = conn.cursor()
    c.execute("INSERT INTO clients (id,name,status,created_at) VALUES (?,?,?,?)",
        (cid, req.name, "active", datetime.utcnow().isoformat()))
    conn.commit(); conn.close()
    return {"client_id": cid, "message": f"Client {req.name} created", "widget_code": f'<script src="https://chatbot.arjism.com/widget.js" data-client-id="{cid}"></script>'}

@app.get("/clients")
async def list_clients():
    conn = get_db()
    clients = [dict(r) for r in conn.execute("SELECT * FROM clients ORDER BY created_at DESC").fetchall()]
    conn.close(); return {"clients": clients}

@app.post("/clients/{cid}/train")
async def train(cid: str, req: TrainReq):
    conn = get_db()
    c = conn.cursor()
    c.execute("SELECT * FROM clients WHERE id=?", (cid,))
    if not c.fetchone(): conn.close(); raise HTTPException(404, "Client not found")
    did = str(uuid.uuid4())
    c.execute("INSERT INTO knowledge_docs (id,client_id,content,doc_type,created_at) VALUES (?,?,?,?,?)",
        (did, cid, req.content, req.doc_type, datetime.utcnow().isoformat()))
    conn.commit(); conn.close()
    return {"doc_id": did, "message": "Document trained"}

@app.post("/clients/{cid}/chat")
async def chat(cid: str, req: ChatReq):
    conn = get_db()
    c = conn.cursor()
    c.execute("SELECT * FROM clients WHERE id=?", (cid,))
    if not c.fetchone(): conn.close(); raise HTTPException(404, "Client not found")
    c.execute("SELECT content FROM knowledge_docs WHERE client_id=? ORDER BY created_at", (cid,))
    docs = [r["content"] for r in c.fetchall()]
    context = "\n".join(docs[:5]) if docs else "No knowledge base yet."
    response = f"I understand you're asking about: {req.message}\n\nBased on my knowledge: {context[:200]}..."
    conv_id = str(uuid.uuid4())
    c.execute("INSERT INTO conversations (id,client_id,user_message,bot_response,created_at) VALUES (?,?,?,?,?)",
        (conv_id, cid, req.message, response, datetime.utcnow().isoformat()))
    conn.commit(); conn.close()
    return {"response": response}

@app.get("/clients/{cid}/widget")
async def widget(cid: str):
    return {"widget_code": f'<script src="https://chatbot.arjism.com/widget.js" data-client-id="{cid}"></script>'}

@app.get("/clients/{cid}/analytics")
async def analytics(cid: str):
    conn = get_db()
    c = conn.cursor()
    c.execute("SELECT COUNT(*) FROM conversations WHERE client_id=?", (cid,))
    total = c.fetchone()[0]
    c.execute("SELECT COUNT(*) FROM knowledge_docs WHERE client_id=?", (cid,))
    docs = c.fetchone()[0]
    conn.close()
    return {"total_messages": total, "total_docs": docs}

@app.on_event("startup")
async def startup(): init_db()
