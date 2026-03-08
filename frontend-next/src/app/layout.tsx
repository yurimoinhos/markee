import type { Metadata } from 'next';
import type { ReactNode } from 'react';
import '@/lib/server/local-storage-shim';
import './globals.css';

export const metadata: Metadata = {
  title: 'AggiPay Console',
  description: 'Console financeiro com contratos, cobrança, pagamentos e automações.',
  openGraph: {
    title: 'AggiPay Console',
    description: 'Gestão financeira B2B com operações em tempo real.',
  },
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="pt-BR">
      <body>{children}</body>
    </html>
  );
}
