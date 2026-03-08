# Frontend Next.js (`/app`)

- Base path: `/app`
- Runtime: Next.js SSR/ISR (Node)
- Auth: cookie HttpOnly via rotas BFF (`/app/api/*`)
- DTOs: gerados de `../docs/swagger.json` (`npm run generate:types`)

## Comandos

```bash
npm install
npm run generate:types
npm run build
npm run start
```

## Endpoints internos (BFF)

- `POST /app/api/auth/login`
- `POST /app/api/auth/logout`
- `GET /app/api/session`
- `ALL /app/api/bff/*` (proxy para `/api/v1/*` com Bearer cookie)
