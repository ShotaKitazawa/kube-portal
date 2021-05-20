import axios from 'axios';

export type LinkInfo = {
  name: string
  url: string
  icon_url: string
}


class IngressInfo {
  private url: string
  private jwt: string

  constructor(origin: string) {
    this.url = origin + "/api/list"
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

export default IngressInfo
