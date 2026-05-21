"use client";

import React, { useCallback, useEffect, useRef, useState } from "react";
import { AutoComplete, Input, Typography, Empty, Grid, theme } from "antd";
import {
  BankOutlined,
  EnvironmentOutlined,
  FileTextOutlined,
  SearchOutlined,
  ShopOutlined,
  RightOutlined,
} from "@ant-design/icons";
import { useRouter } from "next/navigation";
import { apiClient } from "@/lib/api-client";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

type DebouncedFunction<T extends (...args: any[]) => void> = ((...args: Parameters<T>) => void) & {
  cancel: () => void;
};

const debounce = <T extends (...args: any[]) => void>(fn: T, delay: number): DebouncedFunction<T> => {
  let timeout: ReturnType<typeof setTimeout> | undefined;
  const debounced = ((...args: Parameters<T>) => {
    if (timeout) {
      clearTimeout(timeout);
    }
    timeout = setTimeout(() => fn(...args), delay);
  }) as DebouncedFunction<T>;

  debounced.cancel = () => {
    if (timeout) {
      clearTimeout(timeout);
      timeout = undefined;
    }
  };

  return debounced;
};

type SearchApiItem = {
  id: number | string;
  type: string;
  label?: string;
  name?: string;
  target_path?: string;
};

type SearchOption = {
  value: string;
  type: string;
  id: number | string;
  displayLabel: string;
  targetPath?: string;
  label: React.ReactNode;
};

const TYPE_CONFIG: Record<string, { icon: React.ReactNode; color: string; tag: string; route: (id: number | string) => string }> = {
  vendors: {
    icon: <ShopOutlined />,
    color: "#1677ff",
    tag: "Vendor",
    route: (id) => `/vendors/edit/${id}`,
  },
  mills: {
    icon: <BankOutlined />,
    color: "#52c41a",
    tag: "Mill",
    route: (id) => `/mills/edit/${id}`,
  },
  zones: {
    icon: <EnvironmentOutlined />,
    color: "#faad14",
    tag: "Zone",
    route: (id) => `/zones/edit/${id}`,
  },
  "contracts/dedicated-fix": {
    icon: <FileTextOutlined />,
    color: "#eb2f96",
    tag: "Dedicated Fix",
    route: (id) => `/contracts/dedicated-fix/edit/${id}`,
  },
  "contracts/dedicated-fix/search": {
    icon: <SearchOutlined />,
    color: "#eb2f96",
    tag: "Dedicated Fix Results",
    route: () => "/contracts/dedicated-fix",
  },
  "contracts/dedicated-var": {
    icon: <FileTextOutlined />,
    color: "#722ed1",
    tag: "Dedicated Var",
    route: (id) => `/contracts/dedicated-var/edit/${id}`,
  },
  "contracts/dedicated-var/search": {
    icon: <SearchOutlined />,
    color: "#722ed1",
    tag: "Dedicated Var Results",
    route: () => "/contracts/dedicated-var",
  },
  "contracts/oncall": {
    icon: <FileTextOutlined />,
    color: "#13c2c2",
    tag: "Oncall",
    route: (id) => `/contracts/oncall/edit/${id}`,
  },
  "contracts/oncall/search": {
    icon: <SearchOutlined />,
    color: "#13c2c2",
    tag: "Oncall Results",
    route: () => "/contracts/oncall",
  },
};

const buildSearchPath = (path: string, query: string) => `${path}?q=${encodeURIComponent(query)}`;

const buildContractSearchItems = (query: string): SearchApiItem[] => [
  {
    id: "dedicated-var-search",
    type: "contracts/dedicated-var/search",
    label: `Search Dedicated Var contracts for "${query}"`,
    target_path: buildSearchPath("/contracts/dedicated-var", query),
  },
  {
    id: "oncall-search",
    type: "contracts/oncall/search",
    label: `Search Oncall contracts for "${query}"`,
    target_path: buildSearchPath("/contracts/oncall", query),
  },
  {
    id: "dedicated-fix-search",
    type: "contracts/dedicated-fix/search",
    label: `Search Dedicated Fix contracts for "${query}"`,
    target_path: buildSearchPath("/contracts/dedicated-fix", query),
  },
];

const getConfig = (type: string) => {
  return TYPE_CONFIG[type] || {
    icon: <SearchOutlined />,
    color: "#8c8c8c",
    tag: type,
    route: (id: number | string) => `/${type}/edit/${id}`,
  };
};

/** Highlight matching text in the label */
const highlightMatch = (text: string, query: string, tokenColor: string) => {
  if (!query || query.length < 2) return text;
  const idx = text.toLowerCase().indexOf(query.toLowerCase());
  if (idx === -1) return text;
  const before = text.slice(0, idx);
  const match = text.slice(idx, idx + query.length);
  const after = text.slice(idx + query.length);
  return (
    <>
      {before}
      <span style={{ color: tokenColor, fontWeight: 600 }}>{match}</span>
      {after}
    </>
  );
};

