import React from 'react';
import Link from 'next/link';
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import {
  Toolbar,
  AppBar,
  Typography,
  IconButton,
  Button,
} from '@material-ui/core';
import {
  AccountCircle,
} from '@material-ui/icons';


export const Header: React.FC = ({
  children,
}) => {
  const [auth, setAuth] = React.useState(false);

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
            {auth ? (
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
        </Toolbar>
      </AppBar>
    </div >
  )
}
