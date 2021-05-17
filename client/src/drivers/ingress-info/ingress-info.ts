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

  WithJWT(jwt: string): IngressInfo {
    this.jwt = jwt
    return this
  }

  async List(): Promise<LinkInfo[]> {
    var res
    try {
      res = await axios.get<LinkInfo[]>(this.url, {
        headers: {
          Cookie: this.jwt != "" ? "jwt=" + this.jwt + ";" : ""
        }
      })
    } catch {
      res = { data: null }
    }
    return res.data
  }
}

export default IngressInfo
