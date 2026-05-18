"use client";

import { useEffect, useMemo, useState } from "react";
import { useList } from "@refinedev/core";
import {
  App,
  Button,
  Card,
  Checkbox,
  DatePicker,
  Flex,
  Input,
  Popover,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  theme,
} from "antd";
import type { CheckboxOptionType, TableColumnsType } from "antd";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import { useSearchParams } from "next/navigation";
import {
  DownloadOutlined,
  ReloadOutlined,
  SearchOutlined,
  SettingOutlined,
} from "@ant-design/icons";
import ContractModeSwitch, { type ContractMode } from "@/components/contracts/ContractModeSwitch";

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;
const DATA_EXPLORER_DRAFT_STORAGE_KEY = "data_explorer_draft_v1";

type RelatedEntity = {
  id?: number | string;
  name?: string;
  code?: string;
};

type RawContract = {
  id?: number | string;
  spk_number?: string;
  fa_number?: string;
  vendor_id?: number | string;
  mill_id?: number | string;
  vendor?: RelatedEntity;
  mill?: RelatedEntity;
  mot?: RelatedEntity;
  origin_zone?: RelatedEntity;
  dest_zone?: RelatedEntity;
  distributed_cost?: number | string;
  cost_idr?: number | string;
  validity_start?: string;
  validity_end?: string;
};

type ExplorerRow = {
  key: string;
  contractType: "Oncall" | "Dedicated Fix" | "Dedicated Var";
  spkNumber: string;
  faNumber: string;
  vendorId?: number | string;
  vendorName: string;
  vendorCode: string;
  millId?: number | string;
  millName: string;
  route: string;
  mot: string;
  cost: number | null;
  validityStart: string | null;
  validityEnd: string | null;
};

type ColumnDefinition = {
  key: string;
  label: string;
  column: TableColumnsType<ExplorerRow>[number];
};

type DataExplorerDraft = {
  selectedMills: Array<string | number>;
  selectedVendors: Array<string | number>;
  selectedContractTypes: ExplorerRow["contractType"][];
  validityStart: string | null;
  validityEnd: string | null;
  searchText: string;
  visibleColumns: string[];
  contractMode: ContractMode;
};

const formatIDR = (value: number | null) => {
  if (value === null || Number.isNaN(value)) return "-";
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);
};

const formatDate = (value: string | null) => {
  if (!value) return "-";
  const parsed = dayjs(value);
  return parsed.isValid() ? parsed.format("DD/MM/YYYY") : "-";
};

const toNumber = (value: number | string | undefined) => {
  if (value === undefined || value === null || value === "") return null;
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : null;
};

const DEFAULT_VISIBLE_COLUMNS = [
  "contractType",
  "spkNumber",
  "vendorName",
  "millName",
  "route",
  "mot",
  "cost",
  "validity",
];

