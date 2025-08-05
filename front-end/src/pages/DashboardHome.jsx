import {
  Container,
  Title,
  Text,
  Group,
  Button,
  Card,
  SimpleGrid,
  Paper,
  Badge,
  Stack,
} from "@mantine/core";
import {
  IconArrowRight,
  IconUsers,
  IconChartLine,
  IconFileAnalytics,
  IconUserPlus,
} from "@tabler/icons-react";
import { useNavigate } from "react-router-dom";
import useAppStore from "../store/useAppStore";
import PermissionGuard from "../components/PermissionGuard";

const DashboardHome = () => {
  const navigate = useNavigate();
  const { developers, performanceReports } = useAppStore();

  const totalScores = performanceReports.reduce(
    (sum, report) => sum + report.weightedAverageScore,
    0
  );
  const overallAverageScore =
    performanceReports.length > 0 ? totalScores / performanceReports.length : 0;

  const sortedDevelopers = [...developers].sort(
    (a, b) => b.latestPerformanceScore - a.latestPerformanceScore
  );
  const topPerformers = sortedDevelopers.slice(0, 3);
  const bottomPerformers = sortedDevelopers.slice(-3);

  return (
    <Container size="xl">
      <Group justify="space-between" mb="xl">
        <div>
          <Title order={1} mb="xs">
            Visão Geral da Equipe
          </Title>
          <Text c="dimmed">
            Métricas consolidadas e insights rápidos sobre a performance do
            time.
          </Text>
        </div>
        <PermissionGuard action="create" resource="developers">
          <Button
            leftSection={<IconUserPlus size={16} />}
            onClick={() => navigate("/add-developer")}
          >
            Adicionar Novo Membro
          </Button>
        </PermissionGuard>
      </Group>

      <SimpleGrid cols={{ base: 1, sm: 2, lg: 4 }} spacing="lg" mb="xl">
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Group justify="space-between" mb="xs">
            <Text size="xs" c="dimmed" fw={700}>
              TOTAL DE MEMBROS
            </Text>
            <IconUsers size={20} color="var(--mantine-color-blue-6)" />
          </Group>
          <Text size="xl" fw={700}>
            {developers.length}
          </Text>
        </Card>

        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Group justify="space-between" mb="xs">
            <Text size="xs" c="dimmed" fw={700}>
              MÉDIA GERAL DA EQUIPE
            </Text>
            <IconChartLine size={20} color="var(--mantine-color-green-6)" />
          </Group>
          <Text size="xl" fw={700}>
            {overallAverageScore.toFixed(1)}/10
          </Text>
        </Card>

        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Group justify="space-between" mb="xs">
            <Text size="xs" c="dimmed" fw={700}>
              AVALIAÇÕES REGISTRADAS
            </Text>
            <IconFileAnalytics size={20} color="var(--mantine-color-grape-6)" />
          </Group>
          <Text size="xl" fw={700}>
            {performanceReports.length}
          </Text>
        </Card>

        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Group justify="space-between" mb="xs">
            <Text size="xs" c="dimmed" fw={700}>
              ÚLTIMA ATUALIZAÇÃO
            </Text>
            <IconFileAnalytics size={20} color="var(--mantine-color-gray-6)" />
          </Group>
          <Text size="xl" fw={700}>
            {performanceReports.length > 0
              ? new Date(performanceReports[0].createdAt).toLocaleDateString(
                  "pt-BR"
                )
              : "N/A"}
          </Text>
        </Card>
      </SimpleGrid>

      <SimpleGrid cols={{ base: 1, md: 2 }} spacing="lg" mb="xl">
        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Title order={3} mb="md">
            Melhores Performances
          </Title>
          <Stack>
            {topPerformers.length > 0 ? (
              topPerformers.map((dev) => (
                <Paper key={dev.id} p="xs" withBorder>
                  <Group justify="space-between">
                    <Text
                      style={{ cursor: "pointer" }}
                      c="blue"
                      onClick={() => navigate(`/developer/${dev.id}`)}
                    >
                      {dev.name}
                    </Text>
                    <Badge color="green">
                      {dev.latestPerformanceScore.toFixed(1)}/10
                    </Badge>
                  </Group>
                </Paper>
              ))
            ) : (
              <Text c="dimmed">Nenhum desenvolvedor para exibir.</Text>
            )}
          </Stack>
        </Card>

        <Card shadow="sm" padding="lg" radius="md" withBorder>
          <Title order={3} mb="md">
            A Desenvolver
          </Title>
          <Stack>
            {bottomPerformers.length > 0 ? (
              bottomPerformers.map((dev) => (
                <Paper key={dev.id} p="xs" withBorder>
                  <Group justify="space-between">
                    <Text
                      style={{ cursor: "pointer" }}
                      c="blue"
                      onClick={() => navigate(`/developer/${dev.id}`)}
                    >
                      {dev.name}
                    </Text>
                    <Badge color="orange">
                      {dev.latestPerformanceScore.toFixed(1)}/10
                    </Badge>
                  </Group>
                </Paper>
              ))
            ) : (
              <Text c="dimmed">Nenhum desenvolvedor para exibir.</Text>
            )}
          </Stack>
        </Card>
      </SimpleGrid>

      <Group justify="center" mt="xl">
        <Button
          variant="outline"
          rightSection={<IconArrowRight size={16} />}
          onClick={() => navigate("/team-dashboard")}
        >
          Ver Dashboard Detalhado da Equipe
        </Button>
      </Group>
    </Container>
  );
};

export default DashboardHome;
