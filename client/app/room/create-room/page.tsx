'use client'

import { ArrowRight } from 'lucide-react'
import { toast } from 'sonner'
import { useRouter } from 'next/navigation'

import amaLogo from '@/src/assets/ama-logo.svg'
import { createRoom } from '@/src/http/create-room'
import { setHostToken } from '@/src/utils/host-token'
import Image from 'next/image'

export default function CreateRoom() {
  const router = useRouter()

  async function handleCreateRoom(data: FormData) {
    const theme = data.get('theme')?.toString()

    if (!theme) {
      return
    }

    try {
      const { roomId, hostToken } = await createRoom({ theme })

      console.log('CreateRoom - roomId received:', roomId, typeof roomId)
      console.log('CreateRoom - hostToken received:', hostToken)
      console.log('CreateRoom - navigating to:', `/room/${roomId}`)

      // Armazenar o host token no localStorage
      setHostToken(roomId, hostToken)

      router.push(`/room/${roomId}`)
    } catch {
      toast.error('Falha ao criar sala!')
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <Image src={amaLogo} alt="AMA" className="h-10 mx-auto" />

      <p className="leading-relaxed text-zinc-300 text-center">
        Crie uma sala pública de AMA (Ask me anything) e priorize as perguntas mais importantes para a comunidade.
      </p>

      <form 
        action={handleCreateRoom}
        className="flex items-center gap-2 bg-zinc-900 p-2 rounded-xl border border-zinc-800 ring-orange-400 ring-offset-2 ring-offset-zinc-950 focus-within:ring-1"
      >
        <input 
          type="text" 
          name="theme"
          placeholder="Nome da sala"
          autoComplete="off"
          className="flex-1 text-sm bg-transparent mx-2 outline-none text-zinc-100 placeholder:text-zinc-500"
          required
        />

        <button 
          type="submit" 
          className="bg-orange-400 text-orange-950 px-3 py-1.5 gap-1.5 flex items-center rounded-lg font-medium text-sm transition-colors hover:bg-orange-500"
        >
          Criar sala
          <ArrowRight className="size-4" />
        </button>
      </form>
    </div>
  )
}