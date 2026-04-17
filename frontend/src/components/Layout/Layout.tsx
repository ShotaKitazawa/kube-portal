import React, { ReactNode } from 'react'

import { Header } from '../Header'

interface LayoutProps {
  children: ReactNode
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <>
      <div className="relative z-20">
        <Header />
      </div>
      <main className="relative z-10 mt-20">{children}</main>
    </>
  )
}
