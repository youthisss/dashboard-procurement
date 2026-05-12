"use client";

import { useForm, Breadcrumb, ListButton, SaveButton } from "@refinedev/antd";
import { Form, Input, Typography, Space, Card, Select } from "antd";
import { ArrowLeftOutlined } from "@ant-design/icons";

export default function ZoneCreate() {
  const { formProps, saveButtonProps } = useForm({});

  return (
    <div className="dashboard-page-form">
      <div style={{ marginBottom: 16 }}><Breadcrumb /></div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
        <Space>
          <ListButton icon={<ArrowLeftOutlined />} shape="circle" type="text" hideText />
          <Typography.Title level={2} style={{ margin: 0, fontWeight: 700 }}>Create Zone</Typography.Title>
        </Space>
      </div>

      <div style={{ maxWidth: 800, margin: "0 auto" }}>
        <Card variant="borderless" className="dashboard-card">
          <Form {...formProps} layout="vertical">
            <Form.Item
              label="Zone Name"
              name={["name"]}
              rules={[{ required: true, message: "Zone name is required" }]}
            >
              <Input placeholder="e.g. JAMBI" />
            </Form.Item>
            <Form.Item
              label="Zone Type"
              name={["type"]}
              rules={[{ required: true, message: "Zone type is required" }]}
            >
              <Select
                placeholder="Select zone type"
                options={[
                  { label: "Origin", value: "origin" },
                  { label: "Destination", value: "destination" },
                ]}
              />
            </Form.Item>
            <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 24 }}>
              <SaveButton {...saveButtonProps} size="large" />
            </div>
          </Form>
        </Card>
      </div>
    </div>
  );
}
