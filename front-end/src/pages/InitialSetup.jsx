import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
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
  LoadingOverlay
} from '@mantine/core';
import { IconAlertCircle, IconUserPlus, IconShield } from '@tabler/icons-react';
import api from '../services/api';

const InitialSetup = () => {
  const [adminData, setAdminData] = useState({
    installKey: 'TIVIX_INSTALL_2024',
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [checking, setChecking] = useState(true);
  
  const navigate = useNavigate();

  useEffect(() => {
    const checkInitStatus = async () => {
      try {
        const response = await api.init.check();
        if (response.initialized || response.userCount > 0) {
          navigate('/login', { replace: true });
        }
      } catch (error) {
        console.error('Error checking init status:', error);
      } finally {
        setChecking(false);
      }
    };

    checkInitStatus();
  }, [navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (adminData.password !== adminData.confirmPassword) {
      setError('As senhas não coincidem');
      setLoading(false);
      return;
    }

    try {
      const { confirmPassword: _, ...createData } = adminData;
      await api.init.createAdmin(createData);
      navigate('/login', { replace: true });
    } catch (error) {
      setError(error.message || 'Erro ao criar administrador');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field) => (e) => {
    setAdminData(prev => ({
      ...prev,
      [field]: e.target.value
    }));
  };

  if (checking) {
    return (
      <Container size="xs" pt={100}>
        <LoadingOverlay visible />
      </Container>
    );
  }

  return (
    <Container size="xs" pt={80}>
      <Card shadow="md" padding="xl" radius="md" withBorder>
        <Stack spacing="md">
          <div style={{ textAlign: 'center' }}>
            <IconShield size={60} color="var(--mantine-color-blue-6)" style={{ marginBottom: 16 }} />
            <Title order={2} mb="xs">
              Configuração Inicial
            </Title>
            <Text c="dimmed" size="sm">
              Crie o primeiro administrador do sistema
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
                label="Chave de instalação"
                value={adminData.installKey}
                onChange={handleChange('installKey')}
                required
                readOnly
                variant="filled"
              />

              <TextInput
                label="Nome do administrador"
                placeholder="Nome completo"
                value={adminData.name}
                onChange={handleChange('name')}
                required
              />

              <TextInput
                label="Email"
                placeholder="admin@tivix.com"
                value={adminData.email}
                onChange={handleChange('email')}
                required
                type="email"
              />

              <PasswordInput
                label="Senha"
                placeholder="Senha do administrador"
                value={adminData.password}
                onChange={handleChange('password')}
                required
              />

              <PasswordInput
                label="Confirmar senha"
                placeholder="Confirme a senha"
                value={adminData.confirmPassword}
                onChange={handleChange('confirmPassword')}
                required
              />

              <Button
                type="submit"
                fullWidth
                loading={loading}
                leftSection={<IconUserPlus size={16} />}
              >
                Criar Administrador
              </Button>
            </Stack>
          </form>
        </Stack>
      </Card>
    </Container>
  );
};

export default InitialSetup;
