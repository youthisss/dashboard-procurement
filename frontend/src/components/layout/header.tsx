"use client";

import React, { useContext } from "react";
import { Layout, Button, Space, Typography } from "antd";
import {
  BellOutlined,
  SunOutlined,
  MoonOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from "@ant-design/icons";
import { ColorModeContext } from "@/contexts/color-mode";
import dayjs from "dayjs";
import { GlobalSearchBar } from "@/components/layout/global-search-bar";

const { Header: AntdHeader } = Layout;

export const Header = ({
  collapsed,
  setCollapsed,
  isMobile,
  onMobileMenuOpen,
}: {
  collapsed: boolean;
  setCollapsed: (val: boolean) => void;
  isMobile: boolean;
  onMobileMenuOpen: () => void;
}) => {
  const { mode, setMode } = useContext(ColorModeContext);

  return (
    <AntdHeader
      style={{
        backgroundColor: "var(--ant-color-bg-container)",
        borderBottom: "1px solid var(--ant-color-border)",
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        flexWrap: isMobile ? "wrap" : "nowrap",
        gap: isMobile ? 8 : 0,
        padding: isMobile ? "8px 12px" : "0 24px",
        minHeight: "64px",
        height: isMobile ? "auto" : "64px",
        lineHeight: isMobile ? "normal" : "64px",
        position: "sticky",
        top: 0,
        zIndex: 998,
      }}
    >
      <Button
        type="text"
        icon={
          isMobile
            ? <MenuUnfoldOutlined style={{ fontSize: "18px" }} />
            : collapsed
              ? <MenuUnfoldOutlined style={{ fontSize: "18px" }} />
              : <MenuFoldOutlined style={{ fontSize: "18px" }} />
        }
        onClick={() => {
          if (isMobile) {
            onMobileMenuOpen();
            return;
          }
          setCollapsed(!collapsed);
        }}
        style={{ flexShrink: 0 }}
      />

      <div
        style={{
          flex: isMobile ? "1 0 calc(100% - 56px)" : 1,
          order: isMobile ? 3 : 0,
          display: "flex",
          alignItems: "center",
          lineHeight: "normal",
          minWidth: 0,
        }}
      >
        <GlobalSearchBar />
      </div>

      <Space size={isMobile ? "small" : "middle"} style={{ alignItems: "center", flexShrink: 0, marginLeft: isMobile ? "auto" : 16 }}>
        <Typography.Text type="secondary" strong style={{ display: isMobile ? "none" : "inline", fontSize: "13px" }}>
          {dayjs().format("dddd, DD MMM YYYY")}
        </Typography.Text>

        <Button
          type="text"
          icon={mode === "light" ? <MoonOutlined style={{ fontSize: "18px" }} /> : <SunOutlined style={{ fontSize: "18px" }} />}
          onClick={() => setMode(mode === "light" ? "dark" : "light")}
        />

        <Button type="text" icon={<BellOutlined style={{ fontSize: "18px" }} />} />
      </Space>
    </AntdHeader>
  );
};
