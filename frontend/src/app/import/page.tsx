"use client";

import { Upload, Button, Card, Typography, Tag, Alert, Space, Divider, App, Row, Col, Flex, Result, Table, Input } from "antd";
import { UploadOutlined, FileExcelOutlined, CheckCircleOutlined, SaveOutlined, TableOutlined } from "@ant-design/icons";
import { useState } from "react";
import { apiClient } from "@/lib/api-client";
import ContractModeSwitch, { type ContractMode } from "@/components/contracts/ContractModeSwitch";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
const { Title, Text } = Typography;
const PREVIEW_PAGE_SIZE = 5;

interface ParsedSheet {
  sheet_name: string;
  sheet_type: string;
  headers: string[];
  rows: Record<string, string>[];
}

interface ImportResponse {
  file_name: string;
  saved_as: string;
  sheets: ParsedSheet[];
  warnings: string[];
}

type PreviewTableRow = Record<string, string | number> & {
  __previewRowIndex: number;
};

export default function ImportWizardPage() {
  const [contractMode, setContractMode] = useState<ContractMode>("plan");
  const [fileList, setFileList] = useState<any[]>([]);
  const [uploading, setUploading] = useState(false);
  const [result, setResult] = useState<ImportResponse | null>(null);
  const [confirming, setConfirming] = useState(false);
  const [confirmed, setConfirmed] = useState(false);
  const [confirmResult, setConfirmResult] = useState<any>(null);
  const { message } = App.useApp();

  const parsedSheets = result?.sheets
    .map((sheet, originalIndex) => ({ ...sheet, originalIndex }))
    .filter((s) => s.sheet_type !== "unknown" && s.rows.length > 0) || [];
  const isActualMode = contractMode === "actual";

  const handleModeChange = (nextMode: ContractMode) => {
    setContractMode(nextMode);
    setResult(null);
    setConfirmed(false);
    setConfirmResult(null);
  };

  const handleUpload = async () => {
    if (fileList.length === 0) {
      message.error("Please select an Excel file first.");
      return;
    }

    const formData = new FormData();
    formData.append("file", fileList[0]);

    setUploading(true);
    setResult(null);
    setConfirmed(false);
    setConfirmResult(null);

    try {
      const response = await apiClient.post(`${API_URL}/import/excel`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      message.success("File uploaded and parsed successfully.");
      setResult(response.data);
    } catch (error: any) {
      message.error(error.response?.data?.error || "Upload failed");
    } finally {
      setUploading(false);
    }
  };

  const handleConfirm = async () => {
    if (!result) return;
    setConfirming(true);
    try {
      const response = await apiClient.post(`${API_URL}/import/confirm`, {
        saved_as: result.saved_as,
        contract_mode: contractMode,
        sheets: result.sheets,
      });
      message.success("Data saved to database successfully.");
      setConfirmed(true);
      setConfirmResult(response.data.result);
    } catch (error: any) {
      const validationResult = error.response?.data?.result;
      if (validationResult) {
        setConfirmed(true);
        setConfirmResult(validationResult);
        message.warning("Import completed with warnings. Review row errors.");
      } else {
        message.error(error.response?.data?.error || "Confirm import failed");
      }
    } finally {
      setConfirming(false);
    }
  };

  const handleReset = () => {
    setFileList([]);
    setResult(null);
    setConfirmed(false);
    setConfirmResult(null);
  };

  const handlePreviewCellChange = (sheetIndex: number, rowIndex: number, field: string, value: string) => {
    setResult((current) => {
      if (!current) return current;
      return {
        ...current,
        sheets: current.sheets.map((sheet, currentSheetIndex) => {
          if (currentSheetIndex !== sheetIndex) return sheet;
          return {
            ...sheet,
            rows: sheet.rows.map((row, currentRowIndex) => {
              if (currentRowIndex !== rowIndex) return row;
              return {
                ...row,
                [field]: value,
              };
            }),
          };
        }),
      };
    });
  };

  const buildPreviewRows = (rows: Record<string, string>[]): PreviewTableRow[] =>
    rows.map((row, rowIndex) => ({
      ...row,
      __previewRowIndex: rowIndex,
    }));

  const buildPreviewColumns = (headers: string[], sheetIndex: number) => {
    return headers
      .filter((header) => header.trim() !== "")
      .map((header) => ({
        title: header,
        dataIndex: header,
        key: header,
        width: 180,
        render: (value: string | number | undefined, record: PreviewTableRow) => (
          <Input
            size="small"
            value={value == null ? "" : String(value)}
            onChange={(event) => handlePreviewCellChange(sheetIndex, record.__previewRowIndex, header, event.target.value)}
            style={{ minWidth: 140 }}
          />
        ),
      }));
  };

  return (
    <div className="dashboard-page">
      <div className="dashboard-page-header">
        <div>
          <Title level={2} style={{ margin: "0 0 8px 0", fontWeight: 700 }}>
            Data Import Wizard
          </Title>
        </div>
        <ContractModeSwitch value={contractMode} actualEnabled onChange={handleModeChange} />
      </div>

      <Card
        title={<Space><FileExcelOutlined /> Upload Excel Template</Space>}
        extra={<Tag color={isActualMode ? "warning" : "processing"}>{isActualMode ? "Actual Contract" : "Plan Contract"}</Tag>}
        variant="borderless"
        className="dashboard-card"
      >
        <div style={{ marginBottom: 24 }}>
          <Text type="secondary" style={{ display: "block", marginBottom: 16 }}>
            Upload your Excel template containing the latest {isActualMode ? "actual" : "plan"} procurement pricing. The system will detect sheet types
            (Dedicated Fix, Dedicated Var, Oncall) and parse them into the required structure.
          </Text>
          <Upload
            beforeUpload={(file) => {
              setFileList([file]);
              setResult(null);
              setConfirmed(false);
              setConfirmResult(null);
              return false;
            }}
            fileList={fileList}
            onRemove={() => {
              setFileList([]);
              setResult(null);
              setConfirmed(false);
              setConfirmResult(null);
            }}
            maxCount={1}
            accept=".xlsx, .xls"
          >
            <Button icon={<UploadOutlined />}>Select Excel File</Button>
          </Upload>
        </div>

        <Button
          onClick={handleUpload}
          disabled={fileList.length === 0}
          loading={uploading}
          type="primary"
          size="large"
          className="dashboard-action-button"
          style={{ width: "100%" }}
        >
          {uploading ? "Parsing document..." : "Upload & Parse Data"}
        </Button>
      </Card>

      {result && (
        <Card
          variant="borderless"
          className="dashboard-card"
          style={{ marginTop: 24 }}
        >
          <Space direction="vertical" size={6} style={{ width: "100%" }}>
            <Title level={4} style={{ margin: 0 }}>Import Status Report</Title>
            <Text type="secondary">Processed file: <strong>{result.file_name}</strong></Text>
          </Space>

          {result.warnings?.length > 0 && (
            <Alert
              message="Data Validation Warnings"
              description={
                <ul style={{ margin: 0, paddingLeft: 18 }}>
                  {result.warnings.map((item, idx) => (
                    <li key={idx}>
                      <Text type="danger">{item}</Text>
                    </li>
                  ))}
                </ul>
              }
              type="warning"
              showIcon
              style={{ marginTop: 16, marginBottom: 16, borderRadius: "8px" }}
            />
          )}

          <div style={{ marginTop: 16 }}>
            {parsedSheets.length === 0 ? (
              <Alert title="No recognizable data rows found in the sheets." type="info" showIcon />
            ) : (
              <>
                <Row gutter={[16, 16]}>
                  {parsedSheets.map((sheet, index) => (
                    <Col xs={24} md={12} key={index}>
                      <Card style={{ borderRadius: "8px", border: "1px solid var(--ant-color-border-secondary)", boxShadow: "none" }}>
                        <Flex vertical style={{ width: "100%" }}>
                          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                            <Text strong>{sheet.sheet_name}</Text>
                            <Tag color="processing">{sheet.sheet_type.replace("_", " ").toUpperCase()}</Tag>
                          </div>
                          <Divider style={{ margin: "12px 0" }} />
                          <div>
                            <CheckCircleOutlined style={{ color: "#10b981", marginRight: 8 }} />
                            <Text>{sheet.rows.length} rows parsed successfully.</Text>
                          </div>
                          <div style={{ marginTop: 8 }}>
                            <Text type="secondary" style={{ fontSize: 12 }}>
                              {sheet.headers.length} columns detected and loaded.
                            </Text>
                          </div>
                        </Flex>
                      </Card>
                    </Col>
                  ))}
                </Row>

                <Divider style={{ marginTop: 24 }}>
                  <Space>
                    <TableOutlined />
                    Editable Preview Data
                  </Space>
                </Divider>
                <Alert
                  message="Review and edit parsed rows before saving"
                  description={`Changes made in the preview table are sent to the import confirmation step and become the data saved to the ${isActualMode ? "actual" : "plan"} contract tables.`}
                  type="info"
                  showIcon
                  style={{ marginBottom: 16, borderRadius: 8 }}
                />

                <Space direction="vertical" size={16} style={{ width: "100%" }}>
                  {parsedSheets.map((sheet, index) => (
                    <Card
                      key={`${sheet.sheet_name}-${index}`}
                      size="small"
                      title={
                        <Space wrap>
                          <Text strong>{sheet.sheet_name}</Text>
                          <Tag color="processing">{sheet.sheet_type.replace("_", " ").toUpperCase()}</Tag>
                          <Text type="secondary">{sheet.rows.length} rows</Text>
                        </Space>
                      }
                      style={{ borderRadius: 8, border: "1px solid var(--ant-color-border-secondary)", boxShadow: "none" }}
                    >
                      <Table
                        rowKey={(record) => `${sheet.sheet_name}-${record.__previewRowIndex}`}
                        size="small"
                        columns={buildPreviewColumns(sheet.headers, sheet.originalIndex)}
                        dataSource={buildPreviewRows(sheet.rows)}
                        scroll={{ x: "max-content" }}
                        pagination={{
                          pageSize: PREVIEW_PAGE_SIZE,
                          placement: ["bottomEnd"],
                          showSizeChanger: false,
                          showTotal: (total) => `${total} rows`,
                        }}
                      />
                    </Card>
                  ))}
                </Space>

                {!confirmed ? (
                  <div style={{ marginTop: 24, display: "flex", gap: 12 }}>
                    <Button
                      icon={<SaveOutlined />}
                      type="primary"
                      size="large"
                      loading={confirming}
                      onClick={handleConfirm}
                      style={{ flex: 1, borderRadius: "6px" }}
                    >
                      {confirming ? "Saving to database..." : "Confirm & Save to Database"}
                    </Button>
                    <Button size="large" onClick={handleReset} style={{ borderRadius: "6px" }}>
                      Reset
                    </Button>
                  </div>
                ) : (
                  <Result
                    status={confirmResult?.errors?.length > 0 ? "error" : "success"}
                    title={confirmResult?.errors?.length > 0 ? "Data imported with warnings" : "Data imported successfully"}
                    subTitle={
                      confirmResult ? (
                        <div>
                          <Text>Dedicated Fix: <strong>{confirmResult.dedicated_fix_inserted}</strong> inserted, <strong>{confirmResult.dedicated_fix_skipped_duplicates || 0}</strong> skipped</Text><br />
                          <Text>Dedicated Var: <strong>{confirmResult.dedicated_var_inserted}</strong> inserted, <strong>{confirmResult.dedicated_var_skipped_duplicates || 0}</strong> skipped</Text><br />
                          <Text>Oncall: <strong>{confirmResult.oncall_inserted}</strong> inserted, <strong>{confirmResult.oncall_skipped_duplicates || 0}</strong> skipped</Text>
                          {confirmResult.errors?.length > 0 && (
                            <Alert
                              message={`${confirmResult.errors.length} row warning/error(s); valid rows were still saved`}
                              description={
                                <div style={{ maxHeight: 200, overflow: "auto" }}>
                                  {confirmResult.errors.map((e: string, i: number) => (
                                    <div key={i}>- {e}</div>
                                  ))}
                                </div>
                              }
                              type="warning"
                              showIcon
                              style={{ marginTop: 16, textAlign: "left" }}
                            />
                          )}
                        </div>
                      ) : "All data has been saved."
                    }
                    extra={[
                      <Button key="new" onClick={handleReset}>Upload Another File</Button>,
                    ]}
                  />
                )}
              </>
            )}
          </div>
        </Card>
      )}
    </div>
  );
}
