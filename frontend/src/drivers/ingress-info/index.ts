import createClient from 'openapi-fetch'
import type { paths } from '../../gen/api'

export type LinkInfo = paths['/list']['get']['responses'][200]['content']['application/json'][number]

class Client {
  private client: ReturnType<typeof createClient<paths>>

  constructor() {
    const baseUrl =
      process.env.NEXT_PUBLIC_BACKEND_URL !== undefined
        ? process.env.NEXT_PUBLIC_BACKEND_URL + '/api'
        : typeof window !== 'undefined'
          ? window.location.origin + '/api'
          : '/api'

    this.client = createClient<paths>({ baseUrl })
  }

  async List(accessToken?: string): Promise<LinkInfo[]> {
    const headers: Record<string, string> = {}
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`
    }
    const { data } = await this.client.GET('/list', { headers }).catch(() => ({ data: null }))
    return data ?? []
  }
}

export default Client
