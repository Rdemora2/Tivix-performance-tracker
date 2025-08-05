const API_BASE_URL = `${
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"
}${import.meta.env.VITE_API_PREFIX || "/api/v1"}`;

export const tokenUtils = {
  get: () => localStorage.getItem("token"),
  set: (token) => localStorage.setItem("token", token),
  remove: () => localStorage.removeItem("token"),

  isValid: () => {
    const token = tokenUtils.get();
    if (!token) return false;

    try {
      const payload = JSON.parse(atob(token.split(".")[1]));
      return payload.exp > Date.now() / 1000;
    } catch {
      return false;
    }
  },

  getPayload: () => {
    const token = tokenUtils.get();
    if (!token) return null;

    try {
      return JSON.parse(atob(token.split(".")[1]));
    } catch {
      return null;
    }
  },
};

export const hasPermission = (action, resource) => {
  const payload = tokenUtils.getPayload();
  if (!payload || !payload.role) return false;

  const { role } = payload;
  if (role === "admin") return true;

  // Exclusão é permitida apenas para administradores
  if (action === "delete") return false;

  if (
    role === "manager" &&
    ["teams", "developers", "reports"].includes(resource)
  ) {
    return ["create", "update", "read"].includes(action);
  }

  return action === "read";
};

const apiRequest = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;

  const token = tokenUtils.get();
  const headers = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  if (token && tokenUtils.isValid()) {
    headers.Authorization = `Bearer ${token}`;
  }

  const config = {
    headers,
    ...options,
  };

  try {
    const response = await fetch(url, config);

    if (response.status === 401) {
      tokenUtils.remove();
      window.location.href = "/login";
      throw new Error("Sessão expirada. Faça login novamente.");
    }

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || `HTTP error! status: ${response.status}`);
    }

    return data;
  } catch (error) {
    console.error(`API request failed: ${endpoint}`, error);
    throw error;
  }
};

export const authAPI = {
  login: (credentials) =>
    apiRequest("/auth/login", {
      method: "POST",
      body: JSON.stringify(credentials),
    }),

  register: (userData) =>
    apiRequest("/auth/register", {
      method: "POST",
      body: JSON.stringify(userData),
    }),

  createUser: (userData) =>
    apiRequest("/auth/create-user", {
      method: "POST",
      body: JSON.stringify(userData),
    }),

  changePassword: (passwordData) =>
    apiRequest("/auth/change-password", {
      method: "POST",
      body: JSON.stringify(passwordData),
    }),
  // Definir nova senha no primeiro acesso
  setNewPassword: (passwordData) =>
    apiRequest("/auth/set-new-password", {
      method: "POST",
      body: JSON.stringify(passwordData),
    }),

  // Listar usuários (Admin apenas)
  getUsers: () => apiRequest("/auth/users"),

  profile: () => apiRequest("/auth/profile"),

  logout: () => {
    tokenUtils.remove();
    window.location.href = "/login";
  },
};

export const initAPI = {
  check: () => apiRequest("/init/check"),

  createAdmin: (adminData) =>
    apiRequest("/init/admin", {
      method: "POST",
      body: JSON.stringify(adminData),
    }),
};

export const teamsAPI = {
  getAll: () => apiRequest("/teams"),

  getById: (id) => apiRequest(`/teams/${id}`),

  create: (teamData) =>
    apiRequest("/teams", {
      method: "POST",
      body: JSON.stringify(teamData),
    }),

  update: (id, teamData) =>
    apiRequest(`/teams/${id}`, {
      method: "PUT",
      body: JSON.stringify(teamData),
    }),

  delete: (id) =>
    apiRequest(`/teams/${id}`, {
      method: "DELETE",
    }),

  getDevelopers: (teamId) => apiRequest(`/teams/${teamId}/developers`),
};

export const developersAPI = {
  getAll: (includeArchived = false) =>
    apiRequest(`/developers?includeArchived=${includeArchived}`),

  getArchived: () => apiRequest("/developers/archived"),

  getById: (id) => apiRequest(`/developers/${id}`),

  create: (developerData) =>
    apiRequest("/developers", {
      method: "POST",
      body: JSON.stringify(developerData),
    }),

  update: (id, developerData) =>
    apiRequest(`/developers/${id}`, {
      method: "PUT",
      body: JSON.stringify(developerData),
    }),

  archive: (id, archive = true) =>
    apiRequest(`/developers/${id}/archive`, {
      method: "PUT",
      body: JSON.stringify({ archive }),
    }),

  delete: (id) =>
    apiRequest(`/developers/${id}`, {
      method: "DELETE",
    }),

  getReports: (developerId) => apiRequest(`/developers/${developerId}/reports`),
};

export const performanceReportsAPI = {
  getAll: () => apiRequest("/performance-reports"),

  getById: (id) => apiRequest(`/performance-reports/${id}`),

  create: (reportData) =>
    apiRequest("/performance-reports", {
      method: "POST",
      body: JSON.stringify(reportData),
    }),

  getAvailableMonths: () => apiRequest("/performance-reports/months"),

  getByMonth: (month) => apiRequest(`/performance-reports/month/${month}`),

  getStats: () => apiRequest("/performance-reports/stats"),
};

export default {
  auth: authAPI,
  init: initAPI,
  teams: teamsAPI,
  developers: developersAPI,
  performanceReports: performanceReportsAPI,
};
