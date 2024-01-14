import React from 'react'
import Link from 'next/link'
import { Toolbar, AppBar, Typography, IconButton } from '@mui/material'

import { Login } from '../Login'

const repositoryLink = 'https://github.com/ShotaKitazawa/kube-portal'

export const Header: React.FC = ({ children }) => {
  return (
    <div className="flex">
      <AppBar position="static">
        <Toolbar>
          <div className="mr-2">
            <a href={repositoryLink}>
              <IconButton edge="start" color="inherit">
                <span className="iconify" data-icon="mdi-kubernetes" />
              </IconButton>
            </a>
          </div>
          <div className="flex-grow">
            <a href={repositoryLink}>
              <Typography variant="h6">kube-portal</Typography>
            </a>
          </div>
          <div className="ml-auto">
            <Login />
          </div>
        </Toolbar>
      </AppBar>
    </div>
  )
}
