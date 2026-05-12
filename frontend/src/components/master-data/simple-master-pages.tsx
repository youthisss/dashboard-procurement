"use client";

import { DeleteButton, EditButton, Breadcrumb, CreateButton, ListButton, SaveButton, useForm, useTable } from "@refinedev/antd";
import { Form, Input, Table, Typography, Card, Space } from "antd";
import { ArrowLeftOutlined, PlusOutlined, SearchOutlined } from "@ant-design/icons";
import { useSearchParams } from "next/navigation";

const { Title, Text } = Typography;

type SimpleMasterListProps = {
  title: string;
  addLabel: string;
  searchPlaceholder: string;
  nameTitle: string;
};

type SimpleMasterFormProps = {
  title: string;
  fieldLabel: string;
  placeholder: string;
  requiredMessage: string;
};

export function SimpleMasterList({ title, addLabel, searchPlaceholder, nameTitle }: SimpleMasterListProps) {
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
    onSearch: (values: any) => [
      {
        field: "q",
        operator: "eq",
        value: values.q,
      },
    ],
  });
  const pagination = { ...((tableProps as any).pagination || {}) };
  delete pagination.position;

  return (
    <div className="dashboard-page">
      <div style={{ marginBottom: 16 }}>
        <Breadcrumb />
      </div>
      <div className="dashboard-page-header">
        <div>
          <Title level={2} style={{ margin: "0 0 8px 0", fontWeight: 700 }}>
            {title}
          </Title>
        </div>
        <CreateButton
          type="default"
          icon={<PlusOutlined />}
          className="dashboard-action-button"
        >
          {addLabel}
        </CreateButton>
      </div>

      <Card
        variant="borderless"
        className="no-padding-card dashboard-table-card"
      >
        <div className="dashboard-search-strip">
          <Form {...(searchFormProps as any)} layout="inline" onValuesChange={() => searchFormProps.form?.submit()}>
            <Form.Item name="q" style={{ margin: 0 }}>
              <Input
                prefix={<SearchOutlined style={{ color: "var(--ant-color-text-secondary)", marginRight: 8 }} />}
                placeholder={searchPlaceholder}
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
          scroll={{ x: "max-content" }}
        >
          <Table.Column dataIndex="id" title="ID" width={80} />
          <Table.Column dataIndex="name" title={nameTitle} render={(val) => <Text strong>{val}</Text>} />
          <Table.Column
            title="Actions"
            dataIndex="actions"
            width={180}
            align="center"
            render={(_, record: any) => (
              <div style={{ display: "flex", gap: "8px", justifyContent: "center", flexWrap: "wrap" }}>
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

export function SimpleMasterCreate({ title, fieldLabel, placeholder, requiredMessage }: SimpleMasterFormProps) {
  const { formProps, saveButtonProps } = useForm({});

  return (
    <SimpleMasterForm
      title={title}
      fieldLabel={fieldLabel}
      placeholder={placeholder}
      requiredMessage={requiredMessage}
      formProps={formProps}
      saveButtonProps={saveButtonProps}
    />
  );
}

export function SimpleMasterEdit({ title, fieldLabel, placeholder, requiredMessage }: SimpleMasterFormProps) {
  const { formProps, saveButtonProps } = useForm({});

  return (
    <SimpleMasterForm
      title={title}
      fieldLabel={fieldLabel}
      placeholder={placeholder}
      requiredMessage={requiredMessage}
      formProps={formProps}
      saveButtonProps={saveButtonProps}
    />
  );
}

function SimpleMasterForm({
  title,
  fieldLabel,
  placeholder,
  requiredMessage,
  formProps,
  saveButtonProps,
}: SimpleMasterFormProps & { formProps: any; saveButtonProps: any }) {
  return (
    <div className="dashboard-page-form">
      <div style={{ marginBottom: 16 }}>
        <Breadcrumb />
      </div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
        <Space>
          <ListButton icon={<ArrowLeftOutlined />} shape="circle" type="text" hideText />
          <Typography.Title level={2} style={{ margin: 0, fontWeight: 700 }}>
            {title}
          </Typography.Title>
        </Space>
      </div>

      <div style={{ maxWidth: 800, margin: "0 auto" }}>
        <Card variant="borderless" className="dashboard-card">
          <Form {...formProps} layout="vertical">
            <Form.Item label={fieldLabel} name={["name"]} rules={[{ required: true, message: requiredMessage }]}>
              <Input placeholder={placeholder} />
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
