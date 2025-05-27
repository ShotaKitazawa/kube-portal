/* eslint @next/next/no-img-element: 0 */

import React, { ReactNode } from 'react'
import { Card, CardContent } from '@mui/material'
import Image from 'next/image'

import IngressInfo, { LinkInfo } from '../../drivers/ingress-info'

interface LinksProps {
  linksList: LinkInfo[]
  children?: ReactNode
}

export const Links: React.FC<LinksProps> = ({ linksList }) => {
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
                      <Image
                        alt="item-icon"
                        src={icon_url}
                        width={64}
                        height={64}
                      />
                    ) : (
                      <Image
                        alt="item-icon"
                        src="/external-link-alt.svg"
                        width={64}
                        height={64}
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
          <React.Fragment></React.Fragment>
        )}
      </div>
    </section>
  )
}
