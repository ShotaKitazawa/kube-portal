import React, { useState, useEffect } from 'react'

import Client, { LinkInfo } from '../drivers/ingress-info'

type contexts = {
  linkInfo?: LinkInfo[]
}

export const GlobalContext = React.createContext({} as contexts)

interface GlobalProviderProps {
  children: React.ReactNode
}

export const GlobalProvider: React.FC<GlobalProviderProps> = (props) => {
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
    const intervalId = setInterval(() => {
      listLinkInfo()
    }, 10000)
    return () => {
      clearInterval(intervalId)
    }
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
