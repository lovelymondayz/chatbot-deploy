import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface Client {
  id: string
  name: string
  status: string
  created_at: string
}

export interface ClientsResponse {
  clients: Client[]
}

export interface CreateClientResponse {
  client_id: string
  message: string
  widget_code: string
}

export interface TrainResponse {
  doc_id: string
  message: string
}

export interface ChatResponse {
  response: string
}

export async function getClients(): Promise<ClientsResponse> {
  const res = await api.get('/clients')
  return res.data
}

export async function createClient(name: string): Promise<CreateClientResponse> {
  const res = await api.post('/clients', { name })
  return res.data
}

export async function trainClient(clientId: string, content: string, docType?: string): Promise<TrainResponse> {
  const res = await api.post(`/clients/${clientId}/train`, { client_id: clientId, content, doc_type: docType || 'text' })
  return res.data
}

export async function chatClient(clientId: string, message: string): Promise<ChatResponse> {
  const res = await api.post(`/clients/${clientId}/chat`, { client_id: clientId, message })
  return res.data
}

export default api