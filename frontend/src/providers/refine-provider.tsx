"use client";

import { Refine } from "@refinedev/core";
import "@refinedev/antd/dist/reset.css";
import { App as AntdApp, theme } from "antd";
import { ColorModeContextProvider } from "@/contexts/color-mode";
import routerProvider from "@refinedev/nextjs-router";
import React, { useEffect, useMemo, useState } from "react";
import dataProvider from "@refinedev/simple-rest";
import { usePathname, useRouter } from "next/navigation";

import { CustomLayout } from "@/components/layout";
import { clearAuth, getAuthUser, setAuth } from "@/lib/auth";
import { apiClient } from "@/lib/api-client";
import AppSpinner from "@/components/common/app-spinner";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export default function RefineProvider({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const [checkingAuth, setCheckingAuth] = useState(true);
  const isLoginPage = pathname === "/login";

  useEffect(() => {
    let cancelled = false;
    const checkSession = async () => {
      try {
        const { data } = await apiClient.get(`${API_URL}/auth/me`);
        setAuth("", data);
        if (!cancelled && isLoginPage) {
          router.replace("/overview");
          return;
        }
      } catch {
        if (!cancelled && !isLoginPage) {
          router.replace("/login");
          return;
        }
      } finally {
        if (!cancelled) setCheckingAuth(false);
      }
    };
    void checkSession();
    return () => {
      cancelled = true;
    };
  }, [isLoginPage, router]);

  const authProvider = useMemo(
    () => ({
      login: async () => ({ success: true }),
      logout: async () => {
        await apiClient.post(`${API_URL}/auth/logout`);
        clearAuth();
        return { success: true, redirectTo: "/login" };
      },
      check: async () => {
        try {
          const { data } = await apiClient.get(`${API_URL}/auth/me`);
          setAuth("", data);
          return { authenticated: true };
        } catch {
          return { authenticated: false, redirectTo: "/login", logout: true };
        }
      },
      getPermissions: async () => null,
      getIdentity: async () => getAuthUser(),
      onError: async () => ({ error: undefined }),
    }),
    [],
  );

  return (
    <ColorModeContextProvider>
      <AntdApp>
        {checkingAuth ? (
          <SessionCheckingState />
        ) : (
          <RefineAppShell authProvider={authProvider} isLoginPage={isLoginPage}>
            {children}
          </RefineAppShell>
        )}
      </AntdApp>
    </ColorModeContextProvider>
  );
}

function SessionCheckingState() {
  const { token } = theme.useToken();

  return (
    <div
      style={{
        minHeight: "100vh",
        display: "grid",
        placeItems: "center",
        backgroundColor: token.colorBgLayout,
      }}
    >
      <AppSpinner text="Checking session..." />
    </div>
  );
}

function RefineAppShell({
  authProvider,
  children,
  isLoginPage,
}: {
  authProvider: any;
  children: React.ReactNode;
  isLoginPage: boolean;
}) {
  const { notification } = AntdApp.useApp();

  const notificationProvider = useMemo(
    () => ({
      open: (params: any) => {
        const { message, ...rest } = params ?? {};
        notification.open({
          ...rest,
          title: rest.title ?? message,
        });
      },
      close: (key: React.Key) => notification.destroy(key),
      destroy: () => notification.destroy(),
    }),
    [notification],
  );

  return (
    <Refine
      dataProvider={dataProvider(API_URL, apiClient)}
      authProvider={authProvider}
      notificationProvider={notificationProvider}
      routerProvider={routerProvider}
      resources={[
        {
          name: "overview",
          list: "/overview",
          meta: {
            label: "Overview",
            icon: "DashboardOutlined",
          },
        },
        {
          name: "import",
          list: "/import",
          meta: {
            label: "Import Wizard",
          },
        },
        {
          name: "admin/users",
          list: "/admin/users",
          meta: {
            label: "Admin Users",
          },
        },
        {
          name: "vendors",
          list: "/vendors",
          create: "/vendors/create",
          edit: "/vendors/edit/:id",
          show: "/vendors/show/:id",
          meta: {
            label: "Vendors",
            parent: "masterData",
          },
        },
        {
          name: "mills",
          list: "/mills",
          create: "/mills/create",
          edit: "/mills/edit/:id",
          show: "/mills/show/:id",
          meta: {
            label: "Mills",
            parent: "masterData",
          },
        },
        {
          name: "zones",
          list: "/zones",
          create: "/zones/create",
          edit: "/zones/edit/:id",
          show: "/zones/show/:id",
          meta: {
            label: "Zones",
            parent: "masterData",
          },
        },
        {
          name: "products",
          list: "/products",
          create: "/products/create",
          edit: "/products/edit/:id",
          meta: {
            label: "Products",
            parent: "masterData",
          },
        },
        {
          name: "mots",
          list: "/mots",
          create: "/mots/create",
          edit: "/mots/edit/:id",
          meta: {
            label: "MOTs",
            parent: "masterData",
          },
        },
        {
          name: "uoms",
          list: "/uoms",
          create: "/uoms/create",
          edit: "/uoms/edit/:id",
          meta: {
            label: "UOMs",
            parent: "masterData",
          },
        },
        {
          name: "contracts/dedicated-fix",
          list: "/contracts/dedicated-fix",
          create: "/contracts/dedicated-fix/create",
          edit: "/contracts/dedicated-fix/edit/:id",
          show: "/contracts/dedicated-fix/edit/:id",
          meta: {
            label: "Dedicated Fix",
            parent: "contracts",
          },
        },
        {
          name: "contracts/dedicated-var",
          list: "/contracts/dedicated-var",
          create: "/contracts/dedicated-var/create",
          edit: "/contracts/dedicated-var/edit/:id",
          show: "/contracts/dedicated-var/edit/:id",
          meta: {
            label: "Dedicated Var",
            parent: "contracts",
          },
        },
        {
          name: "contracts/oncall",
          list: "/contracts/oncall",
          create: "/contracts/oncall/create",
          edit: "/contracts/oncall/edit/:id",
          show: "/contracts/oncall/edit/:id",
          meta: {
            label: "Oncall Routing",
            parent: "contracts",
          },
        },
      ]}
      options={{
        syncWithLocation: true,
        warnWhenUnsavedChanges: true,
        projectId: "rygell-procurement",
      }}
    >
      {isLoginPage ? children : <CustomLayout>{children}</CustomLayout>}
    </Refine>
  );
}
