/* Placeholder generated types for mobile app. Regenerate with npm run generate:types */
export type TokenResponse = { access_token: string; token_type: string };
export type Dashboard = {
  monthly_revenue_cents: number;
  mrr_cents: number;
  pending_payments: number;
  payments_received: number;
  active_contracts: number;
  expiring_contracts: number;
};
export type ProblemDetail = {
  error?: string;
  errorDescription?: string;
  statusCode?: number;
};
