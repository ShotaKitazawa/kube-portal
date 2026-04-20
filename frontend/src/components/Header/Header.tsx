import React, { ReactNode } from "react";
import { Toolbar, AppBar, Typography, IconButton } from "@mui/material";
import { Login } from "../Login";
import type { AuthUser } from "../../drivers/auth";

const repositoryLink = "https://github.com/ShotaKitazawa/kube-portal";

interface HeaderProps {
  user: AuthUser | null;
  onLogin: () => void;
  onLogout: () => void;
  children?: ReactNode;
}

export const Header: React.FC<HeaderProps> = ({ user, onLogin, onLogout }) => {
  return (
    <div className="fixed w-full top-0 left-0 flex">
      <AppBar position="static">
        <Toolbar>
          <div className="mr-2">
            <a href={repositoryLink}>
              <IconButton edge="start" color="inherit">
                <img alt="icon" src="/kube-portal.png" width={32} height={32} />
              </IconButton>
            </a>
          </div>
          <div className="flex-grow">
            <Typography variant="h6">kube-portal</Typography>
          </div>
          <div className="ml-auto">
            <Login user={user} onLogin={onLogin} onLogout={onLogout} />
          </div>
        </Toolbar>
      </AppBar>
    </div>
  );
};
