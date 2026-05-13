"use client";

import { useTable } from "@refinedev/antd";
import { useInvalidate } from "@refinedev/core";
import { Table, Space, Typography, Card, Form, Input } from "antd";
import { PhoneOutlined, PlusOutlined, SearchOutlined } from "@ant-design/icons";
import { EditButton, DeleteButton, Breadcrumb, CreateButton } from "@refinedev/antd";
import { useSearchParams } from "next/navigation";
import ContractModeSwitch from "@/components/contracts/ContractModeSwitch";

const { Title, Text } = Typography;

export default function OncallList() {
  const searchParams = useSearchParams();
  const q = searchParams.get("q") || "";
  const { tableProps, searchFormProps } = useTable({
    syncWithLocation: true,
    meta: {
      populate: ["vendor", "mill", "product", "origin_zone", "dest_zone", "mot", "uom"]
    },
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

  const formatIDR = (val: any) => {
    if (!val) return "-";
    return new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", maximumFractionDigits: 0 }).format(val);
  };

  const formatDate = (dateStr: string | null | undefined) => {
    if (!dateStr || dateStr === "-") return "-";
    const dateOnly = dateStr.split('T')[0];
    const [year, month, day] = dateOnly.split('-');
    if (!year || !month || !day) return dateStr;
    return `${day}/${month}/${year}`;
  };
  const pagination = { ...((tableProps as any).pagination || {}) };
  delete pagination.position;

  return (
    <div className="dashboard-page">
      <div style={{ marginBottom: 16 }}><Breadcrumb /></div>
      <div className="dashboard-page-header">
        <div>
          <Title level={2} style={{ margin: '0 0 8px 0', fontWeight: 700 }}>
            Oncall Routings
          </Title>
        </div>
        <Space size={12} wrap>
          <ContractModeSwitch />
          <CreateButton
            type="default"
            icon={<PlusOutlined />}
            className="dashboard-action-button"
          >
            Create Contract
          </CreateButton>
        </Space>
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
          {...tableProps}
          pagination={{
            ...pagination,
            placement: ["bottomEnd"],
          }}
          rowKey="id"
          scroll={{ x: 'max-content' }}
        >
          <Table.Column dataIndex="id" title="ID" width={60} />
          <Table.Column dataIndex="area_category" title="Area/Category" width={150} />
          <Table.Column dataIndex={["mill", "name"]} title="Mill" width={150} />
          <Table.Column dataIndex={["product", "name"]} title="Product" width={150} />
          <Table.Column title="Contract Type" width={150} render={() => <Text>Oncall</Text>} />
          <Table.Column dataIndex="proposal_cfas" title="Proposal/CFAS" width={150} />
          <Table.Column dataIndex="spk_number" title="SPK Number" width={150} />
          <Table.Column dataIndex="fa_number" title="FA Number" width={150} />

          <Table.Column
            title="Validity"
            width={220}
            render={(_, record: any) => (
              <Space orientation="vertical" size="small">
                <Text style={{ fontSize: 12 }}>Start: {formatDate(record.validity_start)}</Text>
                <Text style={{ fontSize: 12 }}>End: {formatDate(record.validity_end)}</Text>
              </Space>
            )}
          />
          <Table.Column dataIndex={["vendor", "code"]} title="Vendor Code" width={120} />
          <Table.Column dataIndex={["vendor", "name"]} title="Transporter/Carrier" width={220} />

          <Table.Column
            title="Route"
            width={240}
            render={(_, record: any) => (
              <Space orientation="vertical" size="small">
                <Text style={{ fontSize: 12 }}>Origin: <Text strong>{record.origin_zone?.name || "-"}</Text></Text>
                <Text style={{ fontSize: 12 }}>Dest: <Text strong>{record.dest_zone?.name || "-"}</Text></Text>
              </Space>
            )}
          />

          <Table.Column dataIndex={["mot", "name"]} title="MOT" width={120} />
          <Table.Column dataIndex="cost_idr" title="Cost (IDR)" width={150} render={val => <Text strong>{formatIDR(val)}</Text>} />
          <Table.Column dataIndex={["uom", "name"]} title="UoM" width={100} />
          <Table.Column dataIndex="payload" title="Payload" width={120} />

          <Table.Column dataIndex="cost_per_kg" title="Cost/KG" width={150} render={val => <Text>{formatIDR(val)}</Text>} />
          <Table.Column dataIndex="cost_per_ton" title="Cost/Ton" width={150} render={val => <Text>{formatIDR(val)}</Text>} />
          <Table.Column dataIndex="loading_cost" title="Loading Cost (IDR)" width={170} render={val => <Text>{formatIDR(val)}</Text>} />
          <Table.Column dataIndex="unloading_cost" title="Unloading Cost (IDR)" width={180} render={val => <Text>{formatIDR(val)}</Text>} />
          <Table.Column dataIndex="distance" title="Distance (KM)" width={120} />

          <Table.Column dataIndex="running_cost_idr" title="Running Cost (IDR/Ton/KM)" width={220} render={val => <Text strong>{formatIDR(val)}</Text>} />
          <Table.Column dataIndex="running_cost_usd" title="Running Cost (USD/Ton/KM)" width={220} />

          <Table.Column
            dataIndex="notes"
            title="Notes"
            width={250}
            render={(val) => val ? <Text type="secondary" italic>{val}</Text> : "-"}
          />

          <Table.Column
            title="Actions"
            dataIndex="actions"
            fixed="right"
            align="center"
            width={120}
            render={(_, record: any) => (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', justifyContent: 'center' }}>
                <EditButton size="middle" recordItemId={record.id} style={{ borderRadius: '4px' }} />
                <DeleteButton size="middle" recordItemId={record.id} style={{ borderRadius: '4px' }} />
              </div>
            )}
          />
        </Table>
      </Card>
    </div>
  );
}
