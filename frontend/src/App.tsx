import React, { useEffect, useState } from "react";
import { CircularProgress } from "@mui/material";
import { fetchUserinfo, loadOIDCSetup, type AuthUser, type OIDCSetup } from "./drivers/auth";
import { Layout } from "./components/Layout";
import { GlobalProvider } from "./contexts/global";
import { LinkGroups } from "./components/LinkGroups";

function App() {
  const [oidcSetup, setOidcSetup] = useState<OIDCSetup | null>(null);
  // undefined = still loading; null = unauthenticated; AuthUser = authenticated
  const [user, setUser] = useState<AuthUser | null | undefined>(undefined);

  useEffect(() => {
    if (window.location.search.includes("code=")) {
      loadOIDCSetup().then((setup) => {
        if (!setup.configured) {
          window.location.replace("/");
          return;
        }
        setup.userManager
          .signinRedirectCallback()
          .then(() => window.location.replace("/"))
          .catch(console.error);
      });
      return;
    }

    loadOIDCSetup().then(async (setup) => {
      setOidcSetup(setup);
      if (setup.configured) {
        setup.userManager.events.addUserLoaded(async () => {
          const u = await setup.userManager.getUser();
          fetchUserinfo(u?.access_token).then(setUser);
        });
        setup.userManager.events.addUserUnloaded(() => setUser(null));
      }
      // GET /api/userinfo is the single source of truth for auth state:
      // 200 = authenticated, non-200 = unauthenticated.
      const token = setup.configured
        ? (await setup.userManager.getUser())?.access_token
        : undefined;
      const u = await fetchUserinfo(token);
      setUser(u);
    });
  }, []);

  if (window.location.search.includes("code=")) {
    return null;
  }

  if (user === undefined) {
    return (
      <div className="flex justify-center items-center h-screen">
        <CircularProgress />
      </div>
    );
  }

  const getToken = async (): Promise<string | null> => {
    if (!oidcSetup?.configured) return null;
    const u = await oidcSetup.userManager.getUser();
    return u?.access_token ?? null;
  };

  const login = async () => {
    await oidcSetup?.userManager.signinRedirect();
  };

  const logout = async () => {
    if (oidcSetup?.configured) await oidcSetup.userManager.removeUser();
    setUser(null);
  };

  return (
    <Layout user={user} onLogin={login} onLogout={logout}>
      <GlobalProvider getToken={getToken}>
        <LinkGroups />
      </GlobalProvider>
    </Layout>
  );
}

export default App;
