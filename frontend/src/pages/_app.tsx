import 'normalize.css'
import 'tailwindcss/tailwind.css'
import type { AppProps } from 'next/app'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'

import { AuthProvider } from '../contexts/auth'
import { GlobalProvider } from '../contexts/global'

const theme = createTheme({
  palette: {
    primary: {
      main: '#3f51b5',
    },
    secondary: {
      main: '#f50057',
    },
  },
})

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline /> {/* reset MUI based style */}
      <AuthProvider>
        <GlobalProvider>
          <Component {...pageProps} />
        </GlobalProvider>
      </AuthProvider>
    </ThemeProvider>
  )
}

export default MyApp
