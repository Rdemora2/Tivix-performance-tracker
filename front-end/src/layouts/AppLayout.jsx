import {
  AppShell,
  Title,
  Group,
  ActionIcon,
  Text,
  Menu,
  Avatar,
  Badge,
} from "@mantine/core";
import {
  IconSun,
  IconMoon,
  IconLogout,
  IconUser,
  IconChevronDown,
  IconUsers,
} from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import useAppStore from "../store/useAppStore";

const AppLayout = ({ children }) => {
  const { darkMode, toggleDarkMode, user, logout } = useAppStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    window.location.href = "/login";
  };

  const getRoleColor = (role) => {
    switch (role) {
      case "admin":
        return "red";
      case "manager":
        return "blue";
      case "user":
        return "green";
      default:
        return "gray";
    }
  };

  const getRoleLabel = (role) => {
    switch (role) {
      case "admin":
        return "Administrador";
      case "manager":
        return "Gerente";
      case "user":
        return "Usuário";
      default:
        return role;
    }
  };

  return (
    <AppShell header={{ height: 70 }} padding="md">
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <Title
            order={2}
            c="blue"
            style={{ cursor: "pointer" }}
            onClick={() => navigate("/")}
          >
            Performance Tracker
          </Title>

          <Group gap="md">
            <ActionIcon
              variant="outline"
              size="lg"
              onClick={toggleDarkMode}
              title={darkMode ? "Modo Claro" : "Modo Escuro"}
            >
              {darkMode ? <IconSun size={18} /> : <IconMoon size={18} />}
            </ActionIcon>

            {user && (
              <Menu shadow="md" width={200}>
                <Menu.Target>
                  <Group gap="xs" style={{ cursor: "pointer" }}>
                    <Avatar size="sm" color="blue">
                      <IconUser size={16} />
                    </Avatar>
                    <div>
                      <Text size="sm" fw={500}>
                        {user.name}
                      </Text>
                      <Badge size="xs" color={getRoleColor(user.role)}>
                        {getRoleLabel(user.role)}
                      </Badge>
                    </div>
                    <IconChevronDown size={14} />
                  </Group>
                </Menu.Target>

                <Menu.Dropdown>
                  {user.role === "admin" && (
                    <Menu.Item
                      leftSection={<IconUsers size={14} />}
                      onClick={() => navigate("/user-management")}
                    >
                      Gerenciar Usuários
                    </Menu.Item>
                  )}
                  <Menu.Item
                    leftSection={<IconLogout size={14} />}
                    onClick={handleLogout}
                    color="red"
                  >
                    Sair
                  </Menu.Item>
                </Menu.Dropdown>
              </Menu>
            )}
          </Group>
        </Group>
      </AppShell.Header>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
};

export default AppLayout;
