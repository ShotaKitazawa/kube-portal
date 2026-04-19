import React, { useEffect } from "react";
import { loadOIDCSetup } from "./drivers/auth";
import { Layout } from "./components/Layout";
import { LinkGroups } from "./components/LinkGroups";

function App() {
  useEffect(() => {
    if (!window.location.search.includes("code=")) return;
    loadOIDCSetup().then((setup) => {
      if (!setup.configured) {
        window.location.replace("/");
        return;
      }
      setup.userManager
        .signinRedirectCallback()
        .then(() => window.location.replace("/"))
        .catch(console.error);
    });
  }, []);

  if (window.location.search.includes("code=")) return null;

  return (
    <Layout>
      <LinkGroups />
    </Layout>
  );
}

export default App;
