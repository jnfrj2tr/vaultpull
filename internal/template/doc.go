// Package template provides Go text/template-based rendering for .env files.
//
// It allows users to define .env templates that reference Vault secrets and
// local environment variables, which are resolved at sync time.
//
// Example template:
//
//	DB_HOST={{ .Secrets.DB_HOST }}
//	DB_USER={{ .Secrets.DB_USER }}
//	APP_ENV={{ .Env.APP_ENV }}
//
// Templates are rendered via Renderer.RenderFile or Renderer.RenderString.
// Missing secret keys cause an error by default (missingkey=error).
package template