export default function DataExplorerPage() {
  const { token } = theme.useToken();
  const { message } = App.useApp();
  const searchParams = useSearchParams();
  const searchQuery = searchParams.get("q") || "";

  const contractTypeOptions: Array<{ label: string; value: ExplorerRow["contractType"] }> = [
    { label: "Oncall", value: "Oncall" },
    { label: "Dedicated Fix", value: "Dedicated Fix" },
    { label: "Dedicated Var", value: "Dedicated Var" },
  ];

  const [selectedMills, setSelectedMills] = useState<Array<string | number>>([]);
  const [selectedVendors, setSelectedVendors] = useState<Array<string | number>>([]);
  const [selectedContractTypes, setSelectedContractTypes] = useState<ExplorerRow["contractType"][]>([]);
  const [validityRange, setValidityRange] = useState<[Dayjs, Dayjs] | null>(null);
  const [searchText, setSearchText] = useState(searchQuery);
  const [visibleColumns, setVisibleColumns] = useState<string[]>(DEFAULT_VISIBLE_COLUMNS);
  const [contractMode, setContractMode] = useState<ContractMode>("plan");

  useEffect(() => {
    if (searchQuery.trim()) {
      setSearchText(searchQuery);
    }
  }, [searchQuery]);

  useEffect(() => {
    if (typeof window === "undefined") return;
    const rawDraft = window.sessionStorage.getItem(DATA_EXPLORER_DRAFT_STORAGE_KEY);
    if (!rawDraft) return;

    try {
      const draft = JSON.parse(rawDraft) as DataExplorerDraft;
      setSelectedMills(Array.isArray(draft.selectedMills) ? draft.selectedMills : []);
      setSelectedVendors(Array.isArray(draft.selectedVendors) ? draft.selectedVendors : []);
      setSelectedContractTypes(
        Array.isArray(draft.selectedContractTypes) ? draft.selectedContractTypes : [],
      );
      setSearchText(draft.searchText || "");
      setVisibleColumns(
        Array.isArray(draft.visibleColumns) && draft.visibleColumns.length > 0
          ? draft.visibleColumns
          : DEFAULT_VISIBLE_COLUMNS,
      );
      setContractMode(draft.contractMode === "actual" ? "actual" : "plan");

      if (draft.validityStart && draft.validityEnd) {
        const start = dayjs(draft.validityStart);
        const end = dayjs(draft.validityEnd);
        if (start.isValid() && end.isValid()) {
          setValidityRange([start, end]);
        }
      }
    } catch {
      window.sessionStorage.removeItem(DATA_EXPLORER_DRAFT_STORAGE_KEY);
    }
  }, []);

  useEffect(() => {
    if (typeof window === "undefined") return;
    const draft: DataExplorerDraft = {
      selectedMills,
      selectedVendors,
      selectedContractTypes,
      validityStart: validityRange?.[0]?.toISOString() || null,
      validityEnd: validityRange?.[1]?.toISOString() || null,
      searchText,
      visibleColumns,
      contractMode,
    };
    window.sessionStorage.setItem(DATA_EXPLORER_DRAFT_STORAGE_KEY, JSON.stringify(draft));
  }, [
    selectedMills,
    selectedVendors,
    selectedContractTypes,
    validityRange,
    searchText,
    visibleColumns,
    contractMode,
  ]);

  const { query: oncallQuery } = useList<RawContract>({
    resource: "contracts/oncall",
    pagination: { mode: "off" },
    meta: { populate: ["vendor", "mill", "origin_zone", "dest_zone", "mot"] },
  });

  const { query: dedicatedFixQuery } = useList<RawContract>({
    resource: "contracts/dedicated-fix",
    pagination: { mode: "off" },
    meta: { populate: ["vendor", "mill", "mot"] },
  });

  const { query: dedicatedVarQuery } = useList<RawContract>({
    resource: "contracts/dedicated-var",
    pagination: { mode: "off" },
    meta: { populate: ["vendor", "mill", "origin_zone", "dest_zone", "mot"] },
  });

  const { query: millsQuery } = useList<RelatedEntity>({
    resource: "mills",
    pagination: { mode: "off" },
    sorters: [{ field: "name", order: "asc" }],
  });

  const { query: vendorsQuery } = useList<RelatedEntity>({
    resource: "vendors",
    pagination: { mode: "off" },
    sorters: [{ field: "name", order: "asc" }],
  });

  const allRows = useMemo<ExplorerRow[]>(() => {
    const normalize = (
      item: RawContract,
      contractType: ExplorerRow["contractType"],
      index: number,
    ): ExplorerRow => {
      const vendorName = item.vendor?.name ?? "-";
      const originName = item.origin_zone?.name ?? "-";
      const destinationName = item.dest_zone?.name ?? "-";
      const contractTypeKey = contractType.toLowerCase().replace(/\s+/g, "-");

      return {
        key: `${contractTypeKey}-${item.id ?? item.spk_number ?? "row"}-${index}`,
        contractType,
        spkNumber: item.spk_number ?? "-",
        faNumber: item.fa_number ?? "-",
        vendorId: item.vendor_id ?? item.vendor?.id,
        vendorName,
        vendorCode: item.vendor?.code ?? "-",
        millId: item.mill_id ?? item.mill?.id,
        millName: item.mill?.name ?? item.mill?.code ?? "-",
        route:
          item.origin_zone || item.dest_zone
            ? `${originName} - ${destinationName}`
            : "-",
        mot: item.mot?.name ?? "-",
        cost: toNumber(item.distributed_cost) ?? toNumber(item.cost_idr),
        validityStart: item.validity_start ?? null,
        validityEnd: item.validity_end ?? null,
      };
    };

    const oncallRows = (oncallQuery.data?.data ?? []).map((item, index) =>
      normalize(item, "Oncall", index),
    );
    const dedicatedFixRows = (dedicatedFixQuery.data?.data ?? []).map((item, index) =>
      normalize(item, "Dedicated Fix", index),
    );
    const dedicatedVarRows = (dedicatedVarQuery.data?.data ?? []).map((item, index) =>
      normalize(item, "Dedicated Var", index),
    );

    return [...dedicatedFixRows, ...dedicatedVarRows, ...oncallRows];
  }, [oncallQuery.data?.data, dedicatedFixQuery.data?.data, dedicatedVarQuery.data?.data]);

  const filteredRows = useMemo(() => {
    const normalizedSearch = searchText.trim().toLowerCase();
    const selectedMillSet = new Set(selectedMills.map(String));
    const selectedVendorSet = new Set(selectedVendors.map(String));
    const selectedContractTypeSet = new Set(selectedContractTypes.map(String));

    return allRows.filter((row) => {
      if (
        selectedContractTypeSet.size > 0 &&
        !selectedContractTypeSet.has(String(row.contractType))
      ) {
        return false;
      }

      if (selectedMillSet.size > 0 && !selectedMillSet.has(String(row.millId ?? ""))) {
        return false;
      }

      if (selectedVendorSet.size > 0 && !selectedVendorSet.has(String(row.vendorId ?? ""))) {
        return false;
      }

      if (normalizedSearch) {
        const searchableValues = [
          row.contractType,
          row.spkNumber,
          row.faNumber,
          row.vendorName,
          row.vendorCode,
          row.millName,
          row.route,
          row.mot,
          row.cost !== null ? String(row.cost) : "",
          formatIDR(row.cost),
          formatDate(row.validityStart),
          formatDate(row.validityEnd),
        ];
        const searchIndex = searchableValues.join(" ").toLowerCase();
        const isSearchMatch = searchIndex.includes(normalizedSearch);
        if (!isSearchMatch) return false;
      }

      if (validityRange) {
        const [fromDate, toDate] = validityRange;
        const filterStart = fromDate.startOf("day").valueOf();
        const filterEnd = toDate.endOf("day").valueOf();
        const rowStart = dayjs(row.validityStart);
        const rowEnd = dayjs(row.validityEnd);
        const rangeStart = (rowStart.isValid() ? rowStart : rowEnd).valueOf();
        const rangeEnd = (rowEnd.isValid() ? rowEnd : rowStart).valueOf();

        if (!Number.isFinite(rangeStart) || !Number.isFinite(rangeEnd)) {
          return false;
        }

        const intersectsDateRange = rangeStart <= filterEnd && rangeEnd >= filterStart;
        if (!intersectsDateRange) return false;
      }

      return true;
    });
  }, [allRows, selectedMills, selectedVendors, selectedContractTypes, searchText, validityRange]);

  const millOptions = useMemo(
    () =>
      (millsQuery.data?.data ?? [])
        .filter((mill) => mill.id !== undefined)
        .map((mill) => ({
          label: mill.name ?? mill.code ?? `Mill ${mill.id}`,
          value: mill.id as string | number,
        })),
    [millsQuery.data?.data],
  );

  const vendorOptions = useMemo(
    () =>
      (vendorsQuery.data?.data ?? [])
        .filter((vendor) => vendor.id !== undefined)
        .map((vendor) => ({
          label: vendor.name ?? vendor.code ?? `Vendor ${vendor.id}`,
          value: vendor.id as string | number,
        })),
    [vendorsQuery.data?.data],
  );

  const allColumnDefinitions = useMemo<ColumnDefinition[]>(() => {
    const typeColorMap: Record<ExplorerRow["contractType"], string> = {
      Oncall: "cyan",
      "Dedicated Fix": "geekblue",
      "Dedicated Var": "purple",
    };

    return [
      {
        key: "contractType",
        label: "Contract Type",
        column: {
          title: "Contract Type",
          dataIndex: "contractType",
          key: "contractType",
          width: 160,
          render: (value: ExplorerRow["contractType"]) => (
            <Tag color={typeColorMap[value]}>{value}</Tag>
          ),
        },
      },
      {
        key: "spkNumber",
        label: "SPK Number",
        column: {
          title: "SPK Number",
          dataIndex: "spkNumber",
          key: "spkNumber",
          width: 180,
          render: (value: string) => <Text strong>{value}</Text>,
        },
      },
      {
        key: "faNumber",
        label: "FA Number",
        column: {
          title: "FA Number",
          dataIndex: "faNumber",
          key: "faNumber",
          width: 160,
        },
      },
      {
        key: "vendorName",
        label: "Vendor",
        column: {
          title: "Vendor",
          dataIndex: "vendorName",
          key: "vendorName",
          width: 220,
        },
      },
      {
        key: "vendorCode",
        label: "Vendor Code",
        column: {
          title: "Vendor Code",
          dataIndex: "vendorCode",
          key: "vendorCode",
          width: 140,
        },
      },
      {
        key: "millName",
        label: "Mill",
        column: {
          title: "Mill",
          dataIndex: "millName",
          key: "millName",
          width: 170,
        },
      },
      {
        key: "route",
        label: "Route",
        column: {
          title: "Route",
          dataIndex: "route",
          key: "route",
          width: 280,
        },
      },
      {
        key: "mot",
        label: "MOT",
        column: {
          title: "MOT",
          dataIndex: "mot",
          key: "mot",
          width: 120,
        },
      },
      {
        key: "cost",
        label: "Cost",
        column: {
          title: "Cost",
          dataIndex: "cost",
          key: "cost",
          width: 170,
          align: "right",
          render: (value: number | null) => <Text strong>{formatIDR(value)}</Text>,
        },
      },
      {
        key: "validity",
        label: "Validity",
        column: {
          title: "Validity",
          key: "validity",
          width: 240,
          render: (_, record) => (
            <Text type="secondary">
              {formatDate(record.validityStart)} - {formatDate(record.validityEnd)}
            </Text>
          ),
        },
      },
    ];
  }, []);

  const visibleTableColumns = useMemo<TableColumnsType<ExplorerRow>>(
    () =>
      allColumnDefinitions
        .filter((columnDef) => visibleColumns.includes(columnDef.key))
        .map((columnDef) => columnDef.column),
    [allColumnDefinitions, visibleColumns],
  );

  const handleResetFilters = () => {
    setSelectedContractTypes([]);
    setSelectedMills([]);
    setSelectedVendors([]);
    setValidityRange(null);
    setSearchText("");
  };

  const csvEscape = (value: unknown) => {
    const raw = value === null || value === undefined ? "" : String(value);
    return `"${raw.replace(/"/g, '""')}"`;
  };

  const handleExport = () => {
    if (filteredRows.length === 0) {
      message.warning("No data to export.");
      return;
    }

    const headers = [
      "Contract Type",
      "SPK Number",
      "FA Number",
      "Vendor",
      "Vendor Code",
      "Mill",
      "Route",
      "MOT",
      "Cost",
      "Validity Start",
      "Validity End",
    ];

    const lines = filteredRows.map((row) =>
      [
        row.contractType,
        row.spkNumber,
        row.faNumber,
        row.vendorName,
        row.vendorCode,
        row.millName,
        row.route,
        row.mot,
        row.cost ?? "",
        formatDate(row.validityStart),
        formatDate(row.validityEnd),
      ]
        .map(csvEscape)
        .join(","),
    );

    const csvContent = [headers.map(csvEscape).join(","), ...lines].join("\n");
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = `data-explorer-${dayjs().format("YYYYMMDD_HHmmss")}.csv`;
    document.body.appendChild(anchor);
    anchor.click();
    document.body.removeChild(anchor);
    URL.revokeObjectURL(url);

    message.success(`${filteredRows.length} rows exported.`);
  };

  const isTableLoading =
    oncallQuery.isLoading || dedicatedFixQuery.isLoading || dedicatedVarQuery.isLoading;

  const columnOptions: CheckboxOptionType<string>[] = allColumnDefinitions.map((columnDef) => ({
    label: columnDef.label,
    value: columnDef.key,
  }));

  return (
    <div className="dashboard-page-wide">
      <Flex justify="space-between" align="center" wrap="wrap" gap={12} style={{ marginBottom: 24 }}>
        <div>
          <Title level={2} style={{ margin: 0, fontWeight: 700 }}>
            Data Explorer
          </Title>
        </div>
        <ContractModeSwitch value={contractMode} actualEnabled={false} onChange={setContractMode} />
      </Flex>

      <Card
        styles={{ body: { padding: 20 } }}
        style={{
          borderRadius: 12,
          border: `1px solid ${token.colorBorderSecondary}`,
          boxShadow: token.boxShadowTertiary,
          marginBottom: 20,
        }}
      >
        <Flex vertical gap={14}>
          <Flex justify="space-between" align="center" wrap="wrap" gap={8}>
            <Text strong>Filters & Tools</Text>
            <Text type="secondary">
              {filteredRows.length} of {allRows.length} records
            </Text>
          </Flex>

          <Flex wrap="wrap" align="center" gap={10}>
            <Select
              size="large"
              mode="multiple"
              allowClear
              placeholder="Filter Contract Type"
              style={{ flex: "1 1 220px", minWidth: 220 }}
              value={selectedContractTypes}
              options={contractTypeOptions}
              onChange={(values) => setSelectedContractTypes(values)}
              maxTagCount="responsive"
            />

            <Select
              size="large"
              mode="multiple"
              allowClear
              placeholder="Filter Mills"
              style={{ flex: "1 1 220px", minWidth: 220 }}
              value={selectedMills}
              options={millOptions}
              onChange={(values) => setSelectedMills(values)}
              loading={millsQuery.isLoading}
              maxTagCount="responsive"
            />

            <Select
              size="large"
              mode="multiple"
              allowClear
              placeholder="Filter Vendors"
              style={{ flex: "1 1 240px", minWidth: 240 }}
              value={selectedVendors}
              options={vendorOptions}
              onChange={(values) => setSelectedVendors(values)}
              loading={vendorsQuery.isLoading}
              maxTagCount="responsive"
            />

            <RangePicker
              size="large"
              style={{ flex: "1 1 280px", minWidth: 280 }}
              value={validityRange}
              onChange={(dates) => {
                if (dates?.[0] && dates?.[1]) {
                  setValidityRange([dates[0], dates[1]]);
                  return;
                }
                setValidityRange(null);
              }}
            />

            <Button size="large" icon={<ReloadOutlined />} onClick={handleResetFilters}>
              Reset Filter
            </Button>
          </Flex>

          <Flex
            wrap="wrap"
            justify="space-between"
            align="center"
            gap={10}
            style={{ paddingTop: 14, borderTop: `1px solid ${token.colorSplit}` }}
          >
            <Input
              size="large"
              allowClear
              value={searchText}
              onChange={(event) => setSearchText(event.target.value)}
              prefix={<SearchOutlined style={{ color: token.colorTextTertiary }} />}
              placeholder="Search all data..."
              style={{ flex: "1 1 360px", minWidth: 280, maxWidth: "100%" }}
            />

            <Space size={8}>
              <Popover
                trigger="click"
                title="Choose Visible Columns"
                content={
                  <Checkbox.Group
                    style={{ display: "grid", gap: 8, maxWidth: 260 }}
                    options={columnOptions}
                    value={visibleColumns}
                    onChange={(checkedValues) => {
                      const nextValues = checkedValues as string[];
                      if (nextValues.length > 0) {
                        setVisibleColumns(nextValues);
                      }
                    }}
                  />
                }
              >
                <Button size="large" icon={<SettingOutlined />}>
                  Columns
                </Button>
              </Popover>

              <Button size="large" type="primary" icon={<DownloadOutlined />} onClick={handleExport}>
                Export to Excel
              </Button>
            </Space>
          </Flex>
        </Flex>
      </Card>

      <Card
        style={{
          borderRadius: 12,
          border: `1px solid ${token.colorBorderSecondary}`,
          boxShadow: token.boxShadowTertiary,
        }}
      >
        <Table<ExplorerRow>
          rowKey="key"
          loading={isTableLoading}
          dataSource={filteredRows}
          columns={visibleTableColumns}
          scroll={{ x: "max-content" }}
          pagination={{
            showSizeChanger: true,
            pageSizeOptions: ["10", "20", "50", "100"],
            defaultPageSize: 10,
            showTotal: (total) => `${total} records`,
          }}
        />
      </Card>
    </div>
  );
}
