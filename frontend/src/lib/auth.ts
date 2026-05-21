"use client";

const USER_KEY = "auth_user";
const TOKEN_KEY = "auth_token";

export type AuthUser = {
  id: number;
  name: string;
  username: string;
  is_admin?: boolean;
};

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function setAuth(token: string, user: AuthUser) {
  if (typeof window === "undefined") return;
  if (typeof token === "string" && token.trim() !== "") {
    window.localStorage.setItem(TOKEN_KEY, token);
  }
  window.localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearAuth() {
  if (typeof window === "undefined") return;
  window.localStorage.removeItem(TOKEN_KEY);
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
