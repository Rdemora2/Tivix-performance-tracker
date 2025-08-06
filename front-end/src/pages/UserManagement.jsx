import React, { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Plus,
  Edit2,
  Eye,
  EyeOff,
  Copy,
  RefreshCw,
  Users,
  Building,
  AlertTriangle,
  Trash2,
} from "lucide-react";
import useAppStore from "@/store/useAppStore";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { toast } from "sonner";

const UserManagement = () => {
  const {
    user,
    users,
    companies,
    fetchUsers,
    fetchCompanies,
    createUser,
    updateUser,
    deleteUser,
  } = useAppStore();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [userToDelete, setUserToDelete] = useState(null);
  const [loading, setLoading] = useState(false);
  const [loadingUsers, setLoadingUsers] = useState(true);
  const [error, setError] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [generatedPassword, setGeneratedPassword] = useState("");

  const [formData, setFormData] = useState({
    name: "",
    email: "",
    role: "user",
    companyId: "",
  });

  // Carregar dados iniciais
  useEffect(() => {
    const loadData = async () => {
      setLoadingUsers(true);
      try {
        const [usersData, companiesData] = await Promise.all([
          fetchUsers(),
          fetchCompanies(),
        ]);
        console.log("Loaded users:", usersData);
        console.log("Loaded companies:", companiesData);
      } catch {
        setError("Erro ao carregar dados");
      } finally {
        setLoadingUsers(false);
      }
    };

    if (user?.role === "admin" || user?.role === "manager") {
      loadData();
    } else {
      setLoadingUsers(false);
    }
  }, [user, fetchUsers, fetchCompanies]);

  // Automaticamente definir empresa para managers
  useEffect(() => {
    if (user?.role === "manager" && user?.companyId && !formData.companyId) {
      setFormData((prev) => ({
        ...prev,
        companyId: user.companyId,
      }));
    }
  }, [user, formData.companyId]);

  const generateRandomPassword = () => {
    const chars =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%&*";
    let password = "";
    for (let i = 0; i < 12; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return password;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      if (editingUser) {
        // Editar usuário existente
        await updateUser(editingUser.id, formData);
        toast.success("Usuário atualizado com sucesso!");
        setIsEditDialogOpen(false);
        setEditingUser(null);
      } else {
        // Criar novo usuário
        const tempPassword = generateRandomPassword();
        setGeneratedPassword(tempPassword);

        await createUser({
          ...formData,
          temporaryPassword: tempPassword,
        });

        toast.success("Usuário criado com sucesso!");
        // Reset form but keep modal open to show password
        setFormData({ name: "", email: "", role: "user", companyId: "" });
      }
    } catch (err) {
      setError(err.message || "Erro ao salvar usuário");
      toast.error(err.message || "Erro ao salvar usuário");
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (userToEdit) => {
    setEditingUser(userToEdit);
    setFormData({
      name: userToEdit.name,
      email: userToEdit.email,
      role: userToEdit.role,
      companyId: userToEdit.companyId || userToEdit.company_id,
    });
    setIsEditDialogOpen(true);
  };

  const handleDeleteClick = (userToDelete) => {
    setUserToDelete(userToDelete);
    setIsDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!userToDelete) return;

    setLoading(true);
    try {
      await deleteUser(userToDelete.id);
      toast.success(`Usuário ${userToDelete.name} excluído com sucesso!`);
      setIsDeleteDialogOpen(false);
      setUserToDelete(null);
    } catch (err) {
      toast.error(err.message || "Erro ao excluir usuário");
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    toast.success("Senha copiada para a área de transferência!");
  };

  const refreshUsers = async () => {
    setLoadingUsers(true);
    try {
      await fetchUsers();
      toast.success("Lista de usuários atualizada!");
    } catch {
      toast.error("Erro ao recarregar usuários");
    } finally {
      setLoadingUsers(false);
    }
  };

  const resetForm = () => {
    const defaultCompanyId =
      user?.role === "manager" && user?.companyId ? user.companyId : "";
    setFormData({
      name: "",
      email: "",
      role: "user",
      companyId: defaultCompanyId,
    });
    setError("");
    setEditingUser(null);
    setGeneratedPassword("");
    setShowPassword(false);
  };

  const getRoleColor = (role) => {
    switch (role) {
      case "admin":
        return "destructive";
      case "manager":
        return "default";
      case "user":
        return "secondary";
      default:
        return "outline";
    }
  };

  const getRoleLabel = (role) => {
    switch (role) {
      case "admin":
        return "Administrador";
      case "manager":
        return "Gerente";
      case "user":
        return "Usuário";
      default:
        return role;
    }
  };

  const getCompanyName = (companyId) => {
    const company = companies.find((c) => c.id === companyId);
    return company?.name || "N/A";
  };

  const formatDate = (dateString) => {
    if (!dateString) return "N/A";
    return new Date(dateString).toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    });
  };

  // Filtrar empresas baseado na permissão do usuário
  const availableCompanies = () => {
    if (user?.role === "admin") {
      return companies;
    } else if (user?.role === "manager") {
      return companies.filter((company) => company.id === user.companyId);
    }
    return [];
  };

  // Filtrar usuários baseado na permissão
  const filteredUsers = () => {
    if (user?.role === "admin") {
      return users;
    } else if (user?.role === "manager") {
      return users.filter(
        (u) => (u.companyId || u.company_id) === user.companyId
      );
    }
    return [];
  };

  if (!user || (user.role !== "admin" && user.role !== "manager")) {
    return (
      <div className="container mx-auto py-10">
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            Apenas administradores e gerentes podem gerenciar usuários.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            Gerenciamento de Usuários
          </h1>
          <p className="text-muted-foreground">
            Gerencie usuários{" "}
            {user?.role === "manager" ? "da sua empresa" : "do sistema"}
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={refreshUsers}
            disabled={loadingUsers}
          >
            <RefreshCw
              className={`h-4 w-4 mr-2 ${loadingUsers ? "animate-spin" : ""}`}
            />
            Atualizar
          </Button>
          <Dialog
            open={isCreateDialogOpen}
            onOpenChange={(open) => {
              setIsCreateDialogOpen(open);
              if (!open) resetForm();
            }}
          >
            <DialogTrigger asChild>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                Novo Usuário
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[500px]">
              <DialogHeader>
                <DialogTitle>Criar Novo Usuário</DialogTitle>
              </DialogHeader>

              {generatedPassword ? (
                <div className="space-y-4">
                  <Alert>
                    <AlertTriangle className="h-4 w-4" />
                    <AlertDescription>
                      Usuário criado com sucesso! Uma senha temporária foi
                      gerada. O usuário deverá definir uma nova senha no
                      primeiro login.
                    </AlertDescription>
                  </Alert>

                  <Card>
                    <CardHeader>
                      <CardTitle className="text-sm">
                        Senha Temporária
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                      <div className="flex items-center gap-2">
                        <Input
                          type={showPassword ? "text" : "password"}
                          value={generatedPassword}
                          readOnly
                          className="font-mono"
                        />
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={() => setShowPassword(!showPassword)}
                        >
                          {showPassword ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={() => copyToClipboard(generatedPassword)}
                        >
                          <Copy className="h-4 w-4" />
                        </Button>
                      </div>
                      <p className="text-xs text-muted-foreground">
                        Copie esta senha e envie para o usuário de forma segura.
                      </p>
                    </CardContent>
                  </Card>

                  <Button
                    onClick={() => setIsCreateDialogOpen(false)}
                    className="w-full"
                  >
                    Fechar
                  </Button>
                </div>
              ) : (
                <form onSubmit={handleSubmit} className="space-y-4">
                  {error && (
                    <Alert variant="destructive">
                      <AlertTriangle className="h-4 w-4" />
                      <AlertDescription>{error}</AlertDescription>
                    </Alert>
                  )}

                  <div>
                    <Label htmlFor="name" className="mb-2 block">
                      Nome completo *
                    </Label>
                    <Input
                      id="name"
                      value={formData.name}
                      onChange={(e) =>
                        setFormData({ ...formData, name: e.target.value })
                      }
                      required
                      placeholder="Nome do usuário"
                    />
                  </div>

                  <div>
                    <Label htmlFor="email" className="mb-2 block">
                      Email *
                    </Label>
                    <Input
                      id="email"
                      type="email"
                      value={formData.email}
                      onChange={(e) =>
                        setFormData({ ...formData, email: e.target.value })
                      }
                      required
                      placeholder="email@empresa.com"
                    />
                  </div>

                  {/* Mostrar select de empresa apenas para admin ou quando há mais de uma empresa disponível */}
                  {(user?.role === "admin" ||
                    availableCompanies().length > 1) && (
                    <div>
                      <Label htmlFor="company" className="mb-2 block">
                        Empresa *
                      </Label>
                      <Select
                        value={formData.companyId}
                        onValueChange={(value) =>
                          setFormData({ ...formData, companyId: value })
                        }
                        required
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Selecione uma empresa" />
                        </SelectTrigger>
                        <SelectContent>
                          {availableCompanies().map((company) => (
                            <SelectItem key={company.id} value={company.id}>
                              {company.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  )}

                  {/* Mostrar empresa selecionada quando há apenas uma opção */}
                  {user?.role === "manager" &&
                    availableCompanies().length === 1 && (
                      <div>
                        <Label className="mb-2 block">Empresa</Label>
                        <div className="flex items-center gap-2 p-3 border rounded-md bg-muted/50">
                          <Building className="h-4 w-4 text-muted-foreground" />
                          <span className="text-sm font-medium">
                            {availableCompanies()[0]?.name}
                          </span>
                        </div>
                      </div>
                    )}

                  <div>
                    <Label htmlFor="role" className="mb-2 block">
                      Função *
                    </Label>
                    <Select
                      value={formData.role}
                      onValueChange={(value) =>
                        setFormData({ ...formData, role: value })
                      }
                      required
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Selecione uma função" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="user">Usuário</SelectItem>
                        <SelectItem value="manager">Gerente</SelectItem>
                        {user?.role === "admin" && (
                          <SelectItem value="admin">Administrador</SelectItem>
                        )}
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="flex justify-end gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => setIsCreateDialogOpen(false)}
                      disabled={loading}
                    >
                      Cancelar
                    </Button>
                    <Button type="submit" disabled={loading}>
                      {loading ? "Criando..." : "Criar Usuário"}
                    </Button>
                  </div>
                </form>
              )}
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {error && !isCreateDialogOpen && !isEditDialogOpen && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            Usuários Cadastrados ({filteredUsers().length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          {loadingUsers ? (
            <div className="text-center py-8">
              <p className="text-muted-foreground">Carregando usuários...</p>
            </div>
          ) : filteredUsers().length === 0 ? (
            <div className="text-center py-8">
              <Users className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
              <p className="text-muted-foreground">Nenhum usuário encontrado</p>
              <p className="text-sm text-muted-foreground">
                Clique em "Novo Usuário" para criar o primeiro usuário
              </p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Nome</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Empresa</TableHead>
                  <TableHead>Função</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Criado em</TableHead>
                  <TableHead className="text-right">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredUsers().map((userItem) => (
                  <TableRow key={userItem.id}>
                    <TableCell className="font-medium">
                      {userItem.name}
                    </TableCell>
                    <TableCell>{userItem.email}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Building className="h-4 w-4 text-muted-foreground" />
                        {getCompanyName(
                          userItem.companyId || userItem.company_id
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant={getRoleColor(userItem.role)}>
                        {getRoleLabel(userItem.role)}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {userItem.needsPasswordChange ? (
                        <Badge variant="outline" className="text-yellow-600">
                          Aguardando nova senha
                        </Badge>
                      ) : userItem.isActive !== false ? (
                        <Badge variant="outline" className="text-green-600">
                          Ativo
                        </Badge>
                      ) : (
                        <Badge variant="secondary">Inativo</Badge>
                      )}
                    </TableCell>
                    <TableCell>{formatDate(userItem.createdAt)}</TableCell>
                    <TableCell className="text-right">
                      <div className="flex gap-2 justify-end">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleEdit(userItem)}
                        >
                          <Edit2 className="h-4 w-4" />
                        </Button>
                        {(user?.role === "admin" || 
                          (user?.role === "manager" && 
                           userItem.role !== "admin" && 
                           userItem.role !== "manager" &&
                           userItem.id !== user.id)) && (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDeleteClick(userItem)}
                            className="text-destructive hover:text-destructive"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Dialog de edição */}
      <Dialog
        open={isEditDialogOpen}
        onOpenChange={(open) => {
          setIsEditDialogOpen(open);
          if (!open) resetForm();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Editar Usuário</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertTriangle className="h-4 w-4" />
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div>
              <Label htmlFor="edit-name" className="mb-2 block">
                Nome completo *
              </Label>
              <Input
                id="edit-name"
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                required
                placeholder="Nome do usuário"
              />
            </div>

            <div>
              <Label htmlFor="edit-email" className="mb-2 block">
                Email *
              </Label>
              <Input
                id="edit-email"
                type="email"
                value={formData.email}
                onChange={(e) =>
                  setFormData({ ...formData, email: e.target.value })
                }
                required
                placeholder="email@empresa.com"
              />
            </div>

            {/* Mostrar select de empresa apenas para admin ou quando há mais de uma empresa disponível */}
            {(user?.role === "admin" || availableCompanies().length > 1) && (
              <div>
                <Label htmlFor="edit-company" className="mb-2 block">
                  Empresa *
                </Label>
                <Select
                  value={formData.companyId}
                  onValueChange={(value) =>
                    setFormData({ ...formData, companyId: value })
                  }
                  required
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione uma empresa" />
                  </SelectTrigger>
                  <SelectContent>
                    {availableCompanies().map((company) => (
                      <SelectItem key={company.id} value={company.id}>
                        {company.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}

            {/* Mostrar empresa selecionada quando há apenas uma opção */}
            {user?.role === "manager" && availableCompanies().length === 1 && (
              <div>
                <Label className="mb-2 block">Empresa</Label>
                <div className="flex items-center gap-2 p-3 border rounded-md bg-muted/50">
                  <Building className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">
                    {availableCompanies()[0]?.name}
                  </span>
                </div>
              </div>
            )}

            <div>
              <Label htmlFor="edit-role" className="mb-2 block">
                Função *
              </Label>
              <Select
                value={formData.role}
                onValueChange={(value) =>
                  setFormData({ ...formData, role: value })
                }
                required
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione uma função" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="user">Usuário</SelectItem>
                  <SelectItem value="manager">Gerente</SelectItem>
                  {user?.role === "admin" && (
                    <SelectItem value="admin">Administrador</SelectItem>
                  )}
                </SelectContent>
              </Select>
            </div>

            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsEditDialogOpen(false)}
                disabled={loading}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={loading}>
                {loading ? "Salvando..." : "Salvar Alterações"}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Dialog de confirmação de exclusão */}
      <AlertDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Confirmar Exclusão</AlertDialogTitle>
            <AlertDialogDescription>
              Tem certeza que deseja excluir o usuário <strong>{userToDelete?.name}</strong>?
              <br />
              <br />
              <span className="text-destructive">
                ⚠️ Esta ação não pode ser desfeita e removerá permanentemente:
              </span>
              <ul className="mt-2 ml-4 text-sm text-muted-foreground">
                <li>• Todos os dados do usuário</li>
                <li>• Histórico de atividades</li>
                <li>• Acesso ao sistema</li>
              </ul>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={loading}>
              Cancelar
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteConfirm}
              disabled={loading}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {loading ? "Excluindo..." : "Excluir"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
};

export default UserManagement;
