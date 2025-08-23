import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Minhas Salas | AMA',
  description: 'Gerencie suas salas de perguntas e respostas. Visualize estatísticas, modere mensagens e acompanhe a participação dos usuários.',
  openGraph: {
    title: 'Minhas Salas | AMA',
    description: 'Gerencie suas salas de perguntas e respostas. Visualize estatísticas, modere mensagens e acompanhe a participação dos usuários.',
  },
}

interface MyRoomsLayoutProps {
  children: React.ReactNode
}

export default function MyRoomsLayout({ children }: MyRoomsLayoutProps) {
  return (
    <div className="min-h-screen bg-zinc-950">
      {/* Container principal para listagem de salas */}
      <main className="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
        {/* Container com largura máxima para melhor legibilidade */}
        <div className="max-w-4xl mx-auto">
          {children}
        </div>
      </main>
    </div>
  )
}
