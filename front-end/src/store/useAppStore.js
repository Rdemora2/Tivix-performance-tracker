import { create } from "zustand";
import { persist } from "zustand/middleware";
import api, { tokenUtils, hasPermission } from "../services/api";

const useAppStore = create(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,

      developers: [],
      archivedDevelopers: [],
      teams: [],
      performanceReports: [],
      users: [], // Lista de usuários (apenas para admin)
      darkMode: true,
      loading: false,
      error: null,

      login: async (credentials) => {
        try {
          set({ loading: true, error: null });
          const response = await api.auth.login(credentials);

          if (response.data?.token) {
            tokenUtils.set(response.data.token);
            set({
              user: response.data.user,
              isAuthenticated: true,
              loading: false,
            });
            return response.data;
          }
        } catch (error) {
          console.error("Login error:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      register: async (userData) => {
        try {
          set({ loading: true, error: null });
          const response = await api.auth.register(userData);

          if (response.data?.token) {
            tokenUtils.set(response.data.token);
            set({
              user: response.data.user,
              isAuthenticated: true,
              loading: false,
            });
            return response.data;
          }
        } catch (error) {
          console.error("Register error:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      // Admin cria usuário com senha temporária
      createUser: async (userData) => {
        try {
          set({ loading: true, error: null });
          const response = await api.auth.createUser(userData);
          set({ loading: false });
          return response.data;
        } catch (error) {
          console.error("Create user error:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      // Definir nova senha (primeiro acesso)
      setNewPassword: async (passwordData) => {
        try {
          set({ loading: true, error: null });
          const response = await api.auth.setNewPassword(passwordData);
          
          if (response.data?.token) {
            tokenUtils.set(response.data.token);
            set({
              user: response.data.user,
              isAuthenticated: true,
              loading: false,
            });
            return response.data;
          }
        } catch (error) {
          console.error("Set new password error:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      // Alterar senha
      changePassword: async (passwordData) => {
        try {
          set({ loading: true, error: null });
          const response = await api.auth.changePassword(passwordData);
          set({ loading: false });
          return response.data;
        } catch (error) {
          console.error("Change password error:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      // Buscar lista de usuários (Admin apenas)
      fetchUsers: async () => {
        try {
          set({ loading: true, error: null });
          const response = await api.auth.getUsers();
          set({ users: response.data || [], loading: false });
          return response.data;
        } catch (error) {
          console.error("Error fetching users:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      logout: () => {
        tokenUtils.remove();
        set({
          user: null,
          isAuthenticated: false,
          developers: [],
          archivedDevelopers: [],
          teams: [],
          performanceReports: [],
          users: [],
        });
      },

      loadUserProfile: async () => {
        try {
          if (!tokenUtils.isValid()) {
            get().logout();
            return;
          }

          const response = await api.auth.profile();
          set({
            user: response.data,
            isAuthenticated: true,
          });
        } catch (error) {
          console.error("Error loading user profile:", error);
          get().logout();
        }
      },

      hasPermission: (action, resource) => hasPermission(action, resource),

      setLoading: (loading) => set({ loading }),
      setError: (error) => set({ error }),

      fetchTeams: async () => {
        try {
          set({ loading: true, error: null });
          const response = await api.teams.getAll();
          set({ teams: response.data || [], loading: false });
        } catch (error) {
          console.error("Error fetching teams:", error);
          set({ error: error.message, loading: false });
        }
      },

      fetchDevelopers: async (includeArchived = false) => {
        try {
          set({ loading: true, error: null });
          const response = await api.developers.getAll(includeArchived);
          const developers = response.data || [];

          const active = developers.filter((dev) => !dev.archivedAt);
          const archived = developers.filter((dev) => dev.archivedAt);

          set({
            developers: active,
            archivedDevelopers: archived,
            loading: false,
          });
        } catch (error) {
          console.error("Error fetching developers:", error);
          set({ error: error.message, loading: false });
        }
      },

      fetchPerformanceReports: async () => {
        try {
          set({ loading: true, error: null });
          const response = await api.performanceReports.getAll();
          set({ performanceReports: response.data || [], loading: false });
        } catch (error) {
          console.error("Error fetching performance reports:", error);
          set({ error: error.message, loading: false });
        }
      },

      initializeStore: async () => {
        if (tokenUtils.isValid()) {
          await get().loadUserProfile();
        }

        if (get().isAuthenticated) {
          const { fetchTeams, fetchDevelopers, fetchPerformanceReports } =
            get();
          await Promise.all([
            fetchTeams(),
            fetchDevelopers(true),
            fetchPerformanceReports(),
          ]);
        }
      },

      addDeveloper: async (developer) => {
        try {
          set({ loading: true, error: null });
          const response = await api.developers.create(developer);
          const newDeveloper = response.data;

          set((state) => ({
            developers: [...state.developers, newDeveloper],
            loading: false,
          }));

          return newDeveloper;
        } catch (error) {
          console.error("Error adding developer:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      updateDeveloper: async (id, updates) => {
        try {
          set({ loading: true, error: null });
          const response = await api.developers.update(id, updates);
          const updatedDeveloper = response.data;

          set((state) => ({
            developers: state.developers.map((dev) =>
              dev.id === id ? updatedDeveloper : dev
            ),
            loading: false,
          }));

          return updatedDeveloper;
        } catch (error) {
          console.error("Error updating developer:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      archiveDeveloper: async (id) => {
        try {
          set({ loading: true, error: null });
          const response = await api.developers.archive(id, true);
          const archivedDeveloper = response.data;

          set((state) => ({
            developers: state.developers.filter((dev) => dev.id !== id),
            archivedDevelopers: [
              ...state.archivedDevelopers,
              archivedDeveloper,
            ],
            loading: false,
          }));

          return archivedDeveloper;
        } catch (error) {
          console.error("Error archiving developer:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      restoreDeveloper: async (id) => {
        try {
          set({ loading: true, error: null });
          const response = await api.developers.archive(id, false);
          const restoredDeveloper = response.data;

          set((state) => ({
            developers: [...state.developers, restoredDeveloper],
            archivedDevelopers: state.archivedDevelopers.filter(
              (dev) => dev.id !== id
            ),
            loading: false,
          }));

          return restoredDeveloper;
        } catch (error) {
          console.error("Error restoring developer:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      deleteDeveloper: async (id) => {
        try {
          set({ loading: true, error: null });
          const response = await api.developers.delete(id);
          const deletedDeveloper = response.data.deletedDeveloper;

          set((state) => ({
            developers: state.developers.filter((dev) => dev.id !== id),
            archivedDevelopers: state.archivedDevelopers.filter(
              (dev) => dev.id !== id
            ),
            loading: false,
          }));

          return deletedDeveloper;
        } catch (error) {
          console.error("Error deleting developer:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      // Team actions
      addTeam: async (team) => {
        try {
          set({ loading: true, error: null });
          const response = await api.teams.create(team);
          const newTeam = response.data;

          set((state) => ({
            teams: [...state.teams, newTeam],
            loading: false,
          }));

          return newTeam;
        } catch (error) {
          console.error("Error adding team:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      updateTeam: async (id, updates) => {
        try {
          set({ loading: true, error: null });
          const response = await api.teams.update(id, updates);
          const updatedTeam = response.data;

          set((state) => ({
            teams: state.teams.map((team) =>
              team.id === id ? updatedTeam : team
            ),
            loading: false,
          }));

          return updatedTeam;
        } catch (error) {
          console.error("Error updating team:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      removeTeam: async (id) => {
        try {
          set({ loading: true, error: null });
          await api.teams.delete(id);

          set((state) => ({
            teams: state.teams.filter((team) => team.id !== id),
            developers: state.developers.map((dev) =>
              dev.teamId === id ? { ...dev, teamId: null } : dev
            ),
            loading: false,
          }));
        } catch (error) {
          console.error("Error removing team:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      getDevelopersByTeam: (teamId) => {
        const state = get();
        return state.developers.filter((dev) => dev.teamId === teamId);
      },

      addPerformanceReport: async (report) => {
        try {
          set({ loading: true, error: null });
          const response = await api.performanceReports.create(report);
          const newReport = response.data;

          set((state) => {
            const newReports = [...state.performanceReports, newReport];

            const updatedDevelopers = state.developers.map((dev) =>
              dev.id === newReport.developerId
                ? {
                    ...dev,
                    latestPerformanceScore: newReport.weightedAverageScore,
                  }
                : dev
            );

            return {
              performanceReports: newReports,
              developers: updatedDevelopers,
              loading: false,
            };
          });

          return newReport;
        } catch (error) {
          console.error("Error adding performance report:", error);
          set({ error: error.message, loading: false });
          throw error;
        }
      },

      getReportsByDeveloper: (developerId) => {
        const state = get();
        return state.performanceReports
          .filter((report) => report.developerId === developerId)
          .sort((a, b) => new Date(b.month) - new Date(a.month));
      },

      getLatestReport: (developerId) => {
        const reports = get().getReportsByDeveloper(developerId);
        return reports[0] || null;
      },

      toggleDarkMode: () =>
        set((state) => ({
          darkMode: !state.darkMode,
        })),

      setDarkMode: (darkMode) =>
        set(() => ({
          darkMode,
        })),
    }),
    {
      name: "tivix-performance-tracker-storage",
      partialize: (state) => ({
        darkMode: state.darkMode,
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

export default useAppStore;
