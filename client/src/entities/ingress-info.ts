export interface LinksPort {
  List(): LinkInfo[]
}

export type LinkInfo = {
  name: string
  url: string
  icon_url: string
  is_private: boolean
}
