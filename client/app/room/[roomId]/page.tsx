'use client'

import { useParams, useRouter } from "next/navigation"
import { toast } from "sonner"
import { Suspense, useEffect, useState } from "react"
import { Share2 } from "lucide-react"
import amaLogo from '@/src/assets/ama-logo.svg'
import { CreateMessageForm } from "@/src/components/create-message-form"
import { Messages } from "@/src/components/messages"
import Image from "next/image"

export default function RoomPage() {
  const params = useParams()
  const [roomId, setRoomId] = useState<string | null>(null)

  const route = useRouter()

  useEffect(() => {
    console.log('RoomPage - useParams():', params)
    
    if (params.roomId && typeof params.roomId === 'string') {
      console.log('RoomPage - setting roomId:', params.roomId)
      setRoomId(params.roomId)
    } else {
      console.log('RoomPage - invalid roomId:', params.roomId, typeof params.roomId)
      setRoomId(null)
    }
  }, [params.roomId])

  function handleShareRoom() {
    const url = window.location.href.toString()

    if (navigator.share !== undefined && navigator.canShare()) {
      navigator.share({ url })
    } else {
      navigator.clipboard.writeText(url)

      toast.info('O link da sala foi copiado para área de transferência!')
    }
  }

  // Loading state enquanto determina o roomId
  if (roomId === null) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-zinc-400">Carregando sala...</p>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-[640px] flex flex-col gap-6 py-10 px-4">
      <div className="flex items-center gap-3 px-3">
        <Image src={amaLogo} alt="AMA" className="h-5 cursor-pointer" onClick={() => {
          route.push('/')
        }} />

        <span className="text-sm text-zinc-500 truncate">
          Código da sala: <span className="text-zinc-300">{roomId}</span>
        </span>

        <button 
          type="button" 
          onClick={handleShareRoom}
          className="ml-auto bg-zinc-800 text-zinc-300 px-3 py-1.5 gap-1.5 flex items-center rounded-lg font-medium text-sm transition-colors hover:bg-zinc-700"
        >
          Compartilhar
          <Share2 className="size-4" />
        </button>
      </div>

      <div className="h-px w-full bg-zinc-900" />

      <CreateMessageForm roomId={roomId} />

      <Suspense fallback={<p>Carregando...</p>}>
        <Messages roomId={roomId} />
      </Suspense>
    </div>
  )
}
