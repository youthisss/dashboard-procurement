"use client";

import { useTable } from "@refinedev/antd";
import { Table, Typography, Card, Button, Form, Input } from "antd";
import { PlusOutlined, EnvironmentOutlined, SearchOutlined } from "@ant-design/icons";
import { EditButton, DeleteButton, ShowButton, Breadcrumb, CreateButton } from "@refinedev/antd";
import { useSearchParams } from "next/navigation";

const { Title, Text } = Typography;

export default function MillList() {
  const searchParams = useSearchParams();
  const q = searchParams.get("q") || "";
  const { tableProps, searchFormProps } = useTable({
    syncWithLocation: true,
    filters: {
      permanent: q
        ? [
            {
              field: "q",
              operator: "eq",
              value: q,
            },
          ]
        : [],
    },
    onSearch: (values: any) => {
      return [
        {
          field: "q",
          operator: "eq",
          value: values.q,
        },
      ];
    },
  });
  const pagination = { ...((tableProps as any).pagination || {}) };
  delete pagination.position;

  return (
    <div className="dashboard-page">
      <div style={{ marginBottom: 16 }}><Breadcrumb /></div>
      <div className="dashboard-page-header">
        <div>
          <Title level={2} style={{ margin: '0 0 8px 0', fontWeight: 700 }}>
            Factory / Mill Locations
          </Title>
        </div>
        <CreateButton
          type="default"
          icon={<PlusOutlined />}
          className="dashboard-action-button"
        >
          Add Mill
        </CreateButton>
      </div>

      <Card variant="borderless" className="no-padding-card dashboard-table-card">
        <div className="dashboard-search-strip">
          <Form {...(searchFormProps as any)} layout="inline" onValuesChange={() => searchFormProps.form?.submit()}>
            <Form.Item name="q" style={{ margin: 0 }}>
              <Input
                prefix={<SearchOutlined style={{ color: "var(--ant-color-text-secondary)", marginRight: 8 }} />}
                placeholder="Search all data..."
                allowClear
                style={{ width: 360, borderRadius: "8px" }}
                onPressEnter={() => searchFormProps.form?.submit()}
              />
            </Form.Item>
          </Form>
        </div>
        <Table
          {...(tableProps as any)}
          pagination={{
            ...pagination,
            placement: ["bottomEnd"],
          }}
          rowKey="id"
          scroll={{ x: 'max-content' }}
        >
          <Table.Column dataIndex="id" title="ID" width={80} />
          <Table.Column dataIndex="code" title="Mill Code" width={200} render={val => <Text strong>{val}</Text>} />
          <Table.Column dataIndex="name" title="Mill Name" />
          <Table.Column
            title="Actions"
            dataIndex="actions"
            width={240}
            align="center"
            render={(_, record: any) => (
              <div style={{ display: "flex", gap: "8px", justifyContent: "center", flexWrap: "wrap" }}>
                <ShowButton size="middle" recordItemId={record.id} />
                <EditButton size="middle" recordItemId={record.id} />
                <DeleteButton size="middle" recordItemId={record.id} />
              </div>
            )}
          />
        </Table>
      </Card>
    </div>
  );
}
