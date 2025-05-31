export type JwtPayload = {
  name: string
  picture: string
  iss: string
  exp: number
  iat: number
}

class Jwt {
  private payload: JwtPayload

  constructor(JwtBase64ed: string) {
    const base64Url = JwtBase64ed.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    this.payload = JSON.parse(decodeURIComponent(escape(window.atob(base64))))
  }

  GetUsername(): string {
    return this.payload.name
  }
  GetAvatarUrl(): string {
    return this.payload.picture
  }
}

export default Jwt
