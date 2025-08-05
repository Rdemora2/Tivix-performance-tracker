import { Container, Title, Text, Button, Stack } from '@mantine/core';
import { IconLock } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';

const Unauthorized = () => {
  const navigate = useNavigate();

  return (
    <Container size="sm" pt={100}>
      <Stack align="center" spacing="xl">
        <IconLock size={80} color="var(--mantine-color-red-6)" />
        
        <div style={{ textAlign: 'center' }}>
          <Title order={1} mb="md">
            Acesso Negado
          </Title>
          <Text size="lg" c="dimmed">
            Você não tem permissão para acessar esta página.
          </Text>
        </div>

        <Button 
          variant="light" 
          onClick={() => navigate('/')}
        >
          Voltar ao Início
        </Button>
      </Stack>
    </Container>
  );
};

export default Unauthorized;
