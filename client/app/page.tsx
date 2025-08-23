'use client'

import { useRouter } from 'next/navigation'
import { Plus, MessageSquare } from 'lucide-react'
import amaLogo from '@/src/assets/ama-logo.svg'
import Image from 'next/image'

export default function HomePage() {
  const router = useRouter()

  return (
    <div className="min-h-screen flex items-center justify-center px-4">
      <div className="max-w-[450px] flex flex-col gap-8 text-center">
        <Image src={amaLogo} alt="AMA" className="h-12 mx-auto" />

        <div className="space-y-4">
          <h1 className="text-3xl font-bold text-zinc-100">
            Ask Me Anything
          </h1>
          
          <p className="leading-relaxed text-zinc-300">
            Crie salas públicas de AMA (Ask me anything) e priorize as perguntas mais importantes para a comunidade.
          </p>
        </div>

        <div className="flex flex-col gap-4">
          <button
            onClick={() => router.push('/room/create-room')}
            className="bg-orange-400 text-orange-950 px-6 py-3 gap-3 flex items-center justify-center rounded-xl font-medium text-lg transition-colors hover:bg-orange-500"
          >
            <Plus className="size-5" />
            Criar nova sala
          </button>

          <button
            onClick={() => router.push('/my-rooms')}
            className="bg-zinc-800 text-zinc-300 px-6 py-3 gap-3 flex items-center justify-center rounded-xl font-medium text-lg transition-colors hover:bg-zinc-700 border border-zinc-700"
          >
            <MessageSquare className="size-5" />
            Minhas salas
          </button>
        </div>
      </div>
    </div>
  )
}