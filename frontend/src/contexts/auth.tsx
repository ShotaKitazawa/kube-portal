import React, { createContext, useContext, useEffect, useState } from 'react'
import { User } from 'oidc-client-ts'
import userManager from '../drivers/auth'

type AuthContextType = {
  user: User | null
  login: () => Promise<void>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => {},
  logout: async () => {},
})

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    userManager.getUser().then(setUser)

    const handleUserLoaded = (u: User) => setUser(u)
    const handleUserUnloaded = () => setUser(null)

    userManager.events.addUserLoaded(handleUserLoaded)
    userManager.events.addUserUnloaded(handleUserUnloaded)

    return () => {
      userManager.events.removeUserLoaded(handleUserLoaded)
      userManager.events.removeUserUnloaded(handleUserUnloaded)
    }
  }, [])

  const login = () => userManager.signinRedirect()
  const logout = async () => {
    await userManager.removeUser()
  }

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => useContext(AuthContext)
