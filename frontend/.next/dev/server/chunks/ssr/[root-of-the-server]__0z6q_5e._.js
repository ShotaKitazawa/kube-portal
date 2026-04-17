module.exports = [
"[externals]/react/jsx-dev-runtime [external] (react/jsx-dev-runtime, cjs)", ((__turbopack_context__, module, exports) => {

const mod = __turbopack_context__.x("react/jsx-dev-runtime", () => require("react/jsx-dev-runtime"));

module.exports = mod;
}),
"[externals]/react [external] (react, cjs)", ((__turbopack_context__, module, exports) => {

const mod = __turbopack_context__.x("react", () => require("react"));

module.exports = mod;
}),
"[externals]/react/jsx-runtime [external] (react/jsx-runtime, cjs)", ((__turbopack_context__, module, exports) => {

const mod = __turbopack_context__.x("react/jsx-runtime", () => require("react/jsx-runtime"));

module.exports = mod;
}),
"[project]/src/drivers/auth/index.ts [ssr] (ecmascript)", ((__turbopack_context__) => {
"use strict";

return __turbopack_context__.a(async (__turbopack_handle_async_dependencies__, __turbopack_async_result__) => { try {

__turbopack_context__.s([
    "loadOIDCSetup",
    ()=>loadOIDCSetup
]);
var __TURBOPACK__imported__module__$5b$externals$5d2f$oidc$2d$client$2d$ts__$5b$external$5d$__$28$oidc$2d$client$2d$ts$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$oidc$2d$client$2d$ts$40$3$2e$5$2e$0$2f$node_modules$2f$oidc$2d$client$2d$ts$29$__ = __turbopack_context__.i("[externals]/oidc-client-ts [external] (oidc-client-ts, esm_import, [project]/node_modules/.pnpm/oidc-client-ts@3.5.0/node_modules/oidc-client-ts)");
var __turbopack_async_dependencies__ = __turbopack_handle_async_dependencies__([
    __TURBOPACK__imported__module__$5b$externals$5d2f$oidc$2d$client$2d$ts__$5b$external$5d$__$28$oidc$2d$client$2d$ts$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$oidc$2d$client$2d$ts$40$3$2e$5$2e$0$2f$node_modules$2f$oidc$2d$client$2d$ts$29$__
]);
[__TURBOPACK__imported__module__$5b$externals$5d2f$oidc$2d$client$2d$ts__$5b$external$5d$__$28$oidc$2d$client$2d$ts$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$oidc$2d$client$2d$ts$40$3$2e$5$2e$0$2f$node_modules$2f$oidc$2d$client$2d$ts$29$__] = __turbopack_async_dependencies__.then ? (await __turbopack_async_dependencies__)() : __turbopack_async_dependencies__;
;
async function loadOIDCSetup() {
    const baseUrl = ("TURBOPACK compile-time truthy", 1) ? ("TURBOPACK compile-time value", "http://localhost:8080") + '/api' : "TURBOPACK unreachable";
    let config;
    try {
        const res = await fetch(`${baseUrl}/oidc-config`);
        config = await res.json();
    } catch  {
        return {
            configured: false,
            userManager: null
        };
    }
    if (!config.enabled || !config.issuer || !config.client_id) {
        return {
            configured: false,
            userManager: null
        };
    }
    const userManager = new __TURBOPACK__imported__module__$5b$externals$5d2f$oidc$2d$client$2d$ts__$5b$external$5d$__$28$oidc$2d$client$2d$ts$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$oidc$2d$client$2d$ts$40$3$2e$5$2e$0$2f$node_modules$2f$oidc$2d$client$2d$ts$29$__["UserManager"]({
        authority: config.issuer,
        client_id: config.client_id,
        redirect_uri: `${window.location.origin}/callback`,
        scope: 'openid profile email offline_access',
        userStore: new __TURBOPACK__imported__module__$5b$externals$5d2f$oidc$2d$client$2d$ts__$5b$external$5d$__$28$oidc$2d$client$2d$ts$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$oidc$2d$client$2d$ts$40$3$2e$5$2e$0$2f$node_modules$2f$oidc$2d$client$2d$ts$29$__["WebStorageStateStore"]({
            store: window.localStorage
        }),
        automaticSilentRenew: false
    });
    return {
        configured: true,
        userManager
    };
}
__turbopack_async_result__();
} catch(e) { __turbopack_async_result__(e); } }, false);}),
"[project]/src/contexts/auth.tsx [ssr] (ecmascript)", ((__turbopack_context__) => {
"use strict";

return __turbopack_context__.a(async (__turbopack_handle_async_dependencies__, __turbopack_async_result__) => { try {

__turbopack_context__.s([
    "AuthContext",
    ()=>AuthContext,
    "AuthProvider",
    ()=>AuthProvider,
    "useAuth",
    ()=>useAuth
]);
var __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__ = __turbopack_context__.i("[externals]/react/jsx-dev-runtime [external] (react/jsx-dev-runtime, cjs)");
var __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__ = __turbopack_context__.i("[externals]/react [external] (react, cjs)");
var __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$auth$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/src/drivers/auth/index.ts [ssr] (ecmascript)");
var __turbopack_async_dependencies__ = __turbopack_handle_async_dependencies__([
    __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$auth$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__
]);
[__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$auth$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__] = __turbopack_async_dependencies__.then ? (await __turbopack_async_dependencies__)() : __turbopack_async_dependencies__;
;
;
;
const AuthContext = /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["createContext"])({
    user: null,
    login: async ()=>{},
    logout: async ()=>{}
});
const baseUrl = ("TURBOPACK compile-time truthy", 1) ? ("TURBOPACK compile-time value", "http://localhost:8080") + '/api' : "TURBOPACK unreachable";
async function fetchUserinfo() {
    try {
        const res = await fetch(`${baseUrl}/userinfo`);
        if (!res.ok) return null;
        const data = await res.json();
        return {
            access_token: '',
            profile: {
                sub: data.sub,
                name: data.name,
                email: data.email,
                picture: data.picture
            }
        };
    } catch  {
        return null;
    }
}
function toAuthUser(u) {
    return {
        access_token: u.access_token,
        profile: {
            sub: u.profile.sub,
            name: u.profile.name ?? undefined,
            email: u.profile.email ?? undefined,
            picture: u.profile.picture ?? undefined
        }
    };
}
const AuthProvider = ({ children })=>{
    const [user, setUser] = (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["useState"])(null);
    const [userManager, setUserManager] = (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["useState"])(null);
    (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["useEffect"])(()=>{
        (0, __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$auth$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__["loadOIDCSetup"])().then((setup)=>{
            if (!setup.configured) {
                fetchUserinfo().then(setUser);
                return;
            }
            const mgr = setup.userManager;
            setUserManager(mgr);
            mgr.getUser().then((u)=>{
                if (u) setUser(toAuthUser(u));
            });
            const handleUserLoaded = (u)=>setUser(toAuthUser(u));
            const handleUserUnloaded = ()=>setUser(null);
            mgr.events.addUserLoaded(handleUserLoaded);
            mgr.events.addUserUnloaded(handleUserUnloaded);
            return ()=>{
                mgr.events.removeUserLoaded(handleUserLoaded);
                mgr.events.removeUserUnloaded(handleUserUnloaded);
            };
        });
    }, []);
    const login = async ()=>{
        await userManager?.signinRedirect();
    };
    const logout = async ()=>{
        if (userManager) {
            await userManager.removeUser();
        }
        setUser(null);
    };
    return /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__["jsxDEV"])(AuthContext.Provider, {
        value: {
            user,
            login,
            logout
        },
        children: children
    }, void 0, false, {
        fileName: "[project]/src/contexts/auth.tsx",
        lineNumber: 111,
        columnNumber: 5
    }, ("TURBOPACK compile-time value", void 0));
};
const useAuth = ()=>(0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["useContext"])(AuthContext);
__turbopack_async_result__();
} catch(e) { __turbopack_async_result__(e); } }, false);}),
"[project]/src/drivers/ingress-info/index.ts [ssr] (ecmascript)", ((__turbopack_context__) => {
"use strict";

return __turbopack_context__.a(async (__turbopack_handle_async_dependencies__, __turbopack_async_result__) => { try {

__turbopack_context__.s([
    "default",
    ()=>__TURBOPACK__default__export__
]);
var __TURBOPACK__imported__module__$5b$externals$5d2f$openapi$2d$fetch__$5b$external$5d$__$28$openapi$2d$fetch$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$openapi$2d$fetch$40$0$2e$17$2e$0$2f$node_modules$2f$openapi$2d$fetch$29$__ = __turbopack_context__.i("[externals]/openapi-fetch [external] (openapi-fetch, esm_import, [project]/node_modules/.pnpm/openapi-fetch@0.17.0/node_modules/openapi-fetch)");
var __turbopack_async_dependencies__ = __turbopack_handle_async_dependencies__([
    __TURBOPACK__imported__module__$5b$externals$5d2f$openapi$2d$fetch__$5b$external$5d$__$28$openapi$2d$fetch$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$openapi$2d$fetch$40$0$2e$17$2e$0$2f$node_modules$2f$openapi$2d$fetch$29$__
]);
[__TURBOPACK__imported__module__$5b$externals$5d2f$openapi$2d$fetch__$5b$external$5d$__$28$openapi$2d$fetch$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$openapi$2d$fetch$40$0$2e$17$2e$0$2f$node_modules$2f$openapi$2d$fetch$29$__] = __turbopack_async_dependencies__.then ? (await __turbopack_async_dependencies__)() : __turbopack_async_dependencies__;
;
class Client {
    client;
    constructor(){
        const baseUrl = ("TURBOPACK compile-time truthy", 1) ? ("TURBOPACK compile-time value", "http://localhost:8080") + '/api' : "TURBOPACK unreachable";
        this.client = (0, __TURBOPACK__imported__module__$5b$externals$5d2f$openapi$2d$fetch__$5b$external$5d$__$28$openapi$2d$fetch$2c$__esm_import$2c$__$5b$project$5d2f$node_modules$2f2e$pnpm$2f$openapi$2d$fetch$40$0$2e$17$2e$0$2f$node_modules$2f$openapi$2d$fetch$29$__["default"])({
            baseUrl
        });
    }
    async List(accessToken) {
        const headers = {};
        if (accessToken) {
            headers['Authorization'] = `Bearer ${accessToken}`;
        }
        const { data } = await this.client.GET('/list', {
            headers
        }).catch(()=>({
                data: null
            }));
        return data ?? [];
    }
}
const __TURBOPACK__default__export__ = Client;
__turbopack_async_result__();
} catch(e) { __turbopack_async_result__(e); } }, false);}),
"[project]/src/contexts/global.tsx [ssr] (ecmascript)", ((__turbopack_context__) => {
"use strict";

return __turbopack_context__.a(async (__turbopack_handle_async_dependencies__, __turbopack_async_result__) => { try {

__turbopack_context__.s([
    "GlobalContext",
    ()=>GlobalContext,
    "GlobalProvider",
    ()=>GlobalProvider
]);
var __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__ = __turbopack_context__.i("[externals]/react/jsx-dev-runtime [external] (react/jsx-dev-runtime, cjs)");
var __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__ = __turbopack_context__.i("[externals]/react [external] (react, cjs)");
var __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$ingress$2d$info$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/src/drivers/ingress-info/index.ts [ssr] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/src/contexts/auth.tsx [ssr] (ecmascript)");
var __turbopack_async_dependencies__ = __turbopack_handle_async_dependencies__([
    __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$ingress$2d$info$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__,
    __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__
]);
[__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$ingress$2d$info$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__, __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__] = __turbopack_async_dependencies__.then ? (await __turbopack_async_dependencies__)() : __turbopack_async_dependencies__;
;
;
;
;
const GlobalContext = /*#__PURE__*/ __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["default"].createContext({});
const GlobalProvider = (props)=>{
    const [linkInfo, setLinkInfo] = (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["useState"])();
    const { user } = (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["useContext"])(__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__["AuthContext"]);
    const client = new __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$drivers$2f$ingress$2d$info$2f$index$2e$ts__$5b$ssr$5d$__$28$ecmascript$29$__["default"]();
    const listLinkInfo = async ()=>{
        try {
            const resp = await client.List(user?.access_token);
            setLinkInfo(resp);
        } catch (error) {}
    };
    (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react__$5b$external$5d$__$28$react$2c$__cjs$29$__["useEffect"])(()=>{
        listLinkInfo();
        const intervalId = setInterval(()=>{
            listLinkInfo();
        }, 10000);
        return ()=>{
            clearInterval(intervalId);
        };
    }, [
        user
    ]);
    return /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__["jsxDEV"])(GlobalContext.Provider, {
        value: {
            linkInfo
        },
        children: props.children
    }, void 0, false, {
        fileName: "[project]/src/contexts/global.tsx",
        lineNumber: 39,
        columnNumber: 5
    }, ("TURBOPACK compile-time value", void 0));
};
__turbopack_async_result__();
} catch(e) { __turbopack_async_result__(e); } }, false);}),
"[project]/src/pages/_app.tsx [ssr] (ecmascript)", ((__turbopack_context__) => {
"use strict";

return __turbopack_context__.a(async (__turbopack_handle_async_dependencies__, __turbopack_async_result__) => { try {

__turbopack_context__.s([
    "default",
    ()=>__TURBOPACK__default__export__
]);
var __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__ = __turbopack_context__.i("[externals]/react/jsx-dev-runtime [external] (react/jsx-dev-runtime, cjs)");
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f2e$pnpm$2f40$mui$2b$material$40$7$2e$3$2e$10_$40$emotion$2b$react$40$11$2e$14$2e$0_$40$types$2b$react$40$19$2e$2$2e$14_react$40$19$2e$2$2e$4_$5f40$emotion_0bf69b005fb0f7e0fb4059f69d3f0654$2f$node_modules$2f40$mui$2f$material$2f$esm$2f$styles$2f$ThemeProvider$2e$js__$5b$ssr$5d$__$28$ecmascript$29$__$3c$export__default__as__ThemeProvider$3e$__ = __turbopack_context__.i("[project]/node_modules/.pnpm/@mui+material@7.3.10_@emotion+react@11.14.0_@types+react@19.2.14_react@19.2.4__@emotion_0bf69b005fb0f7e0fb4059f69d3f0654/node_modules/@mui/material/esm/styles/ThemeProvider.js [ssr] (ecmascript) <export default as ThemeProvider>");
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f2e$pnpm$2f40$mui$2b$material$40$7$2e$3$2e$10_$40$emotion$2b$react$40$11$2e$14$2e$0_$40$types$2b$react$40$19$2e$2$2e$14_react$40$19$2e$2$2e$4_$5f40$emotion_0bf69b005fb0f7e0fb4059f69d3f0654$2f$node_modules$2f40$mui$2f$material$2f$esm$2f$styles$2f$createTheme$2e$js__$5b$ssr$5d$__$28$ecmascript$29$__$3c$export__default__as__createTheme$3e$__ = __turbopack_context__.i("[project]/node_modules/.pnpm/@mui+material@7.3.10_@emotion+react@11.14.0_@types+react@19.2.14_react@19.2.4__@emotion_0bf69b005fb0f7e0fb4059f69d3f0654/node_modules/@mui/material/esm/styles/createTheme.js [ssr] (ecmascript) <export default as createTheme>");
var __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f2e$pnpm$2f40$mui$2b$material$40$7$2e$3$2e$10_$40$emotion$2b$react$40$11$2e$14$2e$0_$40$types$2b$react$40$19$2e$2$2e$14_react$40$19$2e$2$2e$4_$5f40$emotion_0bf69b005fb0f7e0fb4059f69d3f0654$2f$node_modules$2f40$mui$2f$material$2f$esm$2f$CssBaseline$2f$CssBaseline$2e$js__$5b$ssr$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/node_modules/.pnpm/@mui+material@7.3.10_@emotion+react@11.14.0_@types+react@19.2.14_react@19.2.4__@emotion_0bf69b005fb0f7e0fb4059f69d3f0654/node_modules/@mui/material/esm/CssBaseline/CssBaseline.js [ssr] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/src/contexts/auth.tsx [ssr] (ecmascript)");
var __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$global$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__ = __turbopack_context__.i("[project]/src/contexts/global.tsx [ssr] (ecmascript)");
var __turbopack_async_dependencies__ = __turbopack_handle_async_dependencies__([
    __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__,
    __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$global$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__
]);
[__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__, __TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$global$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__] = __turbopack_async_dependencies__.then ? (await __turbopack_async_dependencies__)() : __turbopack_async_dependencies__;
;
;
;
;
;
;
;
const theme = (0, __TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f2e$pnpm$2f40$mui$2b$material$40$7$2e$3$2e$10_$40$emotion$2b$react$40$11$2e$14$2e$0_$40$types$2b$react$40$19$2e$2$2e$14_react$40$19$2e$2$2e$4_$5f40$emotion_0bf69b005fb0f7e0fb4059f69d3f0654$2f$node_modules$2f40$mui$2f$material$2f$esm$2f$styles$2f$createTheme$2e$js__$5b$ssr$5d$__$28$ecmascript$29$__$3c$export__default__as__createTheme$3e$__["createTheme"])({
    palette: {
        primary: {
            main: '#3f51b5'
        },
        secondary: {
            main: '#f50057'
        }
    }
});
function MyApp({ Component, pageProps }) {
    return /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f2e$pnpm$2f40$mui$2b$material$40$7$2e$3$2e$10_$40$emotion$2b$react$40$11$2e$14$2e$0_$40$types$2b$react$40$19$2e$2$2e$14_react$40$19$2e$2$2e$4_$5f40$emotion_0bf69b005fb0f7e0fb4059f69d3f0654$2f$node_modules$2f40$mui$2f$material$2f$esm$2f$styles$2f$ThemeProvider$2e$js__$5b$ssr$5d$__$28$ecmascript$29$__$3c$export__default__as__ThemeProvider$3e$__["ThemeProvider"], {
        theme: theme,
        children: [
            /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$node_modules$2f2e$pnpm$2f40$mui$2b$material$40$7$2e$3$2e$10_$40$emotion$2b$react$40$11$2e$14$2e$0_$40$types$2b$react$40$19$2e$2$2e$14_react$40$19$2e$2$2e$4_$5f40$emotion_0bf69b005fb0f7e0fb4059f69d3f0654$2f$node_modules$2f40$mui$2f$material$2f$esm$2f$CssBaseline$2f$CssBaseline$2e$js__$5b$ssr$5d$__$28$ecmascript$29$__["default"], {}, void 0, false, {
                fileName: "[project]/src/pages/_app.tsx",
                lineNumber: 24,
                columnNumber: 7
            }, this),
            " ",
            /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$auth$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__["AuthProvider"], {
                children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__["jsxDEV"])(__TURBOPACK__imported__module__$5b$project$5d2f$src$2f$contexts$2f$global$2e$tsx__$5b$ssr$5d$__$28$ecmascript$29$__["GlobalProvider"], {
                    children: /*#__PURE__*/ (0, __TURBOPACK__imported__module__$5b$externals$5d2f$react$2f$jsx$2d$dev$2d$runtime__$5b$external$5d$__$28$react$2f$jsx$2d$dev$2d$runtime$2c$__cjs$29$__["jsxDEV"])(Component, {
                        ...pageProps
                    }, void 0, false, {
                        fileName: "[project]/src/pages/_app.tsx",
                        lineNumber: 27,
                        columnNumber: 11
                    }, this)
                }, void 0, false, {
                    fileName: "[project]/src/pages/_app.tsx",
                    lineNumber: 26,
                    columnNumber: 9
                }, this)
            }, void 0, false, {
                fileName: "[project]/src/pages/_app.tsx",
                lineNumber: 25,
                columnNumber: 7
            }, this)
        ]
    }, void 0, true, {
        fileName: "[project]/src/pages/_app.tsx",
        lineNumber: 23,
        columnNumber: 5
    }, this);
}
const __TURBOPACK__default__export__ = MyApp;
__turbopack_async_result__();
} catch(e) { __turbopack_async_result__(e); } }, false);}),
];

//# sourceMappingURL=%5Broot-of-the-server%5D__0z6q_5e._.js.map