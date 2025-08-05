import { useState, useEffect } from 'react';
import { 
  Container,
  Title,
  Card,
  Button,
  Table,
  Badge,
  Group,
  Modal,
  TextInput,
  Select,
  Stack,
  Text,
  Alert,
  ActionIcon,
  Tooltip
} from '@mantine/core';
import { 
  IconUserPlus, 
  IconCopy, 
  IconEye, 
  IconEyeOff,
  IconCheck,
  IconAlertCircle,
  IconRefresh
} from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import useAppStore from '../store/useAppStore';

const UserManagement = () => {
  const [opened, setOpened] = useState(false);
  const [loading, setLoading] = useState(false);
  const [loadingUsers, setLoadingUsers] = useState(true);
  const [newUserData, setNewUserData] = useState({
    name: '',
    email: '',
    role: 'user'
  });
  const [showPassword, setShowPassword] = useState(false);
  const [generatedPassword, setGeneratedPassword] = useState('');

  const { createUser, fetchUsers, user, users } = useAppStore();

  // Carregar usuários quando o componente montar
  useEffect(() => {
    if (user?.role === 'admin') {
      setLoadingUsers(true);
      fetchUsers()
        .catch(error => {
          notifications.show({
            title: 'Erro ao carregar usuários',
            message: error.message,
            color: 'red'
          });
        })
        .finally(() => {
          setLoadingUsers(false);
        });
    } else {
      setLoadingUsers(false);
    }
  }, [fetchUsers, user]);

  const generateRandomPassword = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%&*';
    let password = '';
    for (let i = 0; i < 12; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return password;
  };

  const handleCreateUser = async () => {
    if (!newUserData.name || !newUserData.email) {
      notifications.show({
        title: 'Erro',
        message: 'Nome e email são obrigatórios',
        color: 'red'
      });
      return;
    }

    const tempPassword = generateRandomPassword();
    setGeneratedPassword(tempPassword);

    try {
      setLoading(true);
      await createUser({
        ...newUserData,
        temporaryPassword: tempPassword
      });

      notifications.show({
        title: 'Usuário criado com sucesso!',
        message: 'Senha temporária gerada. Copie e envie para o usuário.',
        color: 'green'
      });

      // Reset form but keep modal open to show password
      setNewUserData({ name: '', email: '', role: 'user' });
      
      // Recarregar lista de usuários
      await fetchUsers();
    } catch (error) {
      notifications.show({
        title: 'Erro ao criar usuário',
        message: error.message,
        color: 'red'
      });
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    notifications.show({
      title: 'Copiado!',
      message: 'Senha copiada para a área de transferência',
      color: 'blue'
    });
  };

  const refreshUsers = async () => {
    setLoadingUsers(true);
    try {
      await fetchUsers();
      notifications.show({
        title: 'Lista atualizada!',
        message: 'Usuários recarregados com sucesso',
        color: 'green'
      });
    } catch (error) {
      notifications.show({
        title: 'Erro ao recarregar',
        message: error.message,
        color: 'red'
      });
    } finally {
      setLoadingUsers(false);
    }
  };

  const getRoleColor = (role) => {
    switch (role) {
      case 'admin': return 'red';
      case 'manager': return 'blue';
      case 'user': return 'green';
      default: return 'gray';
    }
  };

  const getRoleLabel = (role) => {
    switch (role) {
      case 'admin': return 'Administrador';
      case 'manager': return 'Gerente';
      case 'user': return 'Usuário';
      default: return role;
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    return new Intl.DateTimeFormat('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(new Date(dateString));
  };

  if (user?.role !== 'admin') {
    return (
      <Container size="sm" pt={100}>
        <Alert
          icon={<IconAlertCircle size={16} />}
          color="red"
          title="Acesso Negado"
        >
          Apenas administradores podem gerenciar usuários.
        </Alert>
      </Container>
    );
  }

  return (
    <Container size="xl">
      <Group justify="space-between" mb="xl">
        <div>
          <Title order={2}>Gerenciamento de Usuários</Title>
          <Text c="dimmed">
            Gerencie usuários do sistema
          </Text>
        </div>
        <Group>
          <Button 
            variant="light"
            leftSection={<IconRefresh size={16} />}
            onClick={refreshUsers}
            loading={loadingUsers}
            disabled={loadingUsers}
          >
            Atualizar
          </Button>
          <Button 
            leftSection={<IconUserPlus size={16} />}
            onClick={() => setOpened(true)}
          >
            Novo Usuário
          </Button>
        </Group>
      </Group>

      <Card shadow="sm" padding="lg" radius="md" withBorder>
        {loadingUsers ? (
          <div style={{ textAlign: 'center', padding: '2rem' }}>
            <Text c="dimmed">Carregando usuários...</Text>
          </div>
        ) : (
          <Table>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Nome</Table.Th>
                <Table.Th>Email</Table.Th>
                <Table.Th>Função</Table.Th>
                <Table.Th>Status</Table.Th>
                <Table.Th>Criado em</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {users.map((userItem) => (
                <Table.Tr key={userItem.id}>
                  <Table.Td>{userItem.name}</Table.Td>
                  <Table.Td>{userItem.email}</Table.Td>
                  <Table.Td>
                    <Badge color={getRoleColor(userItem.role)}>
                      {getRoleLabel(userItem.role)}
                    </Badge>
                  </Table.Td>
                  <Table.Td>
                    {userItem.needsPasswordChange ? (
                      <Badge color="yellow">Aguardando nova senha</Badge>
                    ) : userItem.isActive !== false ? (
                      <Badge color="green">Ativo</Badge>
                    ) : (
                      <Badge color="gray">Inativo</Badge>
                    )}
                  </Table.Td>
                  <Table.Td>
                    <Text size="sm" c="dimmed">
                      {formatDate(userItem.createdAt)}
                    </Text>
                  </Table.Td>
                </Table.Tr>
              ))}
              {users.length === 0 && (
                <Table.Tr>
                  <Table.Td colSpan={5} style={{ textAlign: 'center', padding: '2rem' }}>
                    <Text c="dimmed">Nenhum usuário encontrado</Text>
                  </Table.Td>
                </Table.Tr>
              )}
            </Table.Tbody>
          </Table>
        )}
      </Card>

      <Modal
        opened={opened}
        onClose={() => {
          setOpened(false);
          setGeneratedPassword('');
          setNewUserData({ name: '', email: '', role: 'user' });
        }}
        title="Criar Novo Usuário"
        size="md"
      >
        <Stack spacing="md">
          {generatedPassword ? (
            <>
              <Alert
                icon={<IconCheck size={16} />}
                color="green"
                title="Usuário criado com sucesso!"
              >
                Uma senha temporária foi gerada. O usuário deverá definir uma nova senha no primeiro login.
              </Alert>

              <Card padding="md" withBorder>
                <Group justify="space-between" mb="xs">
                  <Text fw={500}>Senha Temporária:</Text>
                  <Group gap="xs">
                    <Tooltip label={showPassword ? "Ocultar senha" : "Mostrar senha"}>
                      <ActionIcon 
                        variant="subtle" 
                        onClick={() => setShowPassword(!showPassword)}
                      >
                        {showPassword ? <IconEyeOff size={16} /> : <IconEye size={16} />}
                      </ActionIcon>
                    </Tooltip>
                    <Tooltip label="Copiar senha">
                      <ActionIcon 
                        variant="subtle" 
                        onClick={() => copyToClipboard(generatedPassword)}
                      >
                        <IconCopy size={16} />
                      </ActionIcon>
                    </Tooltip>
                  </Group>
                </Group>
                <Text 
                  ff="monospace" 
                  size="sm" 
                  p="sm" 
                  bg="gray.1" 
                  style={{ borderRadius: 4 }}
                >
                  {showPassword ? generatedPassword : '*'.repeat(generatedPassword.length)}
                </Text>
                <Text size="xs" c="dimmed" mt="xs">
                  Copie esta senha e envie para o usuário de forma segura.
                </Text>
              </Card>

              <Button 
                onClick={() => {
                  setOpened(false);
                  setGeneratedPassword('');
                  setNewUserData({ name: '', email: '', role: 'user' });
                }}
                fullWidth
              >
                Fechar
              </Button>
            </>
          ) : (
            <>
              <TextInput
                label="Nome completo"
                placeholder="Nome do usuário"
                value={newUserData.name}
                onChange={(e) => setNewUserData(prev => ({ ...prev, name: e.target.value }))}
                required
              />

              <TextInput
                label="Email"
                placeholder="email@tivix.com"
                value={newUserData.email}
                onChange={(e) => setNewUserData(prev => ({ ...prev, email: e.target.value }))}
                required
                type="email"
              />

              <Select
                label="Função"
                value={newUserData.role}
                onChange={(value) => setNewUserData(prev => ({ ...prev, role: value }))}
                data={[
                  { value: 'user', label: 'Usuário' },
                  { value: 'manager', label: 'Gerente' },
                  { value: 'admin', label: 'Administrador' }
                ]}
                required
              />

              <Group justify="space-between" mt="md">
                <Button 
                  variant="subtle" 
                  onClick={() => setOpened(false)}
                >
                  Cancelar
                </Button>
                <Button 
                  onClick={handleCreateUser}
                  loading={loading}
                  leftSection={<IconUserPlus size={16} />}
                >
                  Criar Usuário
                </Button>
              </Group>
            </>
          )}
        </Stack>
      </Modal>
    </Container>
  );
};

export default UserManagement;
