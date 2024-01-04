import axios from 'axios'

export type LinkInfo = {
  name: string
  url: string
  icon_url: string
  tags: string[]
}

class Client {
  private url: string

  constructor() {
    let origin: string
    if (process.env.NEXT_PUBLIC_BACKEND_URL !== undefined) {
      origin = process.env.NEXT_PUBLIC_BACKEND_URL
    } else if (typeof window !== 'undefined') {
      origin = window.location.origin
    }
    this.url = origin + '/api/list'
  }

  async List(): Promise<LinkInfo[]> {
    let res
    try {
      res = await axios.get<LinkInfo[]>(this.url, { withCredentials: true })
    } catch {
      res = { data: null }
    }
    return res.data
  }
}

export default Client
