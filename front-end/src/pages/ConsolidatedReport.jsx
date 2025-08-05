import {
  Container,
  Title,
  Card,
  Text,
  Group,
  Button,
  Select,
  Table,
  Badge,
  Paper,
  Stack,
} from "@mantine/core";
import {
  IconArrowLeft,
  IconDownload,
  IconFileTypePdf,
  IconCheck,
  IconX,
} from "@tabler/icons-react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { notifications } from "@mantine/notifications";
import useAppStore from "../store/useAppStore";
import { EVALUATION_CATEGORIES } from "../types";
import jsPDF from "jspdf";
import html2canvas from "html2canvas";

const ConsolidatedReport = () => {
  const navigate = useNavigate();
  const { developers, performanceReports } = useAppStore();
  const [selectedMonth, setSelectedMonth] = useState("");

  const availableMonths = [
    ...new Set(performanceReports.map((report) => report.month)),
  ]
    .sort((a, b) => new Date(b) - new Date(a))
    .map((month) => ({
      value: month,
      label: new Date(month + "-01").toLocaleDateString("pt-BR", {
        month: "long",
        year: "numeric",
      }),
    }));

  const monthlyReports = selectedMonth
    ? performanceReports.filter((report) => report.month === selectedMonth)
    : [];

  const tableData = monthlyReports
    .map((report) => {
      const developer = developers.find((dev) => dev.id === report.developerId);
      return {
        ...report,
        developerName: developer?.name || "N/A",
        developerRole: developer?.role || "N/A",
      };
    })
    .sort((a, b) => b.weightedAverageScore - a.weightedAverageScore);

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

  const handleExportPDF = () => {
    const input = document.getElementById("consolidated-report-content");
    if (!input) {
      notifications.show({
        title: "Erro",
        message:
          "Não foi possível encontrar o conteúdo do relatório para exportar.",
        color: "red",
      });
      return;
    }

    notifications.show({
      title: "Gerando PDF",
      message: "Por favor, aguarde enquanto o relatório é gerado...",
      color: "blue",
      loading: true,
      autoClose: false,
      id: "pdf-generation",
    });

    const pdfStyle = document.createElement("style");
    pdfStyle.id = "pdf-export-style";
    pdfStyle.textContent = `
      #consolidated-report-content * {
        color: #000000 !important;
        background-color: #ffffff !important;
        border-color: #cccccc !important;
      }
      #consolidated-report-content .mantine-Badge-root {
        background-color: #e3f2fd !important;
        color: #1976d2 !important;
      }
      #consolidated-report-content .mantine-Text-root[data-c="blue"] {
        color: #1976d2 !important;
      }
      #consolidated-report-content .mantine-Text-root[data-c="green"] {
        color: #388e3c !important;
      }
      #consolidated-report-content .mantine-Text-root[data-c="orange"] {
        color: #f57c00 !important;
      }
      #consolidated-report-content .mantine-Text-root[data-c="red"] {
        color: #d32f2f !important;
      }
    `;
    document.head.appendChild(pdfStyle);

    const options = {
      scale: 2,
      useCORS: true,
      allowTaint: true,
      backgroundColor: "#ffffff",
      onclone: (clonedDoc) => {
        const styles = clonedDoc.querySelectorAll(
          'style, link[rel="stylesheet"]'
        );
        styles.forEach((style) => {
          if (style.textContent && style.textContent.includes("oklch")) {
            style.remove();
          }
        });
      },
    };

    html2canvas(input, options)
      .then((canvas) => {
        const tempStyle = document.getElementById("pdf-export-style");
        if (tempStyle) {
          tempStyle.remove();
        }

        const imgData = canvas.toDataURL("image/png");
        const pdf = new jsPDF("p", "mm", "a4");
        const imgWidth = 210;
        const pageHeight = 297;
        const imgHeight = (canvas.height * imgWidth) / canvas.width;
        let heightLeft = imgHeight;
        let position = 0;

        pdf.addImage(imgData, "PNG", 0, position, imgWidth, imgHeight);
        heightLeft -= pageHeight;

        while (heightLeft >= 0) {
          position = heightLeft - imgHeight;
          pdf.addPage();
          pdf.addImage(imgData, "PNG", 0, position, imgWidth, imgHeight);
          heightLeft -= pageHeight;
        }

        pdf.save(`relatorio-consolidado-${selectedMonth}.pdf`);
        notifications.update({
          id: "pdf-generation",
          title: "Sucesso",
          message: "Relatório PDF gerado com sucesso!",
          color: "green",
          icon: <IconCheck size={16} />,
          autoClose: 5000,
        });
      })
      .catch((error) => {
        const tempStyle = document.getElementById("pdf-export-style");
        if (tempStyle) {
          tempStyle.remove();
        }

        console.error("Erro ao gerar PDF:", error);
        notifications.update({
          id: "pdf-generation",
          title: "Erro",
          message: "Ocorreu um erro ao gerar o PDF. Tente novamente.",
          color: "red",
          icon: <IconX size={16} />,
          autoClose: 5000,
        });
      });
  };

  const teamStats =
    monthlyReports.length > 0
      ? {
          averageScore:
            monthlyReports.reduce(
              (sum, report) => sum + report.weightedAverageScore,
              0
            ) / monthlyReports.length,
          highestScore: Math.max(
            ...monthlyReports.map((report) => report.weightedAverageScore)
          ),
          lowestScore: Math.min(
            ...monthlyReports.map((report) => report.weightedAverageScore)
          ),
          totalEvaluations: monthlyReports.length,
        }
      : null;

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

      <Group justify="space-between" mb="xl">
        <div>
          <Title order={1} mb="xs">
            Relatório Consolidado
          </Title>
          <Text c="dimmed">
            Performance da equipe para apresentação à diretoria
          </Text>
        </div>
        {selectedMonth && tableData.length > 0 && (
          <Button
            leftSection={<IconFileTypePdf size={16} />}
            onClick={handleExportPDF}
          >
            Exportar PDF
          </Button>
        )}
      </Group>

      <Card shadow="sm" padding="lg" radius="md" withBorder mb="xl">
        <Group justify="space-between" mb="md">
          <Title order={3}>Selecionar Período</Title>
          <Text size="sm" c="dimmed">
            {availableMonths.length} meses disponíveis
          </Text>
        </Group>

        <Select
          placeholder="Selecione o mês para visualizar o relatório"
          data={availableMonths}
          value={selectedMonth}
          onChange={setSelectedMonth}
          size="md"
          clearable
        />
      </Card>

      {selectedMonth && (
        <div id="consolidated-report-content">
          {teamStats && (
            <Card shadow="sm" padding="lg" radius="md" withBorder mb="xl">
              <Title order={3} mb="md">
                Estatísticas da Equipe
              </Title>
              <Group grow>
                <Paper p="md" withBorder>
                  <Text size="sm" c="dimmed" mb="xs">
                    Performance Média
                  </Text>
                  <Text size="xl" fw={700} c="blue">
                    {teamStats.averageScore.toFixed(1)}/10
                  </Text>
                </Paper>
                <Paper p="md" withBorder>
                  <Text size="sm" c="dimmed" mb="xs">
                    Maior Nota
                  </Text>
                  <Text size="xl" fw={700} c="green">
                    {teamStats.highestScore.toFixed(1)}/10
                  </Text>
                </Paper>
                <Paper p="md" withBorder>
                  <Text size="sm" c="dimmed" mb="xs">
                    Menor Nota
                  </Text>
                  <Text size="xl" fw={700} c="orange">
                    {teamStats.lowestScore.toFixed(1)}/10
                  </Text>
                </Paper>
                <Paper p="md" withBorder>
                  <Text size="sm" c="dimmed" mb="xs">
                    Total de Avaliações
                  </Text>
                  <Text size="xl" fw={700}>
                    {teamStats.totalEvaluations}
                  </Text>
                </Paper>
              </Group>
            </Card>
          )}

          <Card shadow="sm" padding="lg" radius="md" withBorder>
            <Group justify="space-between" mb="md">
              <Title order={3}>
                Performance da Equipe -{" "}
                {availableMonths.find((m) => m.value === selectedMonth)?.label}
              </Title>
              <Text size="sm" c="dimmed">
                {tableData.length} avaliações
              </Text>
            </Group>

            {tableData.length > 0 ? (
              <Table striped highlightOnHover>
                <Table.Thead>
                  <Table.Tr>
                    <Table.Th>Posição</Table.Th>
                    <Table.Th>Desenvolvedor</Table.Th>
                    <Table.Th>Cargo</Table.Th>
                    <Table.Th>Performance Final</Table.Th>
                    <Table.Th>Status</Table.Th>
                    {Object.values(EVALUATION_CATEGORIES).map((category) => (
                      <Table.Th key={category.label}>{category.label}</Table.Th>
                    ))}
                  </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                  {tableData.map((report, index) => (
                    <Table.Tr key={report.id}>
                      <Table.Td>
                        <Text fw={500}>#{index + 1}</Text>
                      </Table.Td>
                      <Table.Td>
                        <Text fw={500}>{report.developerName}</Text>
                      </Table.Td>
                      <Table.Td>
                        <Text c="dimmed">{report.developerRole}</Text>
                      </Table.Td>
                      <Table.Td>
                        <Text fw={700} c="blue">
                          {report.weightedAverageScore.toFixed(1)}/10
                        </Text>
                      </Table.Td>
                      <Table.Td>
                        <Badge
                          color={getPerformanceColor(
                            report.weightedAverageScore
                          )}
                          variant="light"
                        >
                          {getPerformanceLabel(report.weightedAverageScore)}
                        </Badge>
                      </Table.Td>
                      {Object.keys(EVALUATION_CATEGORIES).map((key) => (
                        <Table.Td key={key}>
                          <Text size="sm">
                            {report.categoryScores[key]?.toFixed(1) || "N/A"}
                          </Text>
                        </Table.Td>
                      ))}
                    </Table.Tr>
                  ))}
                </Table.Tbody>
              </Table>
            ) : (
              <Text c="dimmed" ta="center" py="xl">
                Nenhuma avaliação encontrada para este período
              </Text>
            )}
          </Card>

          {tableData.length > 0 && (
            <Card shadow="sm" padding="lg" radius="md" withBorder mt="xl">
              <Title order={3} mb="md">
                Observações Qualitativas
              </Title>
              <Stack gap="md">
                {tableData.map((report) => (
                  <Paper key={report.id} p="md" withBorder>
                    <Text fw={500} mb="xs">
                      {report.developerName}
                    </Text>
                    <Group grow align="flex-start">
                      <div>
                        <Text size="sm" fw={500} c="green" mb="xs">
                          Destaques:
                        </Text>
                        <Text size="sm" c="dimmed">
                          {report.highlights || "Nenhum destaque registrado"}
                        </Text>
                      </div>
                      <div>
                        <Text size="sm" fw={500} c="orange" mb="xs">
                          Pontos a Desenvolver:
                        </Text>
                        <Text size="sm" c="dimmed">
                          {report.pointsToDevelop || "Nenhum ponto registrado"}
                        </Text>
                      </div>
                    </Group>
                  </Paper>
                ))}
              </Stack>
            </Card>
          )}
        </div>
      )}

      {!selectedMonth && (
        <Card shadow="sm" padding="xl" radius="md" withBorder>
          <div style={{ textAlign: "center" }}>
            <IconDownload size={48} color="var(--mantine-color-gray-5)" />
            <Title order={3} mt="md" mb="xs">
              Selecione um Período
            </Title>
            <Text c="dimmed">
              Escolha um mês para visualizar o relatório consolidado da equipe
            </Text>
          </div>
        </Card>
      )}
    </Container>
  );
};

export default ConsolidatedReport;
