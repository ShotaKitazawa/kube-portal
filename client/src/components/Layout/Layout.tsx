import Head from 'next/head'

import { Header } from '../Header'

const name = "Shota Kitazawa (kanata)"

export const Layout: React.FC = ({
  children,
}) => {
  return (
    <>
      <Head>
        <title>{name}</title>
        <link rel="shortcut icon" href="/favicon.ico" />
        <meta charSet="utf-8" />
        <meta name="viewport" content="initial-scale=1.0, width=device-width" />
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-fork-ribbon-css/0.2.3/gh-fork-ribbon.min.css" />
      </Head>
      <div className="relative z-20">
        <Header />
      </div>
      <main className="relative z-10">
        {children}
      </main>
    </>
  )
}
