import { UserManager, WebStorageStateStore } from "oidc-client-ts";

export type OIDCSetup =
  | { configured: false; userManager: null }
  | { configured: true; userManager: UserManager };

export type AuthUser = {
  sub: string;
  name?: string;
  email?: string;
  picture?: string;
};

export async function loadOIDCSetup(): Promise<OIDCSetup> {
  const baseUrl = "/api";

  let config: { enabled: boolean; issuer?: string; client_id?: string; audience?: string };
  try {
    const res = await fetch(`${baseUrl}/oidc-config`);
    config = await res.json();
  } catch {
    return { configured: false, userManager: null };
  }

  if (!config.enabled || !config.issuer || !config.client_id) {
    return { configured: false, userManager: null };
  }

  const userManager = new UserManager({
    authority: config.issuer,
    client_id: config.client_id,
    redirect_uri: `${window.location.origin}/callback`,
    scope: "openid profile email offline_access",
    extraQueryParams: config.audience ? { audience: config.audience } : {},
    userStore: new WebStorageStateStore({ store: window.localStorage }),
    automaticSilentRenew: false,
  });

  return { configured: true, userManager };
}

export async function fetchUserinfo(token?: string): Promise<AuthUser | null> {
  try {
    const headers: HeadersInit = token ? { Authorization: `Bearer ${token}` } : {};
    const res = await fetch("/api/userinfo", { headers });
    if (!res.ok) return null;
    const data = await res.json();
    return {
      sub: data.sub,
      name: data.name,
      email: data.email,
      picture: data.picture,
    };
  } catch {
    return null;
  }
}
