import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Criar Nova Sala - Ask Me Anything',
  description: 'Crie uma nova sala de AMA (Ask Me Anything) e comece a receber perguntas da comunidade.',
  openGraph: {
    title: 'Criar Nova Sala - Ask Me Anything',
    description: 'Crie uma nova sala de AMA e comece a receber perguntas da comunidade',
    type: 'website',
    url: `${process.env.NEXT_PUBLIC_BASE_URL}/room/create-room`,
  },
  twitter: {
    card: 'summary',
    title: 'Criar Nova Sala - Ask Me Anything',
    description: 'Crie uma nova sala de AMA e comece a receber perguntas da comunidade',
  },
  robots: {
    index: true,
    follow: true,
  },
}

interface CreateRoomLayoutProps {
  children: React.ReactNode
}

export default function CreateRoomLayout({ children }: CreateRoomLayoutProps) {
  return (
    <div className="min-h-screen bg-zinc-950">
      {/* Container centralizado para o formulário */}
      <main className="container mx-auto px-4 sm:px-6 lg:px-8">
        {/* Wrapper flexbox para centralização vertical */}
        <div className="min-h-screen flex items-center justify-center">
          {/* Container do formulário com largura máxima */}
          <div className="w-full max-w-[450px]">
            {children}
          </div>
        </div>
      </main>
    </div>
  )
}
