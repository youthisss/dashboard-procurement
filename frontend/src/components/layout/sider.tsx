"use client";

import React from "react";
import { Layout, Menu, Typography, Button } from "antd";
import { usePathname, useRouter } from "next/navigation";
import {
  DashboardOutlined,
  CloudUploadOutlined,
  ShopOutlined,
  EnvironmentOutlined,
  FileDoneOutlined,
  CarOutlined,
  LogoutOutlined,
  NodeIndexOutlined,
  TableOutlined,
  UserAddOutlined,
} from "@ant-design/icons";
import { clearAuth, getAuthUser } from "@/lib/auth";
import { apiClient } from "@/lib/api-client";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

const { Sider } = Layout;
const { Text } = Typography;
const BRAND_ICON_SIZE = 24;
const BRAND_ICON_FRAME = 36;

export const CustomSider = ({ collapsed, onCollapse }: { collapsed: boolean, onCollapse: (val: boolean) => void }) => {
  const pathname = usePathname();
  const router = useRouter();
  const authUser = getAuthUser();
  const profileName = authUser?.name || authUser?.username || "User";
  const profileRole = authUser?.is_admin ? "Admin" : "Staff";

  const menuItems = [
    {
      key: "/overview",
      icon: <DashboardOutlined style={{ fontSize: '18px' }} />,
      label: "Overview",
    },
    {
      key: "/vendors",
      icon: <ShopOutlined style={{ fontSize: '18px' }} />,
      label: "Vendors",
    },
    {
      key: "/mills",
      icon: <EnvironmentOutlined style={{ fontSize: '18px' }} />,
      label: "Mills",
    },
    {
      key: "/contracts/dedicated-fix",
      icon: <FileDoneOutlined style={{ fontSize: '18px' }} />,
      label: "Dedicated Fix",
    },
    {
      key: "/contracts/dedicated-var",
      icon: <NodeIndexOutlined style={{ fontSize: '18px' }} />,
      label: "Dedicated Var",
    },
    {
      key: "/contracts/oncall",
      icon: <CarOutlined style={{ fontSize: '18px' }} />,
      label: "Oncall Routing",
    },
    {
      key: "/explorer",
      icon: <TableOutlined style={{ fontSize: '18px' }} />,
      label: "Data Explorer",
    },
    {
      key: "/import",
      icon: <CloudUploadOutlined style={{ fontSize: '18px' }} />,
      label: "Import Wizard",
    },
    ...(authUser?.is_admin
      ? [
          {
            key: "/admin/users",
            icon: <UserAddOutlined style={{ fontSize: "18px" }} />,
            label: "Admin Users",
          },
        ]
      : []),
  ];

  return (
    <>
      <style>{`
        .custom-menu .ant-menu-item-group-title {
          transition: all 0.2s;
        }
        ${collapsed ? `
        .custom-menu .ant-menu-item-group-title {
          padding: 0 !important;
          height: 0 !important;
          opacity: 0 !important;
          margin: 0 !important;
          overflow: hidden !important;
        }
        ` : ''}
      `}</style>
      <Sider
        collapsible
      collapsed={collapsed}
      trigger={null}
      width={260}
      collapsedWidth={80}
      style={{
        borderRight: "1px solid var(--ant-color-border)",
        backgroundColor: "var(--ant-color-bg-container)",
        height: "100vh",
        position: "fixed",
        left: 0,
        top: 0,
        bottom: 0,
        zIndex: 999,
      }}
    >
      <div style={{ display: 'flex', flexDirection: 'column', height: '100%', overflow: 'hidden' }}>
        <div style={{ 
          padding: collapsed ? "24px 0" : "24px", 
          display: "flex", 
          alignItems: "center", 
          justifyContent: "center", 
          gap: collapsed ? 0 : "12px", 
          height: "80px", 
          flexShrink: 0,
          transition: "all 0.2s"
        }}>
          <div style={{ width: BRAND_ICON_FRAME, height: BRAND_ICON_FRAME, borderRadius: 10, display: "flex", alignItems: "center", justifyContent: "center", flexShrink: 0 }}>
            <svg width={BRAND_ICON_SIZE} height={BRAND_ICON_SIZE} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ color: "var(--ant-color-text)", flexShrink: 0 }}>
              <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
              <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
              <line x1="12" y1="22.08" x2="12" y2="12"></line>
            </svg>
          </div>
          <span style={{ 
            display: "inline-block",
            fontSize: "20px", fontWeight: 700, color: "var(--ant-color-text)", 
            whiteSpace: "nowrap", opacity: collapsed ? 0 : 1, width: collapsed ? 0 : 130, 
            overflow: "hidden", transition: "all 0.2s" 
          }}>
            Procurement
          </span>
        </div>

      <div style={{ flex: 1, padding: "0 12px", overflowY: "auto", overflowX: "hidden" }}>
        <Menu
          className="custom-menu"
          mode="inline"
          selectedKeys={[pathname]}
          onClick={({ key }) => router.push(key)}
          style={{ borderRight: "none", fontSize: "14px", fontWeight: 500 }}
          items={menuItems.map(item => ({
            ...item,
            style: {
              borderRadius: "8px",
              marginBottom: "4px",
              fontWeight: pathname === item.key ? 600 : 500,
              backgroundColor: pathname === item.key ? "rgba(255,255,255,0.08)" : "transparent",
              color: pathname === item.key ? "var(--ant-color-text)" : "var(--ant-color-text-secondary)",
            },
          }))}
        />
      </div>

        <div style={{ padding: collapsed ? "20px 0" : "20px", display: "flex", justifyContent: "center", borderTop: "1px solid var(--ant-color-border)", flexShrink: 0 }}>
          <div style={{ display: "flex", flexDirection: "column", gap: "12px", width: "100%", alignItems: collapsed ? "center" : "stretch" }}>
            <div style={{ 
              opacity: collapsed ? 0 : 1, 
              height: collapsed ? 0 : 'auto', 
              overflow: "hidden", 
              transition: "all 0.2s",
              whiteSpace: "nowrap"
            }}>
              <Text style={{ display: "block", fontWeight: 600 }}>{profileName}</Text>
              <Text type="secondary" style={{ fontSize: "13px" }}>{profileRole}</Text>
            </div>
            
            <Button 
              type="text" 
              icon={<LogoutOutlined style={{ fontSize: '18px', color: '#ef4444' }} />} 
              title="Logout"
              onClick={async () => {
                try {
                  await apiClient.post(`${API_URL}/auth/logout`);
                } finally {
                  clearAuth();
                  router.push("/login");
                }
              }}
              style={{ 
                width: collapsed ? "40px" : "100%", 
                height: collapsed ? "40px" : undefined,
                padding: collapsed ? 0 : "4px 15px",
                textAlign: collapsed ? "center" : "left", 
                justifyContent: collapsed ? "center" : "flex-start",
                background: "transparent",
                color: "#ef4444",
                borderRadius: "8px",
                fontWeight: 500,
                transition: "all 0.2s"
              }}
            >
              <span style={{ 
                display: "inline-block",
                opacity: collapsed ? 0 : 1, 
                width: collapsed ? 0 : 'auto', 
                overflow: 'hidden', 
                transition: "all 0.2s", 
                marginLeft: collapsed ? 0 : 8 
              }}>
                Logout
              </span>
            </Button>
          </div>
        </div>
      </div>
    </Sider>
    </>
  );
};
