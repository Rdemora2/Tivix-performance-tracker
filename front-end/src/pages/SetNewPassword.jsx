import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  Card, 
  Container, 
  Title, 
  Text, 
  PasswordInput, 
  Button, 
  Alert,
  Stack,
  Progress
} from '@mantine/core';
import { IconAlertCircle, IconLock, IconCheck } from '@tabler/icons-react';
import useAppStore from '../store/useAppStore';

const SetNewPassword = () => {
  const [passwords, setPasswords] = useState({
    newPassword: '',
    confirmPassword: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);
  
  const navigate = useNavigate();
  const { setNewPassword, user } = useAppStore();

  const checkPasswordStrength = (password) => {
    let strength = 0;
    if (password.length >= 8) strength += 25;
    if (/[A-Z]/.test(password)) strength += 25;
    if (/[a-z]/.test(password)) strength += 25;
    if (/[0-9]/.test(password)) strength += 12.5;
    if (/[^A-Za-z0-9]/.test(password)) strength += 12.5;
    return Math.min(strength, 100);
  };

  const handlePasswordChange = (value) => {
    setPasswords(prev => ({ ...prev, newPassword: value }));
    setPasswordStrength(checkPasswordStrength(value));
  };

  const getPasswordStrengthColor = (strength) => {
    if (strength < 50) return 'red';
    if (strength < 75) return 'yellow';
    return 'green';
  };

  const getPasswordStrengthLabel = (strength) => {
    if (strength < 25) return 'Muito fraca';
    if (strength < 50) return 'Fraca';
    if (strength < 75) return 'Média';
    if (strength < 90) return 'Forte';
    return 'Muito forte';
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (passwords.newPassword !== passwords.confirmPassword) {
      setError('As senhas não coincidem');
      setLoading(false);
      return;
    }

    if (passwordStrength < 50) {
      setError('A senha deve ser mais forte. Use pelo menos 8 caracteres, incluindo maiúsculas, minúsculas e números.');
      setLoading(false);
      return;
    }

    try {
      await setNewPassword({
        newPassword: passwords.newPassword
      });
      navigate('/', { replace: true });
    } catch (error) {
      setError(error.message || 'Erro ao definir nova senha');
    } finally {
      setLoading(false);
    }
  };

  // Se o usuário não precisa trocar a senha, redireciona
  if (!user?.needsPasswordChange) {
    navigate('/', { replace: true });
    return null;
  }

  return (
    <Container size="xs" pt={80}>
      <Card shadow="md" padding="xl" radius="md" withBorder>
        <Stack spacing="md">
          <div style={{ textAlign: 'center' }}>
            <IconLock size={60} color="var(--mantine-color-blue-6)" style={{ marginBottom: 16 }} />
            <Title order={2} mb="xs">
              Definir Nova Senha
            </Title>
            <Text c="dimmed" size="sm">
              Bem-vindo, {user?.name}! Por segurança, você deve definir uma nova senha antes de continuar.
            </Text>
          </div>

          <Alert
            icon={<IconAlertCircle size={16} />}
            color="blue"
            variant="light"
          >
            Esta é sua primeira vez no sistema. Por favor, defina uma senha segura para sua conta.
          </Alert>

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
              <div>
                <PasswordInput
                  label="Nova senha"
                  placeholder="Digite sua nova senha"
                  value={passwords.newPassword}
                  onChange={(e) => handlePasswordChange(e.target.value)}
                  required
                />
                {passwords.newPassword && (
                  <div style={{ marginTop: 8 }}>
                    <Progress 
                      value={passwordStrength} 
                      color={getPasswordStrengthColor(passwordStrength)}
                      size="sm"
                      mb={4}
                    />
                    <Text size="xs" c={getPasswordStrengthColor(passwordStrength)}>
                      Força da senha: {getPasswordStrengthLabel(passwordStrength)}
                    </Text>
                  </div>
                )}
              </div>

              <PasswordInput
                label="Confirmar nova senha"
                placeholder="Confirme sua nova senha"
                value={passwords.confirmPassword}
                onChange={(e) => setPasswords(prev => ({ ...prev, confirmPassword: e.target.value }))}
                required
              />

              <div>
                <Text size="sm" fw={500} mb="xs">Requisitos da senha:</Text>
                <Stack spacing={4}>
                  <Text size="xs" c={passwords.newPassword.length >= 8 ? 'green' : 'dimmed'}>
                    <IconCheck size={12} style={{ marginRight: 4 }} />
                    Pelo menos 8 caracteres
                  </Text>
                  <Text size="xs" c={/[A-Z]/.test(passwords.newPassword) ? 'green' : 'dimmed'}>
                    <IconCheck size={12} style={{ marginRight: 4 }} />
                    Uma letra maiúscula
                  </Text>
                  <Text size="xs" c={/[a-z]/.test(passwords.newPassword) ? 'green' : 'dimmed'}>
                    <IconCheck size={12} style={{ marginRight: 4 }} />
                    Uma letra minúscula
                  </Text>
                  <Text size="xs" c={/[0-9]/.test(passwords.newPassword) ? 'green' : 'dimmed'}>
                    <IconCheck size={12} style={{ marginRight: 4 }} />
                    Um número
                  </Text>
                  <Text size="xs" c={/[^A-Za-z0-9]/.test(passwords.newPassword) ? 'green' : 'dimmed'}>
                    <IconCheck size={12} style={{ marginRight: 4 }} />
                    Um caractere especial (recomendado)
                  </Text>
                </Stack>
              </div>

              <Button
                type="submit"
                fullWidth
                loading={loading}
                leftSection={<IconCheck size={16} />}
                disabled={passwordStrength < 50 || passwords.newPassword !== passwords.confirmPassword}
              >
                Definir Nova Senha
              </Button>
            </Stack>
          </form>
        </Stack>
      </Card>
    </Container>
  );
};

export default SetNewPassword;
