import React, { ReactNode } from "react";
import { Header } from "../Header";
import type { AuthUser } from "../../drivers/auth";

interface LayoutProps {
  user: AuthUser | null;
  onLogin: () => void;
  onLogout: () => void;
  children: ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ user, onLogin, onLogout, children }) => {
  return (
    <>
      <div className="relative z-20">
        <Header user={user} onLogin={onLogin} onLogout={onLogout} />
      </div>
      <main className="relative z-10 mt-20">{children}</main>
    </>
  );
};
