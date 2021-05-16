import useSWR from 'swr'

import { LinksPort, LinkInfo } from '../../entities/ingress-info'

class IngressInfo implements LinksPort {
  private url: string

  constructor() {
    this.url = "http://localhost:8080/api/list" // TODO
  }

  List(): LinkInfo[] {
    const fetcher = async (url: string): Promise<LinkInfo[] | null> => {
      const response = await fetch(url)
      const links = await response.json()
      console.log(links)
      return links
    }

    const { data } = useSWR(this.url, fetcher)

    if (data === undefined) {
      return null
    }

    return data
  }
}

export default IngressInfo
