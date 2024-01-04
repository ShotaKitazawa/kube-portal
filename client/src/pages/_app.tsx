import 'normalize.css'
import 'tailwindcss/tailwind.css'

import { GlobalProvider } from '../contexts/global'

function MyApp({ Component, pageProps }) {
  return (
    <GlobalProvider>
      <Component {...pageProps} />
    </GlobalProvider>
  )
}

export default MyApp
