/* eslint @next/next/no-img-element: 0 */

import React, { useState } from 'react'
import Image from 'next/image'
import { IconButton, Button, Menu, MenuItem } from '@mui/material'

import { useAuth } from '../../contexts/auth'

export const Login: React.FC = () => {
  const { user, login, logout } = useAuth()
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const open = Boolean(anchorEl)

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }
  const handleClose = () => {
    setAnchorEl(null)
  }

  if (!user) {
    return (
      <Button color="inherit" onClick={login}>
        Login
      </Button>
    )
  }

  return (
    <div>
      <IconButton
        aria-label="account of current user"
        aria-controls="menu-appbar"
        aria-haspopup="true"
        onClick={handleMenu}
        color="inherit"
      >
        <Image
          src={user.profile.picture}
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
        disableScrollLock={true}
        keepMounted
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        open={open}
        onClose={handleClose}
      >
        <MenuItem
          onClick={() => {
            handleClose()
            logout()
          }}
        >
          Logout
        </MenuItem>
      </Menu>
    </div>
  )
}
