import { useEffect } from 'react'
import { useRouter } from 'next/router'
import userManager from '../drivers/auth'

export default function Callback() {
  const router = useRouter()

  useEffect(() => {
    userManager
      .signinRedirectCallback()
      .then(() => router.replace('/'))
      .catch(console.error)
  }, [])

  return null
}
