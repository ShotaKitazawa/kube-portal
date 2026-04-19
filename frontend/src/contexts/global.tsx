import React, { useState, useEffect, useContext } from "react";

import Client, { LinkInfo } from "../drivers/ingress-info";
import { AuthContext } from "./auth";

type contexts = {
  linkInfo?: LinkInfo[];
};

export const GlobalContext = React.createContext({} as contexts);

interface GlobalProviderProps {
  children: React.ReactNode;
}

export const GlobalProvider: React.FC<GlobalProviderProps> = (props) => {
  const [linkInfo, setLinkInfo] = useState<LinkInfo[]>();
  const { user } = useContext(AuthContext);

  const client = new Client();
  const listLinkInfo = async () => {
    try {
      const resp = await client.List(user?.access_token);
      setLinkInfo(resp);
    } catch {}
  };

  useEffect(() => {
    listLinkInfo();
    const intervalId = setInterval(() => {
      listLinkInfo();
    }, 10000);
    return () => {
      clearInterval(intervalId);
    };
  }, [user]);

  return (
    <GlobalContext.Provider
      value={{
        linkInfo,
      }}
    >
      {props.children}
    </GlobalContext.Provider>
  );
};
