import React, { createContext, useContext, useEffect, useState } from 'react'
import userManager, { disableOIDC } from '../drivers/auth'

export type AuthUser = {
  access_token: string
  profile: {
    sub: string
    name?: string
    email?: string
    picture?: string
  }
}

type AuthContextType = {
  user: AuthUser | null
  login: () => Promise<void>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => {},
  logout: async () => {},
})

const baseUrl =
  process.env.NEXT_PUBLIC_BACKEND_URL !== undefined
    ? process.env.NEXT_PUBLIC_BACKEND_URL + '/api'
    : typeof window !== 'undefined'
      ? window.location.origin + '/api'
      : '/api'

async function fetchMockUser(): Promise<AuthUser | null> {
  try {
    const res = await fetch(`${baseUrl}/userinfo`)
    if (!res.ok) return null
    const data = await res.json()
    return {
      access_token: '',
      profile: {
        sub: data.sub,
        name: data.name,
        email: data.email,
        picture: data.picture,
      },
    }
  } catch {
    return null
  }
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<AuthUser | null>(null)

  useEffect(() => {
    if (disableOIDC) {
      fetchMockUser().then(setUser)
      return
    }

    if (!userManager) return

    userManager.getUser().then((u) => {
      if (u) {
        setUser({
          access_token: u.access_token,
          profile: {
            sub: u.profile.sub,
            name: u.profile.name ?? undefined,
            email: u.profile.email ?? undefined,
            picture: u.profile.picture ?? undefined,
          },
        })
      }
    })

    const handleUserLoaded = (u: import('oidc-client-ts').User) => {
      setUser({
        access_token: u.access_token,
        profile: {
          sub: u.profile.sub,
          name: u.profile.name ?? undefined,
          email: u.profile.email ?? undefined,
          picture: u.profile.picture ?? undefined,
        },
      })
    }
    const handleUserUnloaded = () => setUser(null)

    userManager.events.addUserLoaded(handleUserLoaded)
    userManager.events.addUserUnloaded(handleUserUnloaded)

    return () => {
      userManager!.events.removeUserLoaded(handleUserLoaded)
      userManager!.events.removeUserUnloaded(handleUserUnloaded)
    }
  }, [])

  const login = async () => {
    if (disableOIDC) return
    await userManager?.signinRedirect()
  }

  const logout = async () => {
    if (disableOIDC) {
      setUser(null)
      return
    }
    await userManager?.removeUser()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)
