import type {
  AutomationRunResult,
  BillingCreateChargeInput,
  BillingPaymentLinkResponse,
  BillingPaymentQRResponse,
  ContractsCreateContractInput,
  ContractsGenerateContractInput,
  ContractsSendSignatureInput,
  CustomersCreateCustomerInput,
  CustomersUpdateCustomerInput,
  FinanceCashFlowPoint,
  FinanceDashboard,
  FinanceDefaultMetrics,
  HttpTokenResponse,
  PaymentsConfirmPaymentInput,
  PaymentsEvidenceInput,
  ProblemBaseErr,
  ProjectsCreateMilestoneInput,
  ProjectsCreateProjectInput,
  ProjectsCreateWorklogInput,
} from '@/generated/api';

export type ProblemDetail = ProblemBaseErr;
export type TokenResponse = HttpTokenResponse;

export type Dashboard = FinanceDashboard;
export type CashFlowPoint = FinanceCashFlowPoint;
export type DefaultMetrics = FinanceDefaultMetrics;

export type Customer = Record<string, unknown>;
export type CustomerFinancialSummary = Record<string, unknown>;
export type CustomerDetail = Record<string, any>;
export type CreateCustomerInput = CustomersCreateCustomerInput;
export type UpdateCustomerInput = CustomersUpdateCustomerInput;

export type Contract = Record<string, unknown>;
export type CreateContractInput = ContractsCreateContractInput;
export type GenerateContractInput = ContractsGenerateContractInput;
export type SendSignatureInput = ContractsSendSignatureInput;

export type Charge = Record<string, unknown>;
export type CreateChargeInput = BillingCreateChargeInput;
export type PaymentLinkResponse = BillingPaymentLinkResponse;
export type PaymentQRResponse = BillingPaymentQRResponse;

export type Payment = Record<string, unknown>;
export type ConfirmPaymentInput = PaymentsConfirmPaymentInput;
export type EvidenceInput = PaymentsEvidenceInput;

export type Project = Record<string, unknown>;
export type CreateProjectInput = ProjectsCreateProjectInput;
export type CreateMilestoneInput = ProjectsCreateMilestoneInput;
export type CreateWorklogInput = ProjectsCreateWorklogInput;

export type AutomationRun = Record<string, unknown>;
export type RunResult = AutomationRunResult;
