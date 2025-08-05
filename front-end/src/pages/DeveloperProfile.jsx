import {
  Container,
  Title,
  Card,
  Text,
  Group,
  Button,
  Badge,
  Grid,
  Paper,
  Stack,
  Modal,
  Table,
  Progress,
} from "@mantine/core";
import {
  IconArrowLeft,
  IconPlus,
  IconTrendingUp,
  IconCalendar,
  IconEye,
} from "@tabler/icons-react";
import { useParams, useNavigate } from "react-router-dom";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
} from "recharts";
import { useState } from "react";
import useAppStore from "../store/useAppStore";
import { EVALUATION_CATEGORIES } from "../types";

const DeveloperProfile = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { developers, getReportsByDeveloper } = useAppStore();
  const [selectedReportId, setSelectedReportId] = useState(null);
  const [detailsModalOpened, setDetailsModalOpened] = useState(false);
  const [reportForDetails, setReportForDetails] = useState(null);

  const developer = developers.find((dev) => dev.id === id);
  const reports = getReportsByDeveloper(id);
  const selectedReport =
    reports.find((report) => report.id === selectedReportId) || reports[0];

  if (!developer) {
    return (
      <Container size="xl">
        <Text>Desenvolvedor não encontrado</Text>
      </Container>
    );
  }

  const performanceData = reports
    .map((report) => ({
      month: new Date(report.month + "-01").toLocaleDateString("pt-BR", {
        month: "short",
        year: "numeric",
      }),
      score: report.weightedAverageScore,
    }))
    .reverse();

  const radarData = selectedReport
    ? Object.entries(
        selectedReport.categoryScores || selectedReport.scores || {}
      ).map(([key, value]) => ({
        category: EVALUATION_CATEGORIES[key]?.label || key,
        score: value,
        fullMark: 10,
      }))
    : [];

  const getPerformanceColor = (score) => {
    if (score >= 8) return "green";
    if (score >= 6) return "yellow";
    if (score >= 4) return "orange";
    return "red";
  };

  const getPerformanceLabel = (score) => {
    if (score >= 8) return "Excelente";
    if (score >= 6) return "Bom";
    if (score >= 4) return "Regular";
    return "Precisa Melhorar";
  };

  const handleViewDetails = (report) => {
    setReportForDetails(report);
    setDetailsModalOpened(true);
  };

  const getScoreColor = (score) => {
    if (score >= 8) return "green";
    if (score >= 6) return "blue";
    if (score >= 4) return "yellow";
    return "red";
  };

  return (
    <Container size="xl">
      <Group mb="xl">
        <Button
          variant="subtle"
          leftSection={<IconArrowLeft size={16} />}
          onClick={() => navigate("/")}
        >
          Voltar ao Dashboard
        </Button>
      </Group>

      <Grid>
        <Grid.Col span={12}>
          <Card shadow="sm" padding="lg" radius="md" withBorder>
            <Group justify="space-between" mb="md">
              <div>
                <Title order={2} mb="xs">
                  {developer.name}
                </Title>
                <Text c="dimmed" size="lg">
                  {developer.role}
                </Text>
              </div>
              <div style={{ textAlign: "right" }}>
                <Badge
                  color={getPerformanceColor(developer.latestPerformanceScore)}
                  variant="light"
                  size="lg"
                >
                  {getPerformanceLabel(developer.latestPerformanceScore)}
                </Badge>
                <Group gap="xs" mt="xs" justify="flex-end">
                  <IconTrendingUp
                    size={16}
                    color="var(--mantine-color-blue-6)"
                  />
                  <Text size="lg" fw={700} c="blue">
                    {developer.latestPerformanceScore.toFixed(1)}/10
                  </Text>
                </Group>
              </div>
            </Group>

            <Button
              leftSection={<IconPlus size={16} />}
              onClick={() => navigate(`/developer/${id}/create-report`)}
            >
              Nova Avaliação
            </Button>
          </Card>
        </Grid.Col>

        {performanceData.length > 0 && (
          <Grid.Col span={12}>
            <Card shadow="sm" padding="lg" radius="md" withBorder>
              <Title order={3} mb="md">
                Evolução da Performance
              </Title>
              <div style={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={performanceData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="month" />
                    <YAxis domain={[0, 10]} />
                    <Tooltip
                      formatter={(value) => [
                        `${value.toFixed(1)}/10`,
                        "Performance",
                      ]}
                      labelFormatter={(label) => `Mês: ${label}`}
                    />
                    <Line
                      type="monotone"
                      dataKey="score"
                      stroke="var(--mantine-color-blue-6)"
                      strokeWidth={3}
                      dot={{
                        fill: "var(--mantine-color-blue-6)",
                        strokeWidth: 2,
                        r: 6,
                      }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </Card>
          </Grid.Col>
        )}

        <Grid.Col span={{ base: 12, md: 6 }}>
          <Card shadow="sm" padding="lg" radius="md" withBorder h="100%">
            <Title order={3} mb="md">
              Histórico de Relatórios
            </Title>
            {reports.length > 0 ? (
              <Stack gap="sm">
                {reports.map((report) => (
                  <Paper
                    key={report.id}
                    p="md"
                    withBorder
                    style={{
                      backgroundColor:
                        selectedReport?.id === report.id
                          ? "var(--mantine-color-blue-light)"
                          : undefined,
                    }}
                  >
                    <Group justify="space-between">
                      <div
                        style={{ flex: 1, cursor: "pointer" }}
                        onClick={() => setSelectedReportId(report.id)}
                      >
                        <Group gap="xs" mb="xs">
                          <IconCalendar size={16} />
                          <Text fw={500}>
                            {new Date(report.month + "-01").toLocaleDateString(
                              "pt-BR",
                              { month: "long", year: "numeric" }
                            )}
                          </Text>
                        </Group>
                        <Text size="sm" c="dimmed">
                          Performance: {report.weightedAverageScore.toFixed(1)}
                          /10
                        </Text>
                      </div>
                      <Stack gap="xs" align="flex-end">
                        <Badge
                          color={getPerformanceColor(report.weightedAverageScore)}
                          variant="light"
                        >
                          {getPerformanceLabel(report.weightedAverageScore)}
                        </Badge>
                        <Button
                          size="xs"
                          variant="subtle"
                          leftSection={<IconEye size={14} />}
                          onClick={(e) => {
                            e.stopPropagation();
                            handleViewDetails(report);
                          }}
                        >
                          Ver Notas
                        </Button>
                      </Stack>
                    </Group>
                  </Paper>
                ))}
              </Stack>
            ) : (
              <Text c="dimmed" ta="center" py="xl">
                Nenhum relatório encontrado
              </Text>
            )}
          </Card>
        </Grid.Col>

        {selectedReport && (
          <Grid.Col span={{ base: 12, md: 6 }}>
            <Card shadow="sm" padding="lg" radius="md" withBorder h="100%">
              <Title order={3} mb="md">
                Detalhes -{" "}
                {new Date(selectedReport.month + "-01").toLocaleDateString(
                  "pt-BR",
                  { month: "long", year: "numeric" }
                )}
              </Title>

              <div style={{ height: 250, marginBottom: "1rem" }}>
                <ResponsiveContainer width="100%" height="100%">
                  <RadarChart data={radarData}>
                    <PolarGrid />
                    <PolarAngleAxis dataKey="category" />
                    <PolarRadiusAxis angle={90} domain={[0, 10]} />
                    <Radar
                      name="Performance"
                      dataKey="score"
                      stroke="var(--mantine-color-blue-6)"
                      fill="var(--mantine-color-blue-6)"
                      fillOpacity={0.3}
                    />
                  </RadarChart>
                </ResponsiveContainer>
              </div>

              <Stack gap="sm">
                <div>
                  <Text fw={500} size="sm" mb="xs">
                    Destaques do Mês:
                  </Text>
                  <Text size="sm" c="dimmed">
                    {selectedReport.highlights || "Nenhum destaque registrado"}
                  </Text>
                </div>
                <div>
                  <Text fw={500} size="sm" mb="xs">
                    Pontos a Desenvolver:
                  </Text>
                  <Text size="sm" c="dimmed">
                    {selectedReport.pointsToDevelop ||
                      "Nenhum ponto registrado"}
                  </Text>
                </div>
              </Stack>
            </Card>
          </Grid.Col>
        )}
      </Grid>

      {/* Modal de Detalhes das Notas */}
      <Modal
        opened={detailsModalOpened}
        onClose={() => setDetailsModalOpened(false)}
        title={
          reportForDetails && (
            <Group gap="xs">
              <Text fw={600}>Notas Detalhadas -</Text>
              <Text>
                {new Date(reportForDetails.month + "-01").toLocaleDateString(
                  "pt-BR",
                  { month: "long", year: "numeric" }
                )}
              </Text>
            </Group>
          )
        }
        size="lg"
        centered
      >
        {reportForDetails && (
          <Stack gap="lg">
            {/* Resumo Geral */}
            <Card withBorder p="md">
              <Group justify="space-between" mb="md">
                <Text fw={600} size="lg">
                  Performance Geral
                </Text>
                <Badge
                  size="lg"
                  color={getPerformanceColor(reportForDetails.weightedAverageScore)}
                >
                  {reportForDetails.weightedAverageScore.toFixed(1)}/10
                </Badge>
              </Group>
              <Progress
                value={(reportForDetails.weightedAverageScore / 10) * 100}
                color={getScoreColor(reportForDetails.weightedAverageScore)}
                size="lg"
                radius="md"
              />
            </Card>

            {/* Notas por Categoria */}
            <Stack gap="md">
              <Text fw={600} size="md">
                Notas por Categoria
              </Text>
              
              {Object.entries(EVALUATION_CATEGORIES).map(([categoryKey, category]) => {
                const categoryScore = reportForDetails.categoryScores?.[categoryKey] || 0;
                
                return (
                  <Card key={categoryKey} withBorder p="md">
                    <Group justify="space-between" mb="sm">
                      <Text fw={500}>{category.label}</Text>
                      <Badge color={getScoreColor(categoryScore)}>
                        {categoryScore.toFixed(1)}/10
                      </Badge>
                    </Group>
                    
                    <Progress
                      value={(categoryScore / 10) * 100}
                      color={getScoreColor(categoryScore)}
                      mb="md"
                    />

                    {/* Perguntas da Categoria */}
                    <Stack gap="xs">
                      {category.questions.map((question) => {
                        const questionScore = reportForDetails.questionScores?.[question.key] || 0;
                        
                        return (
                          <Group key={question.key} justify="space-between">
                            <Text size="sm" c="dimmed" style={{ flex: 1 }}>
                              {question.label}
                            </Text>
                            <Group gap="xs">
                              <Text size="sm" fw={500}>
                                {questionScore.toFixed(1)}/10
                              </Text>
                              <Text size="xs" c="dimmed">
                                (Peso: {question.weight})
                              </Text>
                            </Group>
                          </Group>
                        );
                      })}
                    </Stack>
                  </Card>
                );
              })}
            </Stack>

            {/* Observações */}
            {(reportForDetails.highlights || reportForDetails.pointsToDevelop) && (
              <Stack gap="md">
                <Text fw={600} size="md">
                  Observações
                </Text>
                
                {reportForDetails.highlights && (
                  <Card withBorder p="md">
                    <Text fw={500} size="sm" mb="xs" c="green">
                      🌟 Destaques do Mês
                    </Text>
                    <Text size="sm">{reportForDetails.highlights}</Text>
                  </Card>
                )}
                
                {reportForDetails.pointsToDevelop && (
                  <Card withBorder p="md">
                    <Text fw={500} size="sm" mb="xs" c="orange">
                      📈 Pontos a Desenvolver
                    </Text>
                    <Text size="sm">{reportForDetails.pointsToDevelop}</Text>
                  </Card>
                )}
              </Stack>
            )}
          </Stack>
        )}
      </Modal>
    </Container>
  );
};

export default DeveloperProfile;
