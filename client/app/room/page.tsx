'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

export default function RoomsPage() {
  const router = useRouter()

  useEffect(() => {
    // Redireciona para a página de criar sala ou página inicial
    router.push('/room/create-room')
  }, [router])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <p className="text-zinc-400">Redirecionando...</p>
    </div>
  )
}