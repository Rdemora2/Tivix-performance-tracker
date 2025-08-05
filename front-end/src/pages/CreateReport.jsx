import {
  Container,
  Title,
  Card,
  Text,
  Group,
  Button,
  Stepper,
  Slider,
  Textarea,
  Select,
  Divider,
  Stack,
  Badge,
} from "@mantine/core";
import { IconArrowLeft, IconCheck } from "@tabler/icons-react";
import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import useAppStore from "../store/useAppStore";
import {
  createPerformanceReport,
  EVALUATION_CATEGORIES,
  getAllQuestions,
} from "../types";

const CreateReport = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { developers, addPerformanceReport } = useAppStore();
  const [active, setActive] = useState(0);

  const developer = developers.find((dev) => dev.id === id);
  const allQuestions = getAllQuestions();

  const initialQuestionValues = {};
  allQuestions.forEach((question) => {
    initialQuestionValues[question.key] = 5;
  });

  const form = useForm({
    initialValues: {
      month: new Date()
        .toLocaleDateString("en-CA", { year: "numeric", month: "2-digit" })
        .slice(0, 7),
      ...initialQuestionValues,
      highlights: "",
      pointsToDevelop: "",
    },
    validate: {
      month: (value) => (!value ? "Mês é obrigatório" : null),
      highlights: (value) =>
        value.length < 10
          ? "Destaques devem ter pelo menos 10 caracteres"
          : null,
      pointsToDevelop: (value) =>
        value.length < 10
          ? "Pontos a desenvolver devem ter pelo menos 10 caracteres"
          : null,
    },
  });

  if (!developer) {
    return (
      <Container size="xl">
        <Text>Desenvolvedor não encontrado</Text>
      </Container>
    );
  }

  const nextStep = () => {
    if (active === 2) {
      const errors = form.validate();
      if (errors.hasErrors) {
        return;
      }
    }
    setActive((current) => (current < 3 ? current + 1 : current));
  };

  const prevStep = () =>
    setActive((current) => (current > 0 ? current - 1 : current));

  const handleSubmit = async (values) => {
    try {
      const questionScores = {};
      allQuestions.forEach((question) => {
        questionScores[question.key] = values[question.key];
      });

      const reportCalculation = createPerformanceReport(
        "",
        id,
        values.month,
        questionScores,
        values.highlights,
        values.pointsToDevelop
      );

      const reportData = {
        developerId: id,
        month: values.month,
        questionScores,
        categoryScores: reportCalculation.categoryScores,
        weightedAverageScore: reportCalculation.weightedAverageScore,
        highlights: values.highlights,
        pointsToDevelop: values.pointsToDevelop,
      };

      await addPerformanceReport(reportData);

      notifications.show({
        title: "Sucesso",
        message: "Relatório de performance criado com sucesso!",
        color: "green",
      });

      navigate(`/developer/${id}`);
    } catch (error) {
      console.error("Erro ao salvar relatório:", error);
      notifications.show({
        title: "Erro",
        message: "Ocorreu um erro ao salvar o relatório. Tente novamente.",
        color: "red",
      });
    }
  };

  const getSliderColor = (value) => {
    if (value >= 8) return "green";
    if (value >= 6) return "blue";
    if (value >= 4) return "yellow";
    return "red";
  };

  const getSliderLabel = (value) => {
    if (value >= 8) return "Excelente";
    if (value >= 6) return "Bom";
    if (value >= 4) return "Regular";
    return "Precisa Melhorar";
  };

  const monthOptions = Array.from({ length: 12 }, (_, i) => {
    const date = new Date();
    date.setMonth(date.getMonth() - i);
    const value = date.toISOString().slice(0, 7);
    const label = date.toLocaleDateString("pt-BR", {
      month: "long",
      year: "numeric",
    });
    return { value, label };
  });

  const calculatePreviewScores = () => {
    const questionScores = {};
    allQuestions.forEach((question) => {
      questionScores[question.key] = form.values[question.key];
    });

    return createPerformanceReport("", "", "", questionScores);
  };

  return (
    <Container size="lg">
      <Group mb="xl">
        <Button
          variant="subtle"
          leftSection={<IconArrowLeft size={16} />}
          onClick={() => navigate(`/developer/${id}`)}
        >
          Voltar ao Perfil
        </Button>
      </Group>

      <Card shadow="sm" padding="lg" radius="md" withBorder>
        <Title order={2} mb="xs">
          Nova Avaliação de Performance
        </Title>
        <Text c="dimmed" mb="xl">
          Avaliando: <strong>{developer.name}</strong> - {developer.role}
        </Text>

        <form onSubmit={form.onSubmit(handleSubmit)}>
          <Stepper
            active={active}
            onStepClick={setActive}
            breakpoint="sm"
            mb="xl"
          >
            <Stepper.Step label="Período" description="Selecione o mês">
              <Select
                label="Mês de Avaliação"
                placeholder="Selecione o mês"
                data={monthOptions}
                {...form.getInputProps("month")}
                mb="xl"
              />
            </Stepper.Step>

            <Stepper.Step label="Avaliações" description="Pontue cada pergunta">
              <Stack gap="xl">
                {Object.entries(EVALUATION_CATEGORIES).map(
                  ([categoryKey, category]) => (
                    <Card key={categoryKey} withBorder p="md">
                      <Group justify="space-between" mb="md">
                        <Title order={4}>{category.label}</Title>
                        <Badge variant="light" color="blue">
                          Peso: {(category.weight * 100).toFixed(0)}%
                        </Badge>
                      </Group>

                      <Stack gap="lg">
                        {category.questions.map((question) => (
                          <div key={question.key}>
                            <Group justify="space-between" mb="xs">
                              <div>
                                <Text fw={500} size="sm">
                                  {question.label}
                                </Text>
                                <Text size="xs" c="dimmed">
                                  Peso na categoria: {question.weight}
                                </Text>
                              </div>
                              <Group gap="xs">
                                <Text
                                  size="sm"
                                  c={getSliderColor(form.values[question.key])}
                                >
                                  {form.values[question.key]}/10
                                </Text>
                                <Text size="xs" c="dimmed">
                                  ({getSliderLabel(form.values[question.key])})
                                </Text>
                              </Group>
                            </Group>
                            <Slider
                              {...form.getInputProps(question.key)}
                              min={0}
                              max={10}
                              step={0.5}
                              marks={[
                                { value: 0, label: "0" },
                                { value: 2.5, label: "2.5" },
                                { value: 5, label: "5" },
                                { value: 7.5, label: "7.5" },
                                { value: 10, label: "10" },
                              ]}
                              color={getSliderColor(form.values[question.key])}
                            />
                          </div>
                        ))}
                      </Stack>
                    </Card>
                  )
                )}
              </Stack>
            </Stepper.Step>

            <Stepper.Step
              label="Comentários"
              description="Feedback qualitativo"
            >
              <Textarea
                label="Destaques do Mês"
                placeholder="Descreva os principais pontos positivos e conquistas do desenvolvedor neste período..."
                {...form.getInputProps("highlights")}
                minRows={4}
                mb="md"
              />

              <Textarea
                label="Pontos a Desenvolver"
                placeholder="Identifique áreas de melhoria e sugestões para o desenvolvimento profissional..."
                {...form.getInputProps("pointsToDevelop")}
                minRows={4}
              />
            </Stepper.Step>

            <Stepper.Completed>
              <div style={{ textAlign: "center", padding: "2rem" }}>
                <IconCheck size={48} color="var(--mantine-color-green-6)" />
                <Title order={3} mt="md" mb="xs">
                  Avaliação Concluída
                </Title>
                <Text c="dimmed" mb="lg">
                  Revise as informações e clique em "Salvar Relatório" para
                  finalizar.
                </Text>

                <Card withBorder p="md" mb="lg">
                  <Title order={4} mb="md">
                    Resumo da Avaliação
                  </Title>

                  <Stack gap="sm" mb="md">
                    {Object.entries(EVALUATION_CATEGORIES).map(
                      ([categoryKey, category]) => {
                        const previewScores = calculatePreviewScores();
                        const categoryScore =
                          previewScores.categoryScores[categoryKey];

                        return (
                          <Group key={categoryKey} justify="space-between">
                            <Text size="sm">{category.label}:</Text>
                            <Text size="sm" fw={700} c="blue">
                              {categoryScore.toFixed(1)}/10
                            </Text>
                          </Group>
                        );
                      }
                    )}
                  </Stack>

                  <Divider mb="md" />

                  <Group justify="space-between">
                    <Text fw={500}>Performance Final:</Text>
                    <Text fw={700} c="blue" size="lg">
                      {calculatePreviewScores().weightedAverageScore.toFixed(1)}
                      /10
                    </Text>
                  </Group>
                </Card>
              </div>
            </Stepper.Completed>
          </Stepper>

          <Group justify="space-between" mt="xl">
            <Button
              variant="outline"
              onClick={prevStep}
              disabled={active === 0}
            >
              Anterior
            </Button>

            {active < 3 ? (
              <Button onClick={nextStep}>Próximo</Button>
            ) : (
              <Button type="submit" leftSection={<IconCheck size={16} />}>
                Salvar Relatório
              </Button>
            )}
          </Group>
        </form>
      </Card>
    </Container>
  );
};

export default CreateReport;
