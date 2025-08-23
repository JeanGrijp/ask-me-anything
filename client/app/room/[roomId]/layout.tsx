import { Metadata } from 'next'
import { notFound } from 'next/navigation'

interface RoomLayoutProps {
  children: React.ReactNode
  params: {
    roomId: string
  }
}

// Função para buscar informações da sala para metadata
async function getRoomInfo(roomId: string) {
  try {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL
    
    if (!apiUrl) {
      return null
    }

    const response = await fetch(`${apiUrl}/rooms/${roomId}`, {
      cache: 'no-store' // Para garantir dados atualizados
    })

    if (!response.ok) {
      return null
    }

    return await response.json()
  } catch {
    return null
  }
}

export async function generateMetadata({ params }: RoomLayoutProps): Promise<Metadata> {
  const roomInfo = await getRoomInfo(params.roomId)
  
  if (!roomInfo) {
    return {
      title: 'Sala não encontrada - Ask Me Anything',
      description: 'Esta sala de AMA não foi encontrada ou não existe mais.',
    }
  }

  const roomTitle = roomInfo.theme || 'Sala de AMA'
  
  return {
    title: `${roomTitle} - Ask Me Anything`,
    description: `Participe da sala "${roomTitle}" e faça suas perguntas para a comunidade. AMA (Ask Me Anything).`,
    openGraph: {
      title: `${roomTitle} - Ask Me Anything`,
      description: `Participe da sala "${roomTitle}" e faça suas perguntas`,
      type: 'website',
      url: `${process.env.NEXT_PUBLIC_BASE_URL}/room/${params.roomId}`,
    },
    twitter: {
      card: 'summary',
      title: `${roomTitle} - Ask Me Anything`,
      description: `Participe da sala "${roomTitle}" e faça suas perguntas`,
    },
    robots: {
      index: true,
      follow: true,
    },
  }
}

export default function RoomLayout({ children, params }: RoomLayoutProps) {
  return (
    <div className="min-h-screen bg-zinc-950">
      {/* Container principal com padding responsivo */}
      <main className="container mx-auto px-4 sm:px-6 lg:px-8">
        {/* Wrapper para garantir espaçamento adequado */}
        <div className="min-h-screen flex flex-col">
          {/* Header da sala - pode ser expandido futuramente */}
          <div className="py-4">
            {children}
          </div>
        </div>
      </main>
    </div>
  )
}