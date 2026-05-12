"use client";

import { useShow } from "@refinedev/core";
import { Breadcrumb, ListButton, EditButton } from "@refinedev/antd";
import { Typography, Space, Card, Descriptions } from "antd";
import { ArrowLeftOutlined } from "@ant-design/icons";
import AppSpinner from "@/components/common/app-spinner";

const { Title, Text } = Typography;

export default function ZoneShow() {
  const { query } = useShow({});
  const { data, isLoading } = query;
  const record = data?.data;

  if (isLoading) {
    return (
      <div style={{ height: "80vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
        <AppSpinner text="Loading zone details..." />
      </div>
    );
  }

  return (
    <div className="dashboard-page-full">
      <div style={{ marginBottom: 16 }}><Breadcrumb /></div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
        <Space size="middle">
          <ListButton icon={<ArrowLeftOutlined />} shape="circle" type="text" hideText />
          <div style={{ display: "flex", flexDirection: "column" }}>
            <Title level={4} style={{ margin: 0, fontWeight: 700 }}>
              {record?.name || "-"}
            </Title>
            <Text type="secondary" style={{ fontSize: 13 }}>Zone ID: {record?.id ?? "-"}</Text>
          </div>
        </Space>
        <Space>
          <EditButton size="large" style={{ borderRadius: "8px", padding: "0 24px", fontWeight: 600 }} />
        </Space>
      </div>

      <Card styles={{ body: { padding: "24px" } }} style={{ borderRadius: 12, border: "1px solid var(--ant-color-border-secondary)" }}>
        <Descriptions title="Zone Information" bordered size="small" column={1}>
          <Descriptions.Item label="Name">{record?.name || "-"}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
}
