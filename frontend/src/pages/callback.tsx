import { useEffect } from 'react'
import { useRouter } from 'next/router'
import { loadOIDCSetup } from '../drivers/auth'

export default function Callback() {
  const router = useRouter()

  useEffect(() => {
    loadOIDCSetup().then((setup) => {
      if (!setup.configured) {
        router.replace('/')
        return
      }
      setup.userManager
        .signinRedirectCallback()
        .then(() => router.replace('/'))
        .catch(console.error)
    })
  }, [])

  return null
}
