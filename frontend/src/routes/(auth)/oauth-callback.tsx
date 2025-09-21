import oauthService from '@/services/oauth.service'
import { useAuthStore } from '@/stores/auth'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/(auth)/oauth-callback')({
  component: OAuthCallbackComponent,
})

function OAuthCallbackComponent() {
  const navigate = useNavigate()
  const authStore = useAuthStore()
  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        const urlParams = new URLSearchParams(window.location.search)
        const token = urlParams.get('token')
        const success = urlParams.get('success')
        const error = urlParams.get('error')
        const linked = urlParams.get('linked')
        const provider = urlParams.get('provider')

        if (error) {
          let errorMessage = 'OAuth authentication failed'
          switch (error) {
            case 'authentication_required':
              errorMessage = 'You must be logged in to link an account'
              break
            case 'user_not_found':
              errorMessage = 'User not found'
              break
            case 'link_failed':
              errorMessage = 'Failed to link OAuth account'
              break
            case 'auth_failed':
              errorMessage = 'OAuth authentication failed'
              break
            case 'token_generation_failed':
              errorMessage = 'Failed to generate authentication tokens'
              break
            case 'temp_token_failed':
              errorMessage = 'Failed to create temporary token'
              break
          }

          toast.error(errorMessage)

          // Close popup after delay
          setTimeout(() => {
            if (window.opener) {
              window.close()
            } else {
              navigate({ to: '/login' })
            }
          }, 2000)
          return
        }

        if (linked === 'true' && provider) {
          // Account linking successful
          toast.success(`Successfully linked ${provider} account!`)

          // Close popup and refresh parent
          setTimeout(() => {
            if (window.opener) {
              window.opener.postMessage(
                {
                  type: 'OAUTH_LINK_SUCCESS',
                  provider,
                },
                window.location.origin,
              )
              window.close()
            } else {
              navigate({ to: '/me' })
            }
          }, 1000)
          return
        }

        if (success === 'true' && token) {
          console.log('Exchanging token:', token)
          const response = await oauthService.exchangeTempToken(token)
          console.log('OAuth exchange response:', response)

          if (!response || typeof response !== 'object') {
            throw new Error('Invalid response structure')
          }

          if (!response.user) {
            console.error('Response missing user field:', response)
            throw new Error('Response missing user data')
          }

          if (window.opener) {
            window.opener.postMessage(
              {
                type: 'OAUTH_SUCCESS',
                authData: response,
              },
              window.location.origin,
            )
            window.close()
          } else {
            authStore.login(response)
            toast.success('Successfully authenticated!')
            navigate({ to: '/' })
          }
        }
      } catch (error) {
        console.error('OAuth callback error:', error)
        toast.error('Failed to complete OAuth authentication')

        setTimeout(() => {
          if (window.opener) {
            window.close()
          } else {
            navigate({ to: '/login' })
          }
        }, 2000)
      }
    }

    handleOAuthCallback()
  }, [navigate])

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
        <p className="text-muted-foreground">
          Processing OAuth authentication...
        </p>
      </div>
    </div>
  )
}
