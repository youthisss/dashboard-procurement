"use client";

import React, { useEffect, useState } from "react";
import { Grid, Layout, theme } from "antd";
import { CustomSider } from "./sider";
import { Header } from "./header";

const { Content } = Layout;

export const CustomLayout = ({ children }: { children: React.ReactNode }) => {
  const [collapsed, setCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { token } = theme.useToken();
  const screens = Grid.useBreakpoint();
  const isMobile = !screens.md;

  useEffect(() => {
    if (!isMobile) {
      setMobileMenuOpen(false);
    }
  }, [isMobile]);

  return (
    <Layout style={{ minHeight: "100vh", backgroundColor: token.colorBgLayout }}>
      <CustomSider
        collapsed={collapsed}
        onCollapse={setCollapsed}
        isMobile={isMobile}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <Layout
        style={{
          marginLeft: isMobile ? 0 : collapsed ? 80 : 260,
          transition: "margin-left 0.2s",
          backgroundColor: token.colorBgLayout,
          minWidth: 0,
        }}
      >
        <Header
          collapsed={collapsed}
          setCollapsed={setCollapsed}
          isMobile={isMobile}
          onMobileMenuOpen={() => setMobileMenuOpen(true)}
        />
        <Content>
          {children}
        </Content>
      </Layout>
    </Layout>
  );
};
