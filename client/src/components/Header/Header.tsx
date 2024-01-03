import React from 'react'
import Link from 'next/link'
import { Toolbar, AppBar, Typography, IconButton } from '@mui/material'

import { Login } from '../Login'

export const Header: React.FC = ({ children }) => {
  return (
    <div className="flex">
      <AppBar position="static">
        <Toolbar>
          <div className="mr-2">
            <IconButton edge="start" color="inherit">
              <Link href="/">
                <span className="iconify" data-icon="mdi-kubernetes" />
              </Link>
            </IconButton>
          </div>
          <div className="flex-grow">
            <Typography variant="h6">kube-portal</Typography>
          </div>
          <div className="ml-auto">
            <Login />
          </div>
        </Toolbar>
      </AppBar>
    </div>
  )
}
