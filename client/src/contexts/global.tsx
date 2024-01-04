import React, { useState, useEffect } from 'react'

import Client, { LinkInfo } from '../drivers/ingress-info'

type contexts = {
  linkInfo?: LinkInfo[]
}

export const GlobalContext = React.createContext({} as contexts)

type Props = {
  children: React.ReactNode
}

export const GlobalProvider: React.FC<Props> = (props) => {
  const [linkInfo, setLinkInfo] = useState<LinkInfo[]>()

  const client = new Client()
  const listLinkInfo = async () => {
    try {
      const resp = await client.List()
      setLinkInfo(resp)
    } catch (error) {}
  }

  useEffect(() => {
    listLinkInfo()
    setInterval(() => {
      listLinkInfo()
    }, 10000)
  }, [])

  return (
    <GlobalContext.Provider
      value={{
        linkInfo,
      }}
    >
      {props.children}
    </GlobalContext.Provider>
  )
}
