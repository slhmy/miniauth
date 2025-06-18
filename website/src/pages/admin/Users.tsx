import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'
import { Plus, Search, MoreHorizontal, Edit, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { AdminApi } from '@/api/api'
import UserForm from '@/components/admin/UserForm'
import type { HandlersAdminUserInfo, HandlersAdminCreateUserRequest, HandlersAdminUpdateUserRequest } from '@/api/api'
import { formatDistanceToNow } from 'date-fns'

// Initialize API client
const adminApi = new AdminApi(undefined, '/api')

export default function UsersPage() {
  const { t } = useTranslation()
  
  // Set document title for admin users page
  useDocumentTitle(t('admin.userManagement'))
  
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize] = useState(10)
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedUser, setSelectedUser] = useState<HandlersAdminUserInfo | null>(null)
  const [showCreateDialog, setShowCreateDialog] = useState(false)
  const [showEditDialog, setShowEditDialog] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [formError, setFormError] = useState<string>('')

  const queryClient = useQueryClient()

  // Fetch users
  const { 
    data: usersResponse, 
    isLoading, 
    error 
  } = useQuery({
    queryKey: ['admin-users', currentPage, pageSize],
    queryFn: async () => {
      const response = await adminApi.adminUsersGet(currentPage, pageSize)
      return response.data
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })

  // Delete user mutation
  const deleteMutation = useMutation({
    mutationFn: async (userId: number) => {
      const response = await adminApi.adminUsersIdDelete(userId)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      setShowDeleteDialog(false)
      setSelectedUser(null)
    },
  })

  // Create user mutation
  const createMutation = useMutation({
    mutationFn: async (userData: HandlersAdminCreateUserRequest) => {
      const response = await adminApi.adminUsersPost(userData)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      setShowCreateDialog(false)
      setFormError('')
    },
    onError: () => {
      setFormError(t('errors.somethingWentWrong'))
    },
  })

  // Update user mutation
  const updateMutation = useMutation({
    mutationFn: async ({ id, userData }: { id: number; userData: HandlersAdminUpdateUserRequest }) => {
      const response = await adminApi.adminUsersIdPut(id, userData)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      setShowEditDialog(false)
      setSelectedUser(null)
      setFormError('')
    },
    onError: () => {
      setFormError(t('errors.somethingWentWrong'))
    },
  })

  const handleEditUser = (user: HandlersAdminUserInfo) => {
    setSelectedUser(user)
    setFormError('')
    setShowEditDialog(true)
  }

  const handleDeleteUser = (user: HandlersAdminUserInfo) => {
    setSelectedUser(user)
    setShowDeleteDialog(true)
  }

  const handleCreateUser = () => {
    setFormError('')
    setShowCreateDialog(true)
  }

  const handleCreateSubmit = async (data: HandlersAdminCreateUserRequest | HandlersAdminUpdateUserRequest) => {
    await createMutation.mutateAsync(data as HandlersAdminCreateUserRequest)
  }

  const handleEditSubmit = async (data: HandlersAdminCreateUserRequest | HandlersAdminUpdateUserRequest) => {
    if (selectedUser?.id) {
      await updateMutation.mutateAsync({ 
        id: selectedUser.id, 
        userData: data as HandlersAdminUpdateUserRequest 
      })
    }
  }

  const confirmDelete = () => {
    if (selectedUser?.id) {
      deleteMutation.mutate(selectedUser.id)
    }
  }

  // Filter users based on search term
  const filteredUsers = usersResponse?.users?.filter((user: HandlersAdminUserInfo) =>
    user.username?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email?.toLowerCase().includes(searchTerm.toLowerCase())
  ) || []

  const totalUsers = usersResponse?.total || 0
  const totalPages = Math.ceil(totalUsers / pageSize)

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <h3 className="text-lg font-medium text-foreground mb-2">{t('errors.somethingWentWrong')}</h3>
          <p className="text-muted-foreground">{t('errors.networkError')}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-foreground">{t('admin.userManagement')}</h1>
          <p className="text-muted-foreground">{t('admin.userManagementDesc')}</p>
        </div>
        <Button onClick={handleCreateUser} className="w-full sm:w-auto">
          <Plus className="mr-2 h-4 w-4" />
          {t('admin.createUser')}
        </Button>
      </div>

      {/* Search and Filters */}
      <Card>
        <CardHeader>
          <CardTitle>{t('common.users')}</CardTitle>
          <CardDescription>
            {totalUsers} {t('admin.totalUsers')}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col sm:flex-row gap-4 mb-6">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
              <Input
                placeholder={t('admin.searchUsers')}
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          {/* Users Table */}
          <div className="border rounded-lg overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t('common.user')}</TableHead>
                  <TableHead className="hidden sm:table-cell">{t('common.email')}</TableHead>
                  <TableHead className="hidden md:table-cell">{t('common.organizationCount')}</TableHead>
                  <TableHead className="hidden lg:table-cell">{t('common.created')}</TableHead>
                  <TableHead className="w-[70px]">{t('admin.actions')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  // Loading skeleton
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell>
                        <div className="flex items-center space-x-3">
                          <div className="h-8 w-8 bg-gray-200 rounded-full animate-pulse" />
                          <div className="space-y-1">
                            <div className="h-4 w-24 bg-gray-200 rounded animate-pulse" />
                            <div className="h-3 w-32 bg-gray-200 rounded animate-pulse sm:hidden" />
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="hidden sm:table-cell">
                        <div className="h-4 w-32 bg-gray-200 rounded animate-pulse" />
                      </TableCell>
                      <TableCell className="hidden md:table-cell">
                        <div className="h-4 w-16 bg-gray-200 rounded animate-pulse" />
                      </TableCell>
                      <TableCell className="hidden lg:table-cell">
                        <div className="h-4 w-20 bg-gray-200 rounded animate-pulse" />
                      </TableCell>
                      <TableCell>
                        <div className="h-8 w-8 bg-gray-200 rounded animate-pulse" />
                      </TableCell>
                    </TableRow>
                  ))
                ) : filteredUsers.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center py-8">
                      <div className="text-muted-foreground">
                        {searchTerm ? t('admin.noUsersFound') : t('admin.noUsersFound')}
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredUsers.map((user: HandlersAdminUserInfo) => (
                    <TableRow key={user.id}>
                      <TableCell>
                        <div className="flex items-center space-x-3">
                          <div className="h-8 w-8 bg-blue-100 rounded-full flex items-center justify-center">
                            <span className="text-blue-600 font-medium text-sm">
                              {user.username?.charAt(0).toUpperCase()}
                            </span>
                          </div>
                          <div>
                            <div className="font-medium text-foreground">{user.username}</div>
                            <div className="text-sm text-muted-foreground sm:hidden">{user.email}</div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="hidden sm:table-cell text-muted-foreground">
                        {user.email}
                      </TableCell>
                      <TableCell className="hidden md:table-cell">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400">
                          {user.org_count || 0}
                        </span>
                      </TableCell>
                      <TableCell className="hidden lg:table-cell text-muted-foreground text-sm">
                        {user.created_at ? formatDistanceToNow(new Date(user.created_at), { addSuffix: true }) : t('admin.never')}
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" className="h-8 w-8 p-0">
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => handleEditUser(user)}>
                              <Edit className="mr-2 h-4 w-4" />
                              {t('common.edit')}
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => handleDeleteUser(user)} className="text-red-600">
                              <Trash2 className="mr-2 h-4 w-4" />
                              {t('common.delete')}
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between mt-6">
              <div className="text-sm text-muted-foreground">
                {t('common.showing', {
                  from: (currentPage - 1) * pageSize + 1,
                  to: Math.min(currentPage * pageSize, totalUsers),
                  total: totalUsers
                })}
              </div>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                >
                  {t('common.previous')}
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                  disabled={currentPage === totalPages}
                >
                  {t('common.next')}
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t('admin.deleteUser')}</DialogTitle>
            <DialogDescription>
              {t('admin.confirmDelete')} <strong>{selectedUser?.username}</strong>?
              <br />
              <span className="text-sm text-muted-foreground mt-2">
                {t('admin.deleteUserWarning')}
              </span>
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-end gap-3 mt-6">
            <Button
              variant="outline"
              onClick={() => setShowDeleteDialog(false)}
            >
              {t('common.cancel')}
            </Button>
            <Button
              variant="destructive"
              onClick={confirmDelete}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? t('common.loading') : t('common.delete')}
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      {/* Create User Dialog */}
      <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{t('admin.createUser')}</DialogTitle>
            <DialogDescription>
              {t('forms.fillDetails')}
            </DialogDescription>
          </DialogHeader>
          <UserForm
            isCreateMode={true}
            onSubmit={handleCreateSubmit}
            onCancel={() => setShowCreateDialog(false)}
            isLoading={createMutation.isPending}
            error={formError}
          />
        </DialogContent>
      </Dialog>

      {/* Edit User Dialog */}
      <Dialog open={showEditDialog} onOpenChange={setShowEditDialog}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{t('admin.editUser')}: {selectedUser?.username}</DialogTitle>
            <DialogDescription>
              {t('forms.updateUserInfo')}
            </DialogDescription>
          </DialogHeader>
          <UserForm
            user={selectedUser || undefined}
            isCreateMode={false}
            onSubmit={handleEditSubmit}
            onCancel={() => setShowEditDialog(false)}
            isLoading={updateMutation.isPending}
            error={formError}
          />
        </DialogContent>
      </Dialog>
    </div>
  )
}
