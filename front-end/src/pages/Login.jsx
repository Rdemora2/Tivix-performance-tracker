import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import {
  Card,
  Container,
  Title,
  TextInput,
  PasswordInput,
  Button,
  Text,
  Alert,
  Stack,
  Anchor,
} from "@mantine/core";
import { IconAlertCircle, IconLogin } from "@tabler/icons-react";
import useAppStore from "../store/useAppStore";

const Login = () => {
  const [credentials, setCredentials] = useState({
    email: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();
  const login = useAppStore((state) => state.login);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await login(credentials);
      navigate("/", { replace: true });
    } catch (error) {
      setError(error.message || "Erro ao fazer login");
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field) => (e) => {
    setCredentials((prev) => ({
      ...prev,
      [field]: e.target.value,
    }));
  };

  return (
    <Container size="xs" pt={100}>
      <Card shadow="md" padding="xl" radius="md" withBorder>
        <Stack spacing="md">
          <div style={{ textAlign: "center" }}>
            <Title order={2} mb="xs">
              Performance Tracker
            </Title>
            <Text c="dimmed" size="sm">
              Faça login para acessar o sistema
            </Text>
          </div>

          {error && (
            <Alert
              icon={<IconAlertCircle size={16} />}
              color="red"
              variant="light"
            >
              {error}
            </Alert>
          )}

          <form onSubmit={handleSubmit}>
            <Stack spacing="md">
              <TextInput
                label="Email"
                placeholder="seu.email@tivix.com"
                value={credentials.email}
                onChange={handleChange("email")}
                required
                type="email"
              />

              <PasswordInput
                label="Senha"
                placeholder="Sua senha"
                value={credentials.password}
                onChange={handleChange("password")}
                required
              />

              <Button
                type="submit"
                fullWidth
                loading={loading}
                leftSection={<IconLogin size={16} />}
              >
                Entrar
              </Button>
            </Stack>
          </form>

          <Text ta="center" size="sm" c="dimmed" mt="md">
            Esqueceu sua senha? Entre em contato com o administrador do sistema.
          </Text>
        </Stack>
      </Card>
    </Container>
  );
};

export default Login;
