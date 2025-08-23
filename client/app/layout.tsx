import type { Metadata } from 'next'
import './globals.css'
import { ReactQueryProvider } from './providers'
import { Toaster } from '@/src/components/ui/sonner'

export const metadata: Metadata = {
  title: 'Ask Me Anything - Next.js',
  description: 'Uma aplicação de perguntas e respostas em tempo real',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="pt-BR" className="dark">
      <body>
        <ReactQueryProvider>
          {children}
        </ReactQueryProvider>
        <Toaster />
      </body>
    </html>
  )
}
