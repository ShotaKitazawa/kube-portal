import React, { createContext, useContext, useEffect, useState } from "react";
import { type UserManager } from "oidc-client-ts";
import { loadOIDCSetup } from "../drivers/auth";

export type AuthUser = {
  access_token: string;
  profile: {
    sub: string;
    name?: string;
    email?: string;
    picture?: string;
  };
};

type AuthContextType = {
  user: AuthUser | null;
  login: () => Promise<void>;
  logout: () => Promise<void>;
};

export const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => {},
  logout: async () => {},
});

const baseUrl = "/api";

async function fetchUserinfo(): Promise<AuthUser | null> {
  try {
    const res = await fetch(`${baseUrl}/userinfo`);
    if (!res.ok) return null;
    const data = await res.json();
    return {
      access_token: "",
      profile: {
        sub: data.sub,
        name: data.name,
        email: data.email,
        picture: data.picture,
      },
    };
  } catch {
    return null;
  }
}

function toAuthUser(u: import("oidc-client-ts").User): AuthUser {
  return {
    access_token: u.access_token,
    profile: {
      sub: u.profile.sub,
      name: u.profile.name ?? undefined,
      email: u.profile.email ?? undefined,
      picture: u.profile.picture ?? undefined,
    },
  };
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [userManager, setUserManager] = useState<UserManager | null>(null);

  useEffect(() => {
    loadOIDCSetup().then((setup) => {
      if (!setup.configured) {
        fetchUserinfo().then(setUser);
        return;
      }

      const mgr = setup.userManager;
      setUserManager(mgr);

      mgr.getUser().then((u) => {
        if (u) setUser(toAuthUser(u));
      });

      const handleUserLoaded = (u: import("oidc-client-ts").User) => setUser(toAuthUser(u));
      const handleUserUnloaded = () => setUser(null);

      mgr.events.addUserLoaded(handleUserLoaded);
      mgr.events.addUserUnloaded(handleUserUnloaded);

      return () => {
        mgr.events.removeUserLoaded(handleUserLoaded);
        mgr.events.removeUserUnloaded(handleUserUnloaded);
      };
    });
  }, []);

  const login = async () => {
    await userManager?.signinRedirect();
  };

  const logout = async () => {
    if (userManager) {
      await userManager.removeUser();
    }
    setUser(null);
  };

  return <AuthContext.Provider value={{ user, login, logout }}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);
