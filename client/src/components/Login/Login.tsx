/* eslint @next/next/no-img-element: 0 */

import React, { useState, useEffect } from 'react'
import { parseCookies } from 'nookies'
import { IconButton, Button, Menu, MenuItem } from '@mui/material'

import Jwt from '../../drivers/jwt'

type User = {
  user: string
  avatarUrl: string
}

export const Login: React.FC = ({ children }) => {
  const [user, setUser] = useState<User>(null)
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)
  const open = Boolean(anchorEl)
  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }
  const handleClose = () => {
    setAnchorEl(null)
  }

  useEffect(() => {
    const cookies = parseCookies()
    if (cookies.jwt !== undefined) {
      const jwt: Jwt = new Jwt(cookies.jwt)
      setUser({
        user: jwt.GetUsername(),
        avatarUrl: jwt.GetAvatarUrl(),
      })
    }
  }, [])

  return (
    <div>
      {user != null ? (
        <div>
          <IconButton
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            onClick={handleMenu}
            color="inherit"
          >
            <img
              src={user.avatarUrl}
              alt="avatar"
              className="rounded-full"
              width={40}
              height={40}
            />
          </IconButton>
          <Menu
            id="menu-appbar"
            anchorEl={anchorEl}
            anchorOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            keepMounted
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            open={open}
            onClose={handleClose}
          >
            <MenuItem onClick={handleClose}>
              <a href="/auth/logout">Logout</a>
            </MenuItem>
          </Menu>
        </div>
      ) : (
        <a href="/auth/login">
          <Button color="inherit">Login</Button>
        </a>
      )}
    </div>
  )
}
