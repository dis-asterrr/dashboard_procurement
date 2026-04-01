"use client";

import React, { useCallback, useEffect, useRef, useState } from "react";
import { AutoComplete, Input, Tag, Typography, theme } from "antd";
import {
  BankOutlined,
  EnvironmentOutlined,
  FileTextOutlined,
  LoadingOutlined,
  SearchOutlined,
  ShopOutlined,
} from "@ant-design/icons";
import { useNavigation } from "@refinedev/core";
import { apiClient } from "@/lib/api-client";

const debounce = require("lodash/debounce");

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

type SearchApiItem = {
  id: number | string;
  type: string;
  label?: string;
  name?: string;
};

type SearchOption = {
  value: string;
  type: string;
  id: number | string;
  label: React.ReactNode;
};

const getTypeIcon = (type: string) => {
  const normalized = type.toLowerCase();

  if (normalized === "vendors") return <ShopOutlined style={{ color: "#1677ff" }} />;
  if (normalized === "mills") return <BankOutlined style={{ color: "#52c41a" }} />;
  if (normalized === "zones") return <EnvironmentOutlined style={{ color: "#faad14" }} />;
  if (normalized === "contracts" || normalized.includes("contract")) return <FileTextOutlined style={{ color: "#eb2f96" }} />;

  return <SearchOutlined />;
};

const normalizeText = (value: string) => value.trim();

export function GlobalSearchBar() {
  const { token } = theme.useToken();
  const { show } = useNavigation();
  const [loading, setLoading] = useState(false);
  const [options, setOptions] = useState<SearchOption[]>([]);
  const [inputValue, setInputValue] = useState("");

  const mountedRef = useRef(true);
  const abortControllerRef = useRef<AbortController | null>(null);
  const activeRequestIdRef = useRef(0);

  const runSearch = useCallback(async (searchText: string) => {
    const q = normalizeText(searchText);

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

    try {
      const response = await apiClient.get(`${API_URL}/search`, {
        params: { q },
        signal: controller.signal,
      });

      const raw = response.data?.data ?? response.data;
      const items: SearchApiItem[] = Array.isArray(raw) ? raw : [];

      const mappedOptions: SearchOption[] = items
        .filter((item) => item?.type && item?.id !== undefined && item?.id !== null)
        .map((item) => {
          const displayLabel = item.label || item.name || `${item.type} #${item.id}`;
          const itemType = String(item.type);
          const tagType = itemType.startsWith("contracts/") ? "contracts" : itemType;

          return {
            value: `${itemType}-${item.id}`,
            type: itemType,
            id: item.id,
            label: (
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12 }}>
                <div style={{ display: "flex", alignItems: "center", gap: 8, minWidth: 0 }}>
                  {getTypeIcon(itemType)}
                  <Typography.Text style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                    {displayLabel}
                  </Typography.Text>
                </div>
                <Tag style={{ border: "none", fontSize: 10, marginInlineEnd: 0 }}>{tagType.toUpperCase()}</Tag>
              </div>
            ),
          };
        });

      if (mountedRef.current && activeRequestIdRef.current === requestId) {
        setOptions(mappedOptions);
      }
    } catch (error: any) {
      if (
        error?.name !== "CanceledError" &&
        error?.code !== "ERR_CANCELED" &&
        mountedRef.current &&
        activeRequestIdRef.current === requestId
      ) {
        setOptions([]);
      }
    } finally {
      if (mountedRef.current && activeRequestIdRef.current === requestId) {
        setLoading(false);
      }
    }
  }, []);

  const debouncedSearch = useCallback(
    debounce((text: string) => {
      void runSearch(text);
    }, 500),
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
      style={{ width: 400, maxWidth: "100%" }}
      popupMatchSelectWidth={350}
      onChange={(value) => {
        const next = typeof value === "string" ? value : "";
        setInputValue(next);
        debouncedSearch(next);
      }}
      onSelect={(_, option) => {
        const selected = option as SearchOption;
        if (selected?.type && selected?.id !== undefined && selected?.id !== null) {
          show(selected.type, selected.id);
          setInputValue("");
        }
      }}
    >
      <Input
        size="middle"
        placeholder="Search everything (Vendors, Mills, SPK...)"
        allowClear
        prefix={<SearchOutlined style={{ color: token.colorTextTertiary }} />}
        suffix={loading ? <LoadingOutlined spin style={{ color: token.colorTextTertiary }} /> : null}
        style={{
          borderRadius: 8,
          backgroundColor: token.colorFillAlter,
          border: `1px solid ${token.colorBorderSecondary}`,
        }}
      />
    </AutoComplete>
  );
}
