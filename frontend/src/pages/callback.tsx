import { useEffect } from 'react'
import { useRouter } from 'next/router'
import userManager, { disableOIDC } from '../drivers/auth'

export default function Callback() {
  const router = useRouter()

  useEffect(() => {
    if (disableOIDC) {
      router.replace('/')
      return
    }
    userManager!
      .signinRedirectCallback()
      .then(() => router.replace('/'))
      .catch(console.error)
  }, [])

  return null
}