export function GlobalSearchBar() {
  const { token } = theme.useToken();
  const screens = Grid.useBreakpoint();
  const isMobile = !screens.md;
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [options, setOptions] = useState<any[]>([]);
  const [inputValue, setInputValue] = useState("");

  const mountedRef = useRef(true);
  const abortControllerRef = useRef<AbortController | null>(null);
  const activeRequestIdRef = useRef(0);

  const buildGroupedOptions = useCallback((items: SearchApiItem[], q: string) => {
    const grouped: Record<string, SearchApiItem[]> = {};
    for (const item of items) {
      if (!item?.type || item?.id === undefined || item?.id === null) continue;
      const key = item.type;
      if (!grouped[key]) grouped[key] = [];
      if (grouped[key].length < 5) grouped[key].push(item);
    }

    return Object.entries(grouped).map(([type, groupItems]) => {
      const config = getConfig(type);
      return {
        label: (
          <div style={{
            display: "flex",
            alignItems: "center",
            gap: 6,
            padding: "4px 0 2px",
            fontSize: 11,
            fontWeight: 600,
            color: token.colorTextSecondary,
            textTransform: "uppercase" as const,
            letterSpacing: "0.5px",
          }}>
            <span style={{ color: config.color, fontSize: 12 }}>{config.icon}</span>
            {config.tag}
          </div>
        ),
        options: groupItems.map((item) => {
          const displayLabel = item.label || item.name || `${item.type} #${item.id}`;
          return {
            value: `${item.type}-${item.id}`,
            type: item.type,
            id: item.id,
            displayLabel,
            targetPath: item.target_path,
            label: (
              <div
                style={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  padding: "4px 0",
                  gap: 8,
                }}
              >
                <Typography.Text
                  style={{
                    fontSize: 13,
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                    flex: 1,
                    minWidth: 0,
                  }}
                >
                  {highlightMatch(displayLabel, q, config.color)}
                </Typography.Text>
                <RightOutlined style={{ fontSize: 10, color: token.colorTextQuaternary, flexShrink: 0 }} />
              </div>
            ),
          };
        }),
      };
    });
  }, [token]);

  const runSearch = useCallback(async (searchText: string) => {
    const q = searchText.trim();

    if (q.length < 2) {
      setOptions([]);
      setLoading(false);
      abortControllerRef.current?.abort();
      abortControllerRef.current = null;
      return;
    }

    const requestId = activeRequestIdRef.current + 1;
    activeRequestIdRef.current = requestId;
    abortControllerRef.current?.abort();
    const controller = new AbortController();
    abortControllerRef.current = controller;
    setLoading(true);
    const contractSearchItems = buildContractSearchItems(q);
    setOptions(buildGroupedOptions(contractSearchItems, q));

    try {
      const response = await apiClient.get(`${API_URL}/search`, {
        params: { q },
        signal: controller.signal,
      });

      const raw = response.data?.data ?? response.data;
      const items: SearchApiItem[] = Array.isArray(raw) ? raw : [];
      const resultTypes = new Set(items.map((item) => item.type));
      const fallbackItems = contractSearchItems.filter((item) => !resultTypes.has(item.type));
      const groupedOptions = buildGroupedOptions([...fallbackItems, ...items], q);

      if (mountedRef.current && activeRequestIdRef.current === requestId) {
        setOptions(groupedOptions);
      }
    } catch (error: any) {
      if (
        error?.name !== "CanceledError" &&
        error?.code !== "ERR_CANCELED"
      ) {
        console.error("[GlobalSearch] Search failed:", error?.response?.status, error?.response?.data || error?.message);
      }
      if (
        error?.name !== "CanceledError" &&
        error?.code !== "ERR_CANCELED" &&
        mountedRef.current &&
        activeRequestIdRef.current === requestId
      ) {
        setOptions(buildGroupedOptions(contractSearchItems, q));
      }
    } finally {
      if (mountedRef.current && activeRequestIdRef.current === requestId) {
        setLoading(false);
      }
    }
  }, [buildGroupedOptions]);

  const debouncedSearch = useCallback(
    debounce((text: string) => {
      void runSearch(text);
    }, 350),
    [runSearch],
  );

  useEffect(() => {
    return () => {
      mountedRef.current = false;
      debouncedSearch.cancel();
      abortControllerRef.current?.abort();
    };
  }, [debouncedSearch]);

  return (
    <AutoComplete
      value={inputValue}
      options={options}
      filterOption={false}
      style={{ width: "min(420px, 100%)", maxWidth: "100%" }}
      popupMatchSelectWidth={isMobile ? false : 420}
      notFoundContent={
        inputValue.trim().length >= 2 && !loading ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={<Typography.Text type="secondary" style={{ fontSize: 12 }}>No results for &ldquo;{inputValue.trim()}&rdquo;</Typography.Text>}
            style={{ padding: "8px 0" }}
          />
        ) : null
      }
      onChange={(value) => {
        const next = typeof value === "string" ? value : "";
        const q = next.trim();
        setInputValue(next);
        if (q.length >= 2) {
          setOptions(buildGroupedOptions(buildContractSearchItems(q), q));
        } else {
          setOptions([]);
        }
        debouncedSearch(next);
      }}
      onSelect={(_, option) => {
        const selected = option as SearchOption;
        if (selected?.type && selected?.id !== undefined && selected?.id !== null) {
          const config = getConfig(selected.type);
          router.push(selected.targetPath || config.route(selected.id));
          setInputValue("");
          setOptions([]);
        }
      }}
      classNames={{ popup: { root: "global-search-popup" } }}
    >
      <Input
        size="middle"
        placeholder="Search vendors, mills, contracts..."
        allowClear
        onPressEnter={() => {
          const q = inputValue.trim();
          if (q.length >= 2) {
            router.push(buildSearchPath("/contracts/dedicated-var", q));
            setInputValue("");
            setOptions([]);
          }
        }}
        prefix={<SearchOutlined style={{ color: token.colorTextTertiary }} />}
        style={{
          borderRadius: 8,
          backgroundColor: token.colorFillAlter,
          border: `1px solid ${token.colorBorderSecondary}`,
          height: 36,
        }}
      />
    </AutoComplete>
  );
}
