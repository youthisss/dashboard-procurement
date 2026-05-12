"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { App, Button, Card, Form, Input, Modal, Popconfirm, Space, Table, Typography } from "antd";
import { Breadcrumb } from "@refinedev/antd";
import { DeleteOutlined, LockOutlined, UserAddOutlined } from "@ant-design/icons";
import { apiClient } from "@/lib/api-client";
import { getAuthUser } from "@/lib/auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
const { Title, Text } = Typography;

type CreateUserPayload = {
  name: string;
  username: string;
  password: string;
};

type UserRow = {
  id: number;
  name: string;
  username: string;
  created_at?: string;
};

export default function AdminUsersPage() {
  const router = useRouter();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [ready, setReady] = useState(false);
  const [users, setUsers] = useState<UserRow[]>([]);
  const [fetchingUsers, setFetchingUsers] = useState(false);
  const [currentUserId, setCurrentUserId] = useState<number | null>(null);
  const [passwordModalOpen, setPasswordModalOpen] = useState(false);
  const [updatingPassword, setUpdatingPassword] = useState(false);
  const [passwordTargetUser, setPasswordTargetUser] = useState<UserRow | null>(null);
  const [passwordForm] = Form.useForm();

  useEffect(() => {
    const authUser = getAuthUser();
    if (!authUser) {
      router.replace("/login");
      return;
    }
    if (!authUser.is_admin) {
      router.replace("/overview");
      return;
    }
    setCurrentUserId(authUser.id);
    setReady(true);
  }, [router]);

  const fetchUsers = async () => {
    setFetchingUsers(true);
    try {
      const { data } = await apiClient.get(`${API_URL}/auth/users`);
      setUsers(Array.isArray(data) ? data : []);
    } catch (error: any) {
      const err = error?.response?.data?.error || "Failed to fetch users";
      message.error(err);
    } finally {
      setFetchingUsers(false);
    }
  };

  useEffect(() => {
    if (ready) {
      void fetchUsers();
    }
  }, [ready]);

  const onFinish = async (values: CreateUserPayload) => {
    setLoading(true);
    try {
      await apiClient.post(`${API_URL}/auth/users`, values);
      message.success("User created successfully");
      void fetchUsers();
    } catch (error: any) {
      const err = error?.response?.data?.error || "Failed to create user";
      message.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (!ready) return null;

  const handleDeleteUser = async (userId: number) => {
    try {
      await apiClient.delete(`${API_URL}/auth/users/${userId}`);
      message.success("User deleted successfully");
      void fetchUsers();
    } catch (error: any) {
      const err = error?.response?.data?.error || "Failed to delete user";
      message.error(err);
    }
  };

  const openPasswordModal = (user: UserRow) => {
    setPasswordTargetUser(user);
    passwordForm.resetFields();
    setPasswordModalOpen(true);
  };

  const handleUpdatePassword = async () => {
    if (!passwordTargetUser) return;
    try {
      const values = await passwordForm.validateFields();
      setUpdatingPassword(true);
      await apiClient.patch(`${API_URL}/auth/users/${passwordTargetUser.id}/password`, {
        password: values.password,
      });
      message.success(`Password updated for ${passwordTargetUser.username}`);
      setPasswordModalOpen(false);
      passwordForm.resetFields();
    } catch (error: any) {
      if (error?.errorFields) return;
      const err = error?.response?.data?.error || "Failed to update password";
      message.error(err);
    } finally {
      setUpdatingPassword(false);
    }
  };

  return (
    <div className="dashboard-page">
      <div style={{ marginBottom: 16 }}>
        <Breadcrumb />
      </div>

      <Space direction="vertical" size={24} style={{ width: "100%" }}>
        <div>
          <Title level={2} style={{ margin: "0 0 8px 0", fontWeight: 700 }}>
            Admin - User Management
          </Title>
        </div>

        <Card
          variant="borderless"
          className="dashboard-table-card"
        >
          <Form layout="vertical" onFinish={onFinish} requiredMark={false}>
            <Form.Item label="Full Name" name="name" rules={[{ required: true, message: "Please input name" }]}>
              <Input placeholder="e.g. John Doe" />
            </Form.Item>
            <Form.Item label="Username" name="username" rules={[{ required: true, message: "Please input username" }]}>
              <Input placeholder="e.g. john.doe" />
            </Form.Item>
            <Form.Item
              label="Password"
              name="password"
              rules={[
                { required: true, message: "Please input password" },
                { min: 6, message: "Password must be at least 6 characters" },
              ]}
            >
              <Input.Password placeholder="Minimum 6 characters" />
            </Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} icon={<UserAddOutlined />}>
              Create User
            </Button>
          </Form>
        </Card>

        <Card
          variant="borderless"
          className="dashboard-table-card"
        >
          <Table
            rowKey="id"
            loading={fetchingUsers}
            dataSource={users}
            pagination={{ pageSize: 10 }}
            columns={[
              { title: "ID", dataIndex: "id", key: "id", width: 80 },
              { title: "Name", dataIndex: "name", key: "name" },
              { title: "Username", dataIndex: "username", key: "username", width: 220 },
              {
                title: "Action",
                key: "action",
                width: 260,
                render: (_, record: UserRow) => (
                  <Space>
                    <Button type="text" icon={<LockOutlined />} onClick={() => openPasswordModal(record)}>
                      Edit Password
                    </Button>
                    <Popconfirm
                      title="Delete user?"
                      description={`User "${record.username}" will be removed permanently.`}
                      onConfirm={() => handleDeleteUser(record.id)}
                      okText="Delete"
                      cancelText="Cancel"
                      disabled={record.id === currentUserId}
                    >
                      <Button
                        danger
                        type="text"
                        icon={<DeleteOutlined />}
                        disabled={record.id === currentUserId}
                      >
                        Delete
                      </Button>
                    </Popconfirm>
                  </Space>
                ),
              },
            ]}
          />
        </Card>
      </Space>

      <Modal
        title={`Edit Password - ${passwordTargetUser?.username || ""}`}
        open={passwordModalOpen}
        onCancel={() => setPasswordModalOpen(false)}
        onOk={handleUpdatePassword}
        okText="Update Password"
        confirmLoading={updatingPassword}
      >
        <Form form={passwordForm} layout="vertical" requiredMark={false}>
          <Form.Item
            label="New Password"
            name="password"
            rules={[
              { required: true, message: "Please input new password" },
              { min: 6, message: "Password must be at least 6 characters" },
            ]}
          >
            <Input.Password placeholder="Minimum 6 characters" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
