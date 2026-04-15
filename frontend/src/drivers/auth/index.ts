import { UserManager, WebStorageStateStore } from 'oidc-client-ts'

const disableOIDC = process.env.NEXT_PUBLIC_DISABLE_OIDC === 'true'

const userManager: UserManager | null = disableOIDC
  ? null
  : new UserManager({
      authority: process.env.NEXT_PUBLIC_OIDC_AUTHORITY ?? '',
      client_id: process.env.NEXT_PUBLIC_OIDC_CLIENT_ID ?? '',
      redirect_uri:
        typeof window !== 'undefined'
          ? `${window.location.origin}/callback`
          : `${process.env.NEXT_PUBLIC_BACKEND_URL ?? ''}/callback`,
      scope: 'openid profile email offline_access',
      extraQueryParams: {
        audience: process.env.NEXT_PUBLIC_OIDC_AUDIENCE ?? '',
      },
      userStore: new WebStorageStateStore({
        store: typeof window !== 'undefined' ? window.localStorage : undefined,
      }),
      automaticSilentRenew: false,
    })

export { disableOIDC }
export default userManager
