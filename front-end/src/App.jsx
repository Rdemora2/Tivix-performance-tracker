import { Routes, Route } from "react-router-dom";
import {
  MantineProvider,
  createTheme,
  LoadingOverlay,
  Container,
  Alert,
} from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { ModalsProvider } from "@mantine/modals";
import { Toaster } from "sonner";
import { useEffect } from "react";
import { IconAlertCircle } from "@tabler/icons-react";
import AppLayout from "./layouts/AppLayout";
import ProtectedRoute from "./components/ProtectedRoute";
import DashboardHome from "./pages/DashboardHome";
import Dashboard from "./pages/Dashboard";
import DeveloperProfile from "./pages/DeveloperProfile";
import CreateReport from "./pages/CreateReport";
import ConsolidatedReport from "./pages/ConsolidatedReport";
import Login from "./pages/Login";
import InitialSetup from "./pages/InitialSetup";
import Unauthorized from "./pages/Unauthorized";
import UserManagement from "./pages/UserManagement";
import CompanyManagement from "./pages/CompanyManagement";
import SetNewPassword from "./pages/SetNewPassword";
import useAppStore from "./store/useAppStore";
import useInitializeApp from "./hooks/useInitializeApp";
import "./App.css";

const theme = createTheme({
  primaryColor: "blue",
  fontFamily: "Inter, system-ui, Avenir, Helvetica, Arial, sans-serif",
  headings: {
    fontFamily: "Inter, system-ui, Avenir, Helvetica, Arial, sans-serif",
  },
});

function App() {
  const { darkMode } = useAppStore();
  const { loading, error } = useInitializeApp();

  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }, [darkMode]);

  if (loading) {
    return (
      <MantineProvider
        theme={theme}
        forceColorScheme={darkMode ? "dark" : "light"}
      >
        <Container size="lg" py="xl">
          <LoadingOverlay
            visible={true}
            zIndex={1000}
            overlayProps={{ radius: "sm", blur: 2 }}
          />
        </Container>
      </MantineProvider>
    );
  }

  if (error) {
    return (
      <MantineProvider
        theme={theme}
        forceColorScheme={darkMode ? "dark" : "light"}
      >
        <Container size="lg" py="xl">
          <Alert
            variant="light"
            color="red"
            title="Erro de Conexão"
            icon={<IconAlertCircle />}
          >
            Não foi possível conectar com o servidor. Verifique se o backend
            está disponível.
            <br />
            <strong>Erro:</strong> {error}
          </Alert>
        </Container>
      </MantineProvider>
    );
  }

  return (
    <MantineProvider
      theme={theme}
      forceColorScheme={darkMode ? "dark" : "light"}
    >
      <Notifications />
      <Toaster richColors />
      <ModalsProvider>
        <Routes>
          {/* Public routes */}
          <Route path="/login" element={<Login />} />
          <Route path="/initial-setup" element={<InitialSetup />} />
          <Route path="/set-new-password" element={<SetNewPassword />} />
          <Route path="/unauthorized" element={<Unauthorized />} />
          
          {/* Protected routes */}
          <Route path="/*" element={
            <ProtectedRoute>
              <AppLayout>
                <Routes>
                  <Route path="/" element={<DashboardHome />} />
                  <Route path="/team-dashboard" element={<Dashboard />} />
                  <Route path="/developer/:id" element={<DeveloperProfile />} />
                  <Route 
                    path="/developer/:id/create-report" 
                    element={
                      <ProtectedRoute requiredRole="manager">
                        <CreateReport />
                      </ProtectedRoute>
                    } 
                  />
                  <Route 
                    path="/consolidated-report" 
                    element={<ConsolidatedReport />} 
                  />
                  <Route 
                    path="/user-management" 
                    element={
                      <ProtectedRoute requiredRole={["admin", "manager"]}>
                        <UserManagement />
                      </ProtectedRoute>
                    } 
                  />
                  <Route 
                    path="/company-management" 
                    element={
                      <ProtectedRoute requiredRole="admin">
                        <CompanyManagement />
                      </ProtectedRoute>
                    } 
                  />
                </Routes>
              </AppLayout>
            </ProtectedRoute>
          } />
        </Routes>
      </ModalsProvider>
    </MantineProvider>
  );
}

export default App;
