/* eslint @next/next/no-img-element: 0 */

import { useState, useEffect } from 'react'
import { InferGetStaticPropsType, NextPage } from 'next'
import { CircularProgress, Card, CardContent } from '@mui/material'

import IngressInfo, { LinkInfo } from '../../drivers/ingress-info'

type Props = InferGetStaticPropsType<typeof getStaticProps>

export const getStaticProps = async () => {
  return { props: {} }
}

const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL ?? window.location.origin

export const Links: React.FC<Props> = (props) => {
  // get links from API
  const [linksList, setLinksList] = useState<LinkInfo[]>(null)
  const list = async () => {
    if (typeof window !== undefined) {
      setLinksList(await new IngressInfo(backendUrl).List())
    }
  }
  useEffect(() => {
    list()
    setInterval(() => {
      list()
    }, 10000)
  }, [])

  return (
    <section id="links" className="mt-6">
      <div className="mx-auto flex flex-wrap justify-center">
        {linksList != null ? (
          linksList.map(({ name, url, icon_url }) => (
            <div
              key={name}
              className="mx-5 my-5 w-5/6 md:w-4/6 lg:w-3/6 xl:w-2/6 transition duration-300 ease-in-out transform hover:-translate-y-1 hover:scale-110"
            >
              <a href={url} className="hover:no-underline">
                <Card className="flex hover:bg-indigo-200">
                  <div className="flex-none w-16 h-16">
                    {icon_url != '' ? (
                      <img alt="icon" src={icon_url} />
                    ) : (
                      <img
                        alt="icon"
                        src="https://raw.githubusercontent.com/FortAwesome/Font-Awesome/master/svgs/solid/external-link-alt.svg"
                      />
                    )}
                  </div>
                  <div className="flex-grow w-16">
                    <CardContent>
                      <p className="text-center text-xl font-bold">{name}</p>
                    </CardContent>
                  </div>
                </Card>
              </a>
            </div>
          ))
        ) : (
          <CircularProgress />
        )}
      </div>
    </section>
  )
}
