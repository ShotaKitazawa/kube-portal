import React, { useState, useEffect } from "react";
import Client, { LinkInfo } from "../drivers/ingress-info";

type Contexts = {
  linkInfo?: LinkInfo[];
};

export const GlobalContext = React.createContext({} as Contexts);

interface GlobalProviderProps {
  getToken: () => Promise<string | null>;
  children: React.ReactNode;
}

export const GlobalProvider: React.FC<GlobalProviderProps> = ({ getToken, children }) => {
  const [linkInfo, setLinkInfo] = useState<LinkInfo[]>();

  const client = new Client();
  const listLinkInfo = async () => {
    try {
      const token = await getToken();
      const resp = await client.List(token ?? undefined);
      setLinkInfo(resp);
    } catch {}
  };

  useEffect(() => {
    listLinkInfo();
    const intervalId = setInterval(listLinkInfo, 10000);
    return () => {
      clearInterval(intervalId);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return <GlobalContext.Provider value={{ linkInfo }}>{children}</GlobalContext.Provider>;
};
