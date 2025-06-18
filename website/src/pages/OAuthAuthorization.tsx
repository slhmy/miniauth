import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Badge } from '../components/ui/badge';
import { Label } from '../components/ui/label';
import { Shield, User, ExternalLink } from 'lucide-react';

interface AuthorizationData {
  client_name: string;
  client_id: string;
  redirect_uri: string;
  scope: string;
  state: string;
  response_type: string;
  code_challenge?: string;
  code_challenge_method?: string;
  user: {
    id: number;
    username: string;
    email: string;
  };
}

export default function OAuthAuthorization() {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [authData, setAuthData] = useState<AuthorizationData | null>(null);
  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState(false);

  const scopeDescriptions = {
    read: t('oauth.authorization.scopeDescriptions.read'),
    write: t('oauth.authorization.scopeDescriptions.write'),
    profile: t('oauth.authorization.scopeDescriptions.profile'),
    organizations: t('oauth.authorization.scopeDescriptions.organizations'),
    admin: t('oauth.authorization.scopeDescriptions.admin'),
  };

  useEffect(() => {
    const fetchAuthorizationData = async () => {
      const params = new URLSearchParams();
      params.append('response_type', searchParams.get('response_type') || '');
      params.append('client_id', searchParams.get('client_id') || '');
      params.append('redirect_uri', searchParams.get('redirect_uri') || '');
      params.append('scope', searchParams.get('scope') || '');
      params.append('state', searchParams.get('state') || '');
      
      const codeChallenge = searchParams.get('code_challenge');
      const codeChallengeMethod = searchParams.get('code_challenge_method');
      
      if (codeChallenge) {
        params.append('code_challenge', codeChallenge);
      }
      if (codeChallengeMethod) {
        params.append('code_challenge_method', codeChallengeMethod);
      }

      try {
        const response = await fetch(`/api/oauth/authorize?${params.toString()}`, {
          credentials: 'include'
        });

        if (response.ok) {
          const data = await response.json();
          setAuthData(data);
        } else if (response.status === 302) {
          // Handle redirect (trusted app or login required)
          const location = response.headers.get('Location');
          if (location) {
            window.location.href = location;
          }
        } else {
          // Handle error
          navigate('/login');
        }
      } catch (error) {
        console.error('Failed to fetch authorization data:', error);
        navigate('/login');
      } finally {
        setLoading(false);
      }
    };

    fetchAuthorizationData();
  }, [searchParams, navigate]);

  const handleAuthorization = async (authorized: boolean) => {
    if (!authData) return;

    setProcessing(true);

    try {
      const response = await fetch('/api/oauth/authorize', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          authorized,
          response_type: authData.response_type,
          client_id: authData.client_id,
          redirect_uri: authData.redirect_uri,
          scope: authData.scope,
          state: authData.state,
          code_challenge: authData.code_challenge,
          code_challenge_method: authData.code_challenge_method,
        }),
      });

      if (response.ok) {
        const result = await response.json();
        if (result.redirect_url) {
          window.location.href = result.redirect_url;
        }
      }
    } catch (error) {
      console.error('Failed to process authorization:', error);
    } finally {
      setProcessing(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Card className="w-full max-w-md">
          <CardContent className="p-6">
            <div className="flex flex-col items-center space-y-4">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
              <p className="text-muted-foreground">{t('common.loading')}</p>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!authData) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Card className="w-full max-w-md">
          <CardContent className="p-6">
            <div className="flex flex-col items-center space-y-4 text-center">
              <Shield className="h-12 w-12 text-destructive" />
              <h2 className="text-xl font-semibold">{t('errors.unauthorized')}</h2>
              <p className="text-muted-foreground">
                Invalid authorization request or session expired.
              </p>
              <Button onClick={() => navigate('/login')}>
                {t('common.login')}
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  const requestedScopes = authData.scope.split(' ').filter(scope => scope.trim());

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center space-y-4">
          <div className="flex justify-center">
            <div className="relative">
              <Shield className="h-12 w-12 text-primary" />
              <ExternalLink className="h-4 w-4 text-muted-foreground absolute -top-1 -right-1" />
            </div>
          </div>
          <div>
            <CardTitle className="text-xl">{t('oauth.authorization.title')}</CardTitle>
            <CardDescription className="mt-2">
              {t('oauth.authorization.subtitle', { appName: authData.client_name })}
            </CardDescription>
          </div>
        </CardHeader>
        
        <CardContent className="space-y-6">
          {/* User Info */}
          <div className="flex items-center space-x-3 p-3 bg-muted rounded-lg">
            <User className="h-5 w-5 text-muted-foreground" />
            <div>
              <p className="font-medium">{authData.user.username}</p>
              <p className="text-sm text-muted-foreground">{authData.user.email}</p>
            </div>
          </div>

          {/* Requested Permissions */}
          <div>
            <Label className="text-sm font-medium">{t('oauth.authorization.requestedScopes')}</Label>
            <div className="mt-2 space-y-2">
              {requestedScopes.map((scope) => (
                <div key={scope} className="flex items-center justify-between p-3 border rounded-lg">
                  <div>
                    <Badge variant="outline" className="mb-1">
                      {scope}
                    </Badge>
                    <p className="text-sm text-muted-foreground">
                      {scopeDescriptions[scope as keyof typeof scopeDescriptions] || scope}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Application Info */}
          <div className="p-3 bg-muted/50 rounded-lg">
            <p className="text-sm text-muted-foreground">
              Redirecting to: <span className="font-mono text-xs">{authData.redirect_uri}</span>
            </p>
          </div>

          {/* Action Buttons */}
          <div className="grid grid-cols-2 gap-3">
            <Button
              variant="outline"
              onClick={() => handleAuthorization(false)}
              disabled={processing}
              className="w-full"
            >
              {t('oauth.authorization.denyAccess')}
            </Button>
            <Button
              onClick={() => handleAuthorization(true)}
              disabled={processing}
              className="w-full"
            >
              {processing ? t('common.loading') : t('oauth.authorization.allowAccess')}
            </Button>
          </div>

          <p className="text-xs text-muted-foreground text-center">
            By clicking "Allow Access", you authorize {authData.client_name} to access your account 
            with the permissions listed above.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
