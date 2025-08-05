import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
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
  Select
} from '@mantine/core';
import { IconAlertCircle, IconUserPlus } from '@tabler/icons-react';
import useAppStore from '../store/useAppStore';

const Register = () => {
  const [userData, setUserData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
    role: 'user'
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  
  const navigate = useNavigate();
  const register = useAppStore(state => state.register);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (userData.password !== userData.confirmPassword) {
      setError('As senhas não coincidem');
      setLoading(false);
      return;
    }

    try {
      const { confirmPassword: _, ...registrationData } = userData;
      await register(registrationData);
      navigate('/', { replace: true });
    } catch (error) {
      setError(error.message || 'Erro ao criar conta');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field) => (value) => {
    const actualValue = typeof value === 'string' ? value : value.target.value;
    setUserData(prev => ({
      ...prev,
      [field]: actualValue
    }));
  };

  return (
    <Container size="xs" pt={80}>
      <Card shadow="md" padding="xl" radius="md" withBorder>
        <Stack spacing="md">
          <div style={{ textAlign: 'center' }}>
            <Title order={2} mb="xs">
              Criar Conta
            </Title>
            <Text c="dimmed" size="sm">
              Cadastre-se no sistema
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
                label="Nome completo"
                placeholder="Seu nome"
                value={userData.name}
                onChange={handleChange('name')}
                required
              />

              <TextInput
                label="Email"
                placeholder="seu.email@tivix.com"
                value={userData.email}
                onChange={handleChange('email')}
                required
                type="email"
              />

              <Select
                label="Função"
                placeholder="Selecione sua função"
                value={userData.role}
                onChange={handleChange('role')}
                data={[
                  { value: 'user', label: 'Usuário' },
                  { value: 'manager', label: 'Gerente' }
                ]}
                required
              />

              <PasswordInput
                label="Senha"
                placeholder="Sua senha"
                value={userData.password}
                onChange={handleChange('password')}
                required
              />

              <PasswordInput
                label="Confirmar senha"
                placeholder="Confirme sua senha"
                value={userData.confirmPassword}
                onChange={handleChange('confirmPassword')}
                required
              />

              <Button
                type="submit"
                fullWidth
                loading={loading}
                leftSection={<IconUserPlus size={16} />}
              >
                Criar Conta
              </Button>
            </Stack>
          </form>

          <Text ta="center" size="sm" c="dimmed">
            Já tem uma conta?{' '}
            <Anchor component={Link} to="/login">
              Faça login aqui
            </Anchor>
          </Text>
        </Stack>
      </Card>
    </Container>
  );
};

export default Register;
