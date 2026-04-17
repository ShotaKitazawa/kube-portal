import { UserManager, WebStorageStateStore } from 'oidc-client-ts'

export type OIDCSetup =
  | { configured: false; userManager: null }
  | { configured: true; userManager: UserManager }

export async function loadOIDCSetup(): Promise<OIDCSetup> {
  const baseUrl =
    process.env.NEXT_PUBLIC_BACKEND_URL !== undefined
      ? process.env.NEXT_PUBLIC_BACKEND_URL + '/api'
      : window.location.origin + '/api'

  let config: { enabled: boolean; issuer?: string; client_id?: string }
  try {
    const res = await fetch(`${baseUrl}/oidc-config`)
    config = await res.json()
  } catch {
    return { configured: false, userManager: null }
  }

  if (!config.enabled || !config.issuer || !config.client_id) {
    return { configured: false, userManager: null }
  }

  const userManager = new UserManager({
    authority: config.issuer,
    client_id: config.client_id,
    redirect_uri: `${window.location.origin}/callback`,
    scope: 'openid profile email offline_access',
    userStore: new WebStorageStateStore({ store: window.localStorage }),
    automaticSilentRenew: false,
  })

  return { configured: true, userManager }
}
