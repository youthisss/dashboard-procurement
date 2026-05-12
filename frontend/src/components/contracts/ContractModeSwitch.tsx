"use client";

import { Segmented } from "antd";
import { FileDoneOutlined, FieldTimeOutlined } from "@ant-design/icons";
import type { ReactNode } from "react";

export type ContractMode = "plan" | "actual";

type ContractModeSwitchProps = {
  value?: ContractMode;
  actualEnabled?: boolean;
  onChange?: (mode: ContractMode) => void;
};

const optionLabel = (icon: ReactNode, label: string) => (
  <span style={{ display: "inline-flex", alignItems: "center", gap: 6 }}>
    {icon}
    {label}
  </span>
);

export default function ContractModeSwitch({
  value = "plan",
  actualEnabled = false,
  onChange,
}: ContractModeSwitchProps) {
  return (
    <Segmented
      aria-label="Contract mode"
      value={value}
      onChange={(next) => onChange?.(next as ContractMode)}
      options={[
        {
          label: optionLabel(<FileDoneOutlined />, "Plan Contract"),
          value: "plan",
        },
        {
          label: optionLabel(<FieldTimeOutlined />, "Actual Contract"),
          value: "actual",
          disabled: !actualEnabled,
        },
      ]}
    />
  );
}
