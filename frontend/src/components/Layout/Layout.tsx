import React, { ReactNode } from 'react'
import Head from 'next/head'

import { Header } from '../Header'

const name = 'kube-portal'

interface LayoutProps {
  children: ReactNode
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <>
      <Head>
        <title>{name}</title>
        <link rel="shortcut icon" href="/favicon.ico" />
        <meta charSet="utf-8" />
        <meta name="viewport" content="initial-scale=1.0, width=device-width" />
      </Head>
      <div className="relative z-20">
        <Header />
      </div>
      <main className="relative z-10 mt-20">{children}</main>
    </>
  )
}
