import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useDocumentTitle } from '@/hooks/useDocumentTitle';
import { Plus, Settings, Trash2, Eye, EyeOff, Copy, Pencil } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '../components/ui/dialog';
import { Input } from '../components/ui/input';
import { Label } from '../components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Badge } from '../components/ui/badge';

interface OAuthApplication {
  id: number;
  name: string;
  client_id: string;
  client_secret: string;
  redirect_uris: string[];
  scopes: string[];
  trusted: boolean;
  active: boolean;
  created_at: string;
}

interface ApplicationFormData {
  name: string;
  redirect_uris: string[];
  scopes: string[];
  trusted: boolean;
}

export default function OAuthApplications() {
  const { t } = useTranslation();
  
  // Set document title for OAuth applications page
  useDocumentTitle(t('oauth.applications'));

  const [applications, setApplications] = useState<OAuthApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [selectedApp, setSelectedApp] = useState<OAuthApplication | null>(null);
  const [editingApp, setEditingApp] = useState<OAuthApplication | null>(null);
  const [showSecrets, setShowSecrets] = useState<{[key: number]: boolean}>({});
  const [formData, setFormData] = useState<ApplicationFormData>({
    name: '',
    redirect_uris: [''],
    scopes: ['read'],
    trusted: false,
  });

  const availableScopes = [
    { value: 'read', label: t('oauth.authorization.scopeDescriptions.read') },
    { value: 'write', label: t('oauth.authorization.scopeDescriptions.write') },
    { value: 'profile', label: t('oauth.authorization.scopeDescriptions.profile') },
    { value: 'organizations', label: t('oauth.authorization.scopeDescriptions.organizations') },
    { value: 'admin', label: t('oauth.authorization.scopeDescriptions.admin') },
  ];

  const fetchApplications = async () => {
    try {
      const response = await fetch('/api/admin/oauth/applications', {
        credentials: 'include'
      });
      if (response.ok) {
        const data = await response.json();
        setApplications(data);
      }
    } catch (error) {
      console.error('Failed to fetch applications:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApplications();
  }, []);

  const handleCreateApplication = async () => {
    try {
      const response = await fetch('/api/admin/oauth/applications', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          name: formData.name,
          redirect_uris: formData.redirect_uris.filter(uri => uri.trim() !== ''),
          scopes: formData.scopes,
          trusted: formData.trusted,
        }),
      });
      
      if (response.ok) {
        const newApp = await response.json();
        setApplications([...applications, newApp]);
        setCreateDialogOpen(false);
        setFormData({
          name: '',
          redirect_uris: [''],
          scopes: ['read'],
          trusted: false,
        });
      }
    } catch (error) {
      console.error('Failed to create application:', error);
    }
  };

  const handleDeleteApplication = async () => {
    if (!selectedApp) return;

    try {
      const response = await fetch(`/api/admin/oauth/applications/${selectedApp.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      
      if (response.ok) {
        setApplications(applications.filter(app => app.id !== selectedApp.id));
        setDeleteDialogOpen(false);
        setSelectedApp(null);
      }
    } catch (error) {
      console.error('Failed to delete application:', error);
    }
  };

  const handleUpdateApplication = async () => {
    if (!editingApp) return;

    try {
      const response = await fetch(`/api/admin/oauth/applications/${editingApp.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          name: formData.name,
          redirect_uris: formData.redirect_uris.filter(uri => uri.trim() !== ''),
          scopes: formData.scopes,
          trusted: formData.trusted,
        }),
      });
      
      if (response.ok) {
        const updatedApp = await response.json();
        setApplications(applications.map(app => 
          app.id === editingApp.id ? updatedApp : app
        ));
        setEditDialogOpen(false);
        setEditingApp(null);
        setFormData({
          name: '',
          redirect_uris: [''],
          scopes: ['read'],
          trusted: false,
        });
      }
    } catch (error) {
      console.error('Failed to update application:', error);
    }
  };

  const handleToggleApplicationStatus = async (appId: number) => {
    try {
      const response = await fetch(`/api/admin/oauth/applications/${appId}/toggle`, {
        method: 'POST',
        credentials: 'include',
      });
      
      if (response.ok) {
        // Refresh the applications list to get updated status
        fetchApplications();
      }
    } catch (error) {
      console.error('Failed to toggle application status:', error);
    }
  };

  const handleToggleApplicationTrustedStatus = async (appId: number) => {
    try {
      const response = await fetch(`/api/admin/oauth/applications/${appId}/toggle-trusted`, {
        method: 'POST',
        credentials: 'include',
      });
      
      if (response.ok) {
        // Refresh the applications list to get updated status
        fetchApplications();
      }
    } catch (error) {
      console.error('Failed to toggle application trusted status:', error);
    }
  };

  const toggleSecretVisibility = (appId: number) => {
    setShowSecrets(prev => ({
      ...prev,
      [appId]: !prev[appId]
    }));
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const addRedirectUri = () => {
    setFormData(prev => ({
      ...prev,
      redirect_uris: [...prev.redirect_uris, '']
    }));
  };

  const updateRedirectUri = (index: number, value: string) => {
    setFormData(prev => ({
      ...prev,
      redirect_uris: prev.redirect_uris.map((uri, i) => i === index ? value : uri)
    }));
  };

  const removeRedirectUri = (index: number) => {
    setFormData(prev => ({
      ...prev,
      redirect_uris: prev.redirect_uris.filter((_, i) => i !== index)
    }));
  };

  const toggleScope = (scope: string) => {
    setFormData(prev => ({
      ...prev,
      scopes: prev.scopes.includes(scope)
        ? prev.scopes.filter(s => s !== scope)
        : [...prev.scopes, scope]
    }));
  };

  const openEditDialog = (app: OAuthApplication) => {
    setEditingApp(app);
    setFormData({
      name: app.name,
      redirect_uris: app.redirect_uris,
      scopes: app.scopes,
      trusted: app.trusted,
    });
    setEditDialogOpen(true);
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{t('oauth.applications')}</h1>
            <p className="text-muted-foreground">{t('oauth.applicationsDesc')}</p>
          </div>
        </div>
        <div className="grid gap-4">
          {[1, 2, 3].map(i => (
            <Card key={i} className="animate-pulse">
              <CardHeader>
                <div className="h-4 bg-muted rounded w-1/3"></div>
                <div className="h-3 bg-muted rounded w-1/2"></div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="h-3 bg-muted rounded"></div>
                  <div className="h-3 bg-muted rounded w-2/3"></div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{t('oauth.applications')}</h1>
          <p className="text-muted-foreground">{t('oauth.applicationsDesc')}</p>
        </div>
        <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              {t('oauth.createApplication')}
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>{t('oauth.createApplication')}</DialogTitle>
              <DialogDescription>
                Create a new OAuth application to integrate with MiniAuth.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="app-name">{t('oauth.applicationName')}</Label>
                <Input
                  id="app-name"
                  value={formData.name}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                  placeholder="My Application"
                />
              </div>
              
              <div>
                <Label>{t('oauth.redirectUris')}</Label>
                <div className="space-y-2">
                  {formData.redirect_uris.map((uri, index) => (
                    <div key={index} className="flex gap-2">
                      <Input
                        value={uri}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => updateRedirectUri(index, e.target.value)}
                        placeholder="https://yourapp.com/oauth/callback"
                        className="flex-1"
                      />
                      {formData.redirect_uris.length > 1 && (
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={() => removeRedirectUri(index)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                  ))}
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={addRedirectUri}
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Add URI
                  </Button>
                </div>
              </div>

              <div>
                <Label>{t('oauth.scopes')}</Label>
                <div className="grid grid-cols-2 gap-2 mt-2">
                  {availableScopes.map(scope => (
                    <div key={scope.value} className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id={`scope-${scope.value}`}
                        checked={formData.scopes.includes(scope.value)}
                        onChange={() => toggleScope(scope.value)}
                        className="rounded border-gray-300"
                      />
                      <Label htmlFor={`scope-${scope.value}`} className="text-sm">
                        {scope.value}
                      </Label>
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex items-center space-x-2">
                <Switch
                  id="trusted"
                  checked={formData.trusted}
                  onCheckedChange={(checked: boolean) => setFormData(prev => ({ ...prev, trusted: checked }))}
                />
                <Label htmlFor="trusted">{t('oauth.trusted')}</Label>
              </div>
              <p className="text-sm text-muted-foreground">
                {t('oauth.trustedDesc')}
              </p>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setCreateDialogOpen(false)}>
                {t('common.cancel')}
              </Button>
              <Button onClick={handleCreateApplication} disabled={!formData.name.trim()}>
                {t('common.create')}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {applications.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Settings className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-medium mb-2">{t('oauth.noApplications')}</h3>
            <p className="text-muted-foreground mb-4">{t('oauth.createFirstApp')}</p>
            <Button onClick={() => setCreateDialogOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              {t('oauth.createApplication')}
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6">
          {applications.map((app) => (
            <Card key={app.id}>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="flex items-center gap-2">
                      {app.name}
                      {app.trusted && (
                        <Badge variant="secondary">Trusted</Badge>
                      )}
                      <Badge variant={app.active ? "default" : "destructive"}>
                        {app.active ? t('oauth.active') : t('oauth.inactive')}
                      </Badge>
                    </CardTitle>
                    <CardDescription>
                      Created {new Date(app.created_at).toLocaleDateString()}
                    </CardDescription>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleToggleApplicationStatus(app.id)}
                      title={app.active ? t('oauth.deactivate') : t('oauth.activate')}
                    >
                      <Settings className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => openEditDialog(app)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        setSelectedApp(app);
                        setDeleteDialogOpen(true);
                      }}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm font-medium">{t('oauth.clientId')}</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <code className="flex-1 p-2 bg-muted rounded text-sm font-mono">
                        {app.client_id}
                      </code>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => copyToClipboard(app.client_id)}
                      >
                        <Copy className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                  <div>
                    <Label className="text-sm font-medium">{t('oauth.clientSecret')}</Label>
                    <div className="flex items-center gap-2 mt-1">
                      <code className="flex-1 p-2 bg-muted rounded text-sm font-mono">
                        {showSecrets[app.id] ? app.client_secret : '••••••••••••••••'}
                      </code>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => toggleSecretVisibility(app.id)}
                      >
                        {showSecrets[app.id] ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </Button>
                      {showSecrets[app.id] && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => copyToClipboard(app.client_secret)}
                        >
                          <Copy className="h-4 w-4" />
                        </Button>
                      )}
                    </div>
                  </div>
                </div>
                
                <div>
                  <Label className="text-sm font-medium">{t('oauth.redirectUris')}</Label>
                  <div className="mt-1 space-y-1">
                    {app.redirect_uris.map((uri, index) => (
                      <code key={index} className="block p-2 bg-muted rounded text-sm">
                        {uri}
                      </code>
                    ))}
                  </div>
                </div>

                <div>
                  <Label className="text-sm font-medium">{t('oauth.scopes')}</Label>
                  <div className="flex flex-wrap gap-1 mt-1">
                    {app.scopes.map((scope, index) => (
                      <Badge key={index} variant="outline">
                        {scope}
                      </Badge>
                    ))}
                  </div>
                </div>
                
                <div className="flex items-center justify-between pt-2 border-t">
                  <div className="flex items-center space-x-2">
                    <Switch
                      id={`trusted-${app.id}`}
                      checked={app.trusted}
                      onCheckedChange={() => handleToggleApplicationTrustedStatus(app.id)}
                    />
                    <Label htmlFor={`trusted-${app.id}`} className="text-sm font-medium">
                      {t('oauth.trusted')}
                    </Label>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {t('oauth.trustedDesc')}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t('oauth.confirmDeleteApp')}</DialogTitle>
            <DialogDescription>
              {t('oauth.deleteAppWarning')}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>
              {t('common.cancel')}
            </Button>
            <Button variant="destructive" onClick={handleDeleteApplication}>
              {t('common.delete')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{t('oauth.editApplication')}</DialogTitle>
            <DialogDescription>
              Edit the details of your OAuth application.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label htmlFor="app-name">{t('oauth.applicationName')}</Label>
              <Input
                id="app-name"
                value={formData.name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                placeholder="My Application"
              />
            </div>
            
            <div>
              <Label>{t('oauth.redirectUris')}</Label>
              <div className="space-y-2">
                {formData.redirect_uris.map((uri, index) => (
                  <div key={index} className="flex gap-2">
                    <Input
                      value={uri}
                      onChange={(e: React.ChangeEvent<HTMLInputElement>) => updateRedirectUri(index, e.target.value)}
                      placeholder="https://yourapp.com/oauth/callback"
                      className="flex-1"
                    />
                    {formData.redirect_uris.length > 1 && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => removeRedirectUri(index)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                ))}
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={addRedirectUri}
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Add URI
                </Button>
              </div>
            </div>

            <div>
              <Label>{t('oauth.scopes')}</Label>
              <div className="grid grid-cols-2 gap-2 mt-2">
                {availableScopes.map(scope => (
                  <div key={scope.value} className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      id={`scope-${scope.value}`}
                      checked={formData.scopes.includes(scope.value)}
                      onChange={() => toggleScope(scope.value)}
                      className="rounded border-gray-300"
                    />
                    <Label htmlFor={`scope-${scope.value}`} className="text-sm">
                      {scope.value}
                    </Label>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex items-center space-x-2">
              <Switch
                id="trusted"
                checked={formData.trusted}
                onCheckedChange={(checked: boolean) => setFormData(prev => ({ ...prev, trusted: checked }))}
              />
              <Label htmlFor="trusted">{t('oauth.trusted')}</Label>
            </div>
            <p className="text-sm text-muted-foreground">
              {t('oauth.trustedDesc')}
            </p>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditDialogOpen(false)}>
              {t('common.cancel')}
            </Button>
            <Button onClick={handleUpdateApplication} disabled={!formData.name.trim()}>
              {t('common.save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
