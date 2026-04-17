import { UserManager, WebStorageStateStore } from 'oidc-client-ts'

const disableOIDC = process.env.NEXT_PUBLIC_DISABLE_OIDC === 'true'

const userManager: UserManager | null =
  disableOIDC || typeof window === 'undefined'
    ? null
    : new UserManager({
        authority: process.env.NEXT_PUBLIC_OIDC_AUTHORITY ?? '',
        client_id: process.env.NEXT_PUBLIC_OIDC_CLIENT_ID ?? '',
        redirect_uri: `${window.location.origin}/callback`,
        scope: 'openid profile email offline_access',
        extraQueryParams: {
          audience: process.env.NEXT_PUBLIC_OIDC_AUDIENCE ?? '',
        },
        userStore: new WebStorageStateStore({ store: window.localStorage }),
        automaticSilentRenew: false,
      })

export { disableOIDC }
export default userManager
