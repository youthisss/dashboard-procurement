"use client";

const USER_KEY = "auth_user";

export type AuthUser = {
  id: number;
  name: string;
  username: string;
  is_admin?: boolean;
};

export function getToken(): string | null {
  return null;
}

export function setAuth(_token: string, user: AuthUser) {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearAuth() {
  if (typeof window === "undefined") return;
  window.localStorage.removeItem(USER_KEY);
}

export function getAuthUser(): AuthUser | null {
  if (typeof window === "undefined") return null;
  const raw = window.localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as AuthUser;
  } catch {
    return null;
  }
}
