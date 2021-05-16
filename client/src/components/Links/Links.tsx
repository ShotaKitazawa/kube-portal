import React, { useEffect, useState } from 'react'
import {
  CircularProgress,
  Card,
  CardContent,
} from '@material-ui/core';

import IngressInfo from '../../drivers/ingress-info/ingress-info'
import { LinksPort, LinkInfo } from '../../entities/ingress-info'


export const Links: React.FC = ({
  children,
}) => {

  // get links from API
  var linksList = []
  if (typeof window !== 'undefined') {
    linksList = new IngressInfo(window.location.origin).List()
  }

  return (
    <section id="links" className="mt-6">
      <div className="mx-auto flex flex-wrap justify-center w-1/2">
        {
          linksList != null ?
            linksList.map(({ name, url, icon_url }) => (
              <div className="mx-5 my-5 w-5/6 md:w-4/6 lg:w-3/6 xl:w-2/6 transition duration-300 ease-in-out transform hover:-translate-y-1 hover:scale-110">
                <a href={url} className="hover:no-underline">
                  <Card className="flex hover:bg-indigo-200">
                    <div className="flex-none w-16 h-16">
                      {icon_url != "" ?
                        <img src={icon_url} />
                        :
                        <img src="https://raw.githubusercontent.com/FortAwesome/Font-Awesome/master/svgs/solid/external-link-alt.svg" />
                      }
                    </div>
                    <div className="flex-grow w-16">
                      <CardContent >
                        <p className="text-center text-xl font-bold">{name}</p>
                      </CardContent>
                    </div>
                  </Card>
                </a>
              </div>
            ))
            :
            <CircularProgress />
        }
      </div>
    </section>
  )
}
