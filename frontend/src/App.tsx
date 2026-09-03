import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Clients from './pages/Clients'
import ClientDetail from './pages/ClientDetail'
import NewClient from './pages/NewClient'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout><Clients /></Layout>} />
        <Route path="/clients" element={<Layout><Clients /></Layout>} />
        <Route path="/clients/new" element={<Layout><NewClient /></Layout>} />
        <Route path="/clients/:id" element={<Layout><ClientDetail /></Layout>} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}