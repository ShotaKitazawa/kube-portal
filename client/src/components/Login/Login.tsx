import React from 'react';
import Link from 'next/link';
import { parseCookies } from 'nookies'
import {
  IconButton,
  Button,
} from '@material-ui/core';
import {
  AccountCircle,
} from '@material-ui/icons';

type User = {
  name: string
  avatar_url: string
}

export const Login: React.FC = ({
  children,
}) => {
  /* TODO
  const cookies = parseCookies()
  if (cookies.jwt !== 'undefined') {
    var user = {
      name: cookies.jwt,
      avatar_url: cookies.jwt,
    }
  }
  */

  return (
    <div>
      {false ? ( // TODO
        <div>
          <IconButton
            aria-label="account of current user"
            aria-controls="menu-appbar"
            aria-haspopup="true"
            color="inherit"
          >
            <AccountCircle />
          </IconButton>
        </div>
      ) : (
        <a href="/auth/login">
          <Button color="inherit">Login</Button>
        </a>
      )
      }
    </div>
  )
}
