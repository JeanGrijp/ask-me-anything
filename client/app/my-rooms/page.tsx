'use client'

import { UserRoomsList } from '@/src/components/user-rooms-list'
import { LogoutButton } from '@/src/components/logout-button'
import { Plus } from 'lucide-react'
import { useRouter } from 'next/navigation'

export default function MyRoomsPage() {
  const router = useRouter()

  return (
    <>
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold text-zinc-100">Minhas Salas</h1>
        
        <div className="flex items-center gap-4">
          <button
            onClick={() => router.push('/room/create-room')}
            className="flex items-center gap-2 bg-orange-400 text-orange-950 px-4 py-2 rounded-lg font-medium text-sm transition-colors hover:bg-orange-500"
          >
            <Plus className="size-4" />
            Nova Sala
          </button>
          
          <LogoutButton />
        </div>
      </div>

      {/* Lista de salas */}
      <UserRoomsList />
    </>
  )
}
