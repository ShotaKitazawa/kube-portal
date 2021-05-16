import useSWR from 'swr'

import { LinksPort, LinkInfo } from '../../entities/ingress-info'

class IngressInfo implements LinksPort {
  private url: string

  constructor(origin: string) {
    this.url = origin + "/api/list"
  }

  List(): LinkInfo[] {
    const fetcher = async (url: string): Promise<LinkInfo[] | null> => {
      const response = await fetch(url)
      return await response.json()
    }
    console.log(this.url)

    //return null
    const { data } = useSWR(this.url, fetcher)

    if (data === undefined) {
      return null
    }

    return data
  }
}

export default IngressInfo
