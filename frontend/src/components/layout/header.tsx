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
}: {
  collapsed: boolean;
  setCollapsed: (val: boolean) => void;
}) => {
  const { mode, setMode } = useContext(ColorModeContext);

  return (
    <AntdHeader
      style={{
        backgroundColor: "var(--ant-color-bg-container)",
        borderBottom: "1px solid var(--ant-color-border)",
        display: "flex",
        justifyContent: "flex-end",
        alignItems: "center",
        padding: "0 24px",
        height: "64px",
        position: "sticky",
        top: 0,
        zIndex: 998,
      }}
    >
      <Button
        type="text"
        icon={collapsed ? <MenuUnfoldOutlined style={{ fontSize: "18px" }} /> : <MenuFoldOutlined style={{ fontSize: "18px" }} />}
        onClick={() => setCollapsed(!collapsed)}
        style={{ marginRight: 16 }}
      />

      <div style={{ flex: 1, display: "flex", alignItems: "center" }}>
        <GlobalSearchBar />
      </div>

      <Space size="middle" style={{ alignItems: "center" }}>
        <Typography.Text type="secondary" strong style={{ fontSize: "13px" }}>
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
