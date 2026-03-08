/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface AutomationRunResult {
  contracts_near_expiry?: number;
  monthly_report_generated?: number;
  projects_suspended?: number;
  recurring_charges_made?: number;
  reminders_sent?: number;
}

export interface BillingCreateChargeInput {
  amount_cents?: number;
  charge_type?: string;
  contract_id?: string;
  customer_id?: string;
  description?: string;
  due_date?: string;
  milestone_id?: string;
  payment_method?: string;
}

export interface BillingPaymentLinkResponse {
  charge_id?: string;
  link?: string;
}

export interface BillingPaymentQRResponse {
  charge_id?: string;
  qr_code?: string;
}

export interface ContractsCreateContractInput {
  amount_cents?: number;
  auto_renew?: boolean;
  billing_type?: string;
  contract_type?: string;
  customer_id?: string;
  deliverables?: string;
  duration_months?: number;
  end_date?: string;
  payment_terms?: string;
  penalties?: string;
  sla?: string;
  start_date?: string;
  title?: string;
}

export interface ContractsGenerateContractInput {
  editable_content?: string;
  template_name?: string;
}

export interface ContractsSendSignatureInput {
  signer_email?: string;
  signer_name?: string;
}

export interface CustomersCreateCustomerInput {
  address?: string;
  company?: string;
  cpf_cnpj?: string;
  email?: string;
  name?: string;
  phone?: string;
  preferred_payment_method?: string;
}

export interface CustomersUpdateCustomerInput {
  address?: string;
  company?: string;
  email?: string;
  name?: string;
  phone?: string;
  preferred_payment_method?: string;
}

export interface DomainUser {
  active?: boolean;
  balance?: number;
  createdAt?: string;
  email?: string;
  firstName?: string;
  id?: string;
  lastName?: string;
  phoneNumber?: string;
  updatedAt?: string;
}

export interface FinanceCashFlowPoint {
  date?: string;
  in_cents?: number;
  pending_cents?: number;
}

export interface FinanceDashboard {
  active_contracts?: number;
  expiring_contracts?: number;
  monthly_revenue_cents?: number;
  mrr_cents?: number;
  payments_received?: number;
  pending_payments?: number;
}

export interface FinanceDefaultMetrics {
  default_rate_percent?: number;
  growth_percent?: number;
  overdue_charges?: number;
}

export interface HttpCreateOrderRequest {
  /** @min 1 */
  amount_cents?: number;
  description?: string;
  user_id: string;
}

export interface HttpOidcAuthResponse {
  /** @example "https://accounts.google.com/o/oauth2/auth?..." */
  auth_url?: string;
  /** @example "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk" */
  code_verifier?: string;
  /** @example "eyJhbGci..." */
  state?: string;
}

export interface HttpOidcCallbackRequest {
  /** @example "4/0AX4XfWi..." */
  code?: string;
  /** @example "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk" */
  code_verifier?: string;
  /** @example "eyJhbGci..." */
  state?: string;
}

export interface HttpTokenResponse {
  /** @example "eyJhbGci..." */
  access_token?: string;
  /** @example "Bearer" */
  token_type?: string;
}

export interface HttpUpdateUserRequest {
  /** @example "João" */
  firstName?: string;
  /** @example "Silva" */
  lastName?: string;
  /** @example "11999999999" */
  phoneNumber?: string;
}

export interface HttpUserPageResponse {
  items?: DomainUser[];
  page?: number;
  pageSize?: number;
  total?: number;
  totalPages?: number;
}

export interface PaymentsConfirmPaymentInput {
  amount_cents?: number;
  charge_id?: string;
  method?: string;
  paid_at?: string;
  receipt_url?: string;
  tx_hash?: string;
}

export interface PaymentsEvidenceInput {
  file_url?: string;
  note?: string;
  tx_hash?: string;
}

export interface ProblemBaseErr {
  /** Details contém informações adicionais opcionais sobre o erro. */
  details?: any;
  /** Title é o título legível do tipo de problema. */
  error?: string;
  /** ErrorDescription é uma explicação específica desta ocorrência. */
  errorDescription?: string;
  /** Instance é o URI da requisição que originou o erro. */
  instance?: string;
  /** StatusCode é o código HTTP associado ao problema. */
  statusCode?: number;
  /** Timestamp indica quando o erro ocorreu (UTC, ISO 8601). */
  timestamp?: string;
  /** Type é um URI que identifica o tipo do problema (único por categoria de erro). */
  type?: string;
}

export interface ProjectsCreateMilestoneInput {
  amount_cents?: number;
  contract_id?: string;
  deliverables?: string;
  due_date?: string;
  title?: string;
}

export interface ProjectsCreateProjectInput {
  contract_id?: string;
  name?: string;
}

export interface ProjectsCreateWorklogInput {
  description?: string;
  hours?: number;
  milestone_id?: string;
  worked_at?: string;
}

export type GoogleCallbackCreateData = HttpTokenResponse;

export type GoogleLoginListData = HttpOidcAuthResponse;

export type LoginCreateData = HttpTokenResponse;

export type RegisterCreateData = DomainUser;

export interface UsersListParams {
  /** Página (padrão: 1) */
  page?: number;
  /** Tamanho da página (padrão: 10, máx: 100) */
  page_size?: number;
}

export type UsersListData = HttpUserPageResponse;

export interface UsersDetailParams {
  /** ID do usuário */
  id: string;
}

export type UsersDetailData = DomainUser;

export interface UsersUpdateParams {
  /** ID do usuário */
  id: string;
}

export type UsersUpdateData = DomainUser;

export interface UsersDeleteParams {
  /** ID do usuário */
  id: string;
}

export type UsersDeleteData = Record<string, string>;

export type PostAutomationData = AutomationRunResult;

export type RunsListData = Record<string, any>[];

export type ChargesListData = Record<string, any>[];

export type ChargesCreateData = Record<string, any>;

export interface PayLinkCreateParams {
  /** ID da cobrança */
  id: string;
}

export type PayLinkCreateData = BillingPaymentLinkResponse;

export interface PayQrCreateParams {
  /** ID da cobrança */
  id: string;
}

export type PayQrCreateData = BillingPaymentQRResponse;

export type ContractsListData = Record<string, any>[];

export type ContractsCreateData = Record<string, any>;

export interface ContractsDetailParams {
  /** ID do contrato */
  id: string;
}

export type ContractsDetailData = Record<string, any>;

export interface GenerateCreateParams {
  /** ID do contrato */
  id: string;
}

export type GenerateCreateData = Record<string, any>;

export interface SendSignatureCreateParams {
  /** ID do contrato */
  id: string;
}

export type SendSignatureCreateData = Record<string, any>;

export type CustomersListData = Record<string, any>[];

export type CustomersCreateData = Record<string, any>;

export interface CustomersDetailParams {
  /** ID do cliente */
  id: string;
}

export type CustomersDetailData = Record<string, any>;

export interface CustomersUpdateParams {
  /** ID do cliente */
  id: string;
}

export type CustomersUpdateData = Record<string, any>;

export type CashflowListData = FinanceCashFlowPoint[];

export type DashboardListData = FinanceDashboard;

export type DefaultsListData = FinanceDefaultMetrics;

export type OrdersCreateData = Record<string, any>;

export interface OrdersDetailParams {
  /** ID do pedido */
  id: string;
}

export type OrdersDetailData = Record<string, any>;

export type PaymentsListData = Record<string, any>[];

export type ConfirmCreateData = Record<string, any>;

export interface EvidenceCreateParams {
  /** ID do pagamento */
  id: string;
}

export type EvidenceCreateData = Record<string, any>;

export type ProjectsListData = Record<string, any>[];

export type ProjectsCreateData = Record<string, any>;

export interface MilestonesCreateParams {
  /** ID do projeto */
  id: string;
}

export type MilestonesCreateData = Record<string, any>;

export interface WorklogsCreateParams {
  /** ID do projeto */
  id: string;
}

export type WorklogsCreateData = Record<string, any>;

export type AsaasCreateData = Record<string, any>;

export type ClicksignCreateData = Record<string, any>;

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "//localhost:8000/api/v1";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) => {
      if (input instanceof FormData) {
        return input;
      }

      return Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData());
    },
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const responseToParse = responseFormat ? response.clone() : response;
      const data = !responseFormat
        ? r
        : await responseToParse[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title AggiPay API
 * @version 1.0.0
 * @baseUrl //localhost:8000/api/v1
 * @contact
 *
 * API modular de pagamentos com autenticação Bearer Token.
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  auth = {
    /**
     * @description Troca o authorization code pelo JWT da aplicação usando PKCE. O cliente envia o code recebido do Google, o state e o code_verifier gerados no passo anterior.
     *
     * @tags Auth
     * @name GoogleCallbackCreate
     * @summary Finalizar login com Google
     * @request POST:/auth/google/callback
     */
    googleCallbackCreate: (
      body: HttpOidcCallbackRequest,
      params: RequestParams = {},
    ) =>
      this.request<GoogleCallbackCreateData, ProblemBaseErr>({
        path: `/auth/google/callback`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Retorna a URL de autorização do Google e os parâmetros PKCE. O cliente deve armazenar state e code_verifier, depois redirecionar o usuário para auth_url.
     *
     * @tags Auth
     * @name GoogleLoginList
     * @summary Iniciar login com Google
     * @request GET:/auth/google/login
     */
    googleLoginList: (params: RequestParams = {}) =>
      this.request<GoogleLoginListData, ProblemBaseErr>({
        path: `/auth/google/login`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Autentica com email e senha, retorna JWT de acesso
     *
     * @tags Auth
     * @name LoginCreate
     * @summary Login local
     * @request POST:/auth/login
     */
    loginCreate: (
      data: {
        /** E-mail */
        email: string;
        /** Senha */
        password: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<LoginCreateData, ProblemBaseErr>({
        path: `/auth/login`,
        method: "POST",
        body: data,
        type: ContentType.UrlEncoded,
        format: "json",
        ...params,
      }),

    /**
     * @description Cria um novo usuário com email e senha. Usa criptografia Argon2ID com alto uso de processamento.
     *
     * @tags Auth
     * @name RegisterCreate
     * @summary Registrar usuário
     * @request POST:/auth/register
     */
    registerCreate: (
      data: {
        /** Primeiro nome */
        firstName: string;
        /** Sobrenome */
        lastName: string;
        /** E-mail */
        email: string;
        /** Telefone */
        phoneNumber?: string;
        /** Senha */
        password: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<RegisterCreateData, ProblemBaseErr>({
        path: `/auth/register`,
        method: "POST",
        body: data,
        type: ContentType.UrlEncoded,
        format: "json",
        ...params,
      }),

    /**
     * @description Retorna lista paginada de usuários
     *
     * @tags Users
     * @name UsersList
     * @summary Listar usuários
     * @request GET:/auth/users
     * @secure
     */
    usersList: (query: UsersListParams, params: RequestParams = {}) =>
      this.request<UsersListData, ProblemBaseErr>({
        path: `/auth/users`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Retorna um usuário pelo ID
     *
     * @tags Users
     * @name UsersDetail
     * @summary Buscar usuário
     * @request GET:/auth/users/{id}
     * @secure
     */
    usersDetail: (
      { id, ...query }: UsersDetailParams,
      params: RequestParams = {},
    ) =>
      this.request<UsersDetailData, ProblemBaseErr>({
        path: `/auth/users/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Atualiza os dados de um usuário pelo ID
     *
     * @tags Users
     * @name UsersUpdate
     * @summary Atualizar usuário
     * @request PUT:/auth/users/{id}
     * @secure
     */
    usersUpdate: (
      { id, ...query }: UsersUpdateParams,
      body: HttpUpdateUserRequest,
      params: RequestParams = {},
    ) =>
      this.request<UsersUpdateData, ProblemBaseErr>({
        path: `/auth/users/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Realiza soft-delete de um usuário pelo ID
     *
     * @tags Users
     * @name UsersDelete
     * @summary Desativar usuário
     * @request DELETE:/auth/users/{id}
     * @secure
     */
    usersDelete: (
      { id, ...query }: UsersDeleteParams,
      params: RequestParams = {},
    ) =>
      this.request<UsersDeleteData, ProblemBaseErr>({
        path: `/auth/users/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),
  };
  automation = {
    /**
     * No description
     *
     * @tags Automation
     * @name PostAutomation
     * @summary Executar automações financeiras
     * @request POST:/automation/run
     * @secure
     */
    postAutomation: (params: RequestParams = {}) =>
      this.request<PostAutomationData, ProblemBaseErr>({
        path: `/automation/run`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Automation
     * @name RunsList
     * @summary Histórico de automações
     * @request GET:/automation/runs
     * @secure
     */
    runsList: (params: RequestParams = {}) =>
      this.request<RunsListData, ProblemBaseErr>({
        path: `/automation/runs`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),
  };
  charges = {
    /**
     * No description
     *
     * @tags Billing
     * @name ChargesList
     * @summary Listar cobranças
     * @request GET:/charges
     * @secure
     */
    chargesList: (params: RequestParams = {}) =>
      this.request<ChargesListData, ProblemBaseErr>({
        path: `/charges`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Gera cobrança única, mensal ou por milestone
     *
     * @tags Billing
     * @name ChargesCreate
     * @summary Criar cobrança
     * @request POST:/charges
     * @secure
     */
    chargesCreate: (
      body: BillingCreateChargeInput,
      params: RequestParams = {},
    ) =>
      this.request<ChargesCreateData, ProblemBaseErr>({
        path: `/charges`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Billing
     * @name PayLinkCreate
     * @summary Obter link de pagamento
     * @request POST:/charges/{id}/pay-link
     * @secure
     */
    payLinkCreate: (
      { id, ...query }: PayLinkCreateParams,
      params: RequestParams = {},
    ) =>
      this.request<PayLinkCreateData, ProblemBaseErr>({
        path: `/charges/${id}/pay-link`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Billing
     * @name PayQrCreate
     * @summary Obter QR Code de pagamento
     * @request POST:/charges/{id}/pay-qr
     * @secure
     */
    payQrCreate: (
      { id, ...query }: PayQrCreateParams,
      params: RequestParams = {},
    ) =>
      this.request<PayQrCreateData, ProblemBaseErr>({
        path: `/charges/${id}/pay-qr`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),
  };
  contracts = {
    /**
     * No description
     *
     * @tags Contracts
     * @name ContractsList
     * @summary Listar contratos
     * @request GET:/contracts
     * @secure
     */
    contractsList: (params: RequestParams = {}) =>
      this.request<ContractsListData, ProblemBaseErr>({
        path: `/contracts`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Cadastra contrato de serviço de software
     *
     * @tags Contracts
     * @name ContractsCreate
     * @summary Criar contrato
     * @request POST:/contracts
     * @secure
     */
    contractsCreate: (
      body: ContractsCreateContractInput,
      params: RequestParams = {},
    ) =>
      this.request<ContractsCreateData, ProblemBaseErr>({
        path: `/contracts`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Contracts
     * @name ContractsDetail
     * @summary Buscar contrato
     * @request GET:/contracts/{id}
     * @secure
     */
    contractsDetail: (
      { id, ...query }: ContractsDetailParams,
      params: RequestParams = {},
    ) =>
      this.request<ContractsDetailData, ProblemBaseErr>({
        path: `/contracts/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Gera artefatos PDF/editável para assinatura
     *
     * @tags Contracts
     * @name GenerateCreate
     * @summary Gerar versão de contrato
     * @request POST:/contracts/{id}/generate
     * @secure
     */
    generateCreate: (
      { id, ...query }: GenerateCreateParams,
      body: ContractsGenerateContractInput,
      params: RequestParams = {},
    ) =>
      this.request<GenerateCreateData, ProblemBaseErr>({
        path: `/contracts/${id}/generate`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Cria documento no Clicksign e retorna URL de assinatura
     *
     * @tags Contracts
     * @name SendSignatureCreate
     * @summary Enviar para assinatura
     * @request POST:/contracts/{id}/send-signature
     * @secure
     */
    sendSignatureCreate: (
      { id, ...query }: SendSignatureCreateParams,
      body: ContractsSendSignatureInput,
      params: RequestParams = {},
    ) =>
      this.request<SendSignatureCreateData, ProblemBaseErr>({
        path: `/contracts/${id}/send-signature`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  customers = {
    /**
     * No description
     *
     * @tags Customers
     * @name CustomersList
     * @summary Listar clientes
     * @request GET:/customers
     * @secure
     */
    customersList: (params: RequestParams = {}) =>
      this.request<CustomersListData, ProblemBaseErr>({
        path: `/customers`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Cadastra um cliente para gestão financeira/contratos
     *
     * @tags Customers
     * @name CustomersCreate
     * @summary Criar cliente
     * @request POST:/customers
     * @secure
     */
    customersCreate: (
      body: CustomersCreateCustomerInput,
      params: RequestParams = {},
    ) =>
      this.request<CustomersCreateData, ProblemBaseErr>({
        path: `/customers`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Retorna cadastro e resumo financeiro do cliente
     *
     * @tags Customers
     * @name CustomersDetail
     * @summary Detalhar cliente
     * @request GET:/customers/{id}
     * @secure
     */
    customersDetail: (
      { id, ...query }: CustomersDetailParams,
      params: RequestParams = {},
    ) =>
      this.request<CustomersDetailData, ProblemBaseErr>({
        path: `/customers/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Customers
     * @name CustomersUpdate
     * @summary Atualizar cliente
     * @request PUT:/customers/{id}
     * @secure
     */
    customersUpdate: (
      { id, ...query }: CustomersUpdateParams,
      body: CustomersUpdateCustomerInput,
      params: RequestParams = {},
    ) =>
      this.request<CustomersUpdateData, ProblemBaseErr>({
        path: `/customers/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  finance = {
    /**
     * No description
     *
     * @tags Finance
     * @name CashflowList
     * @summary Fluxo de caixa
     * @request GET:/finance/cashflow
     * @secure
     */
    cashflowList: (params: RequestParams = {}) =>
      this.request<CashflowListData, ProblemBaseErr>({
        path: `/finance/cashflow`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Finance
     * @name DashboardList
     * @summary Dashboard financeiro
     * @request GET:/finance/dashboard
     * @secure
     */
    dashboardList: (params: RequestParams = {}) =>
      this.request<DashboardListData, ProblemBaseErr>({
        path: `/finance/dashboard`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Finance
     * @name DefaultsList
     * @summary Métricas de inadimplência e crescimento
     * @request GET:/finance/defaults
     * @secure
     */
    defaultsList: (params: RequestParams = {}) =>
      this.request<DefaultsListData, ProblemBaseErr>({
        path: `/finance/defaults`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),
  };
  orders = {
    /**
     * No description
     *
     * @tags Orders
     * @name OrdersCreate
     * @summary Cria um pedido de pagamento
     * @request POST:/orders
     */
    ordersCreate: (body: HttpCreateOrderRequest, params: RequestParams = {}) =>
      this.request<OrdersCreateData, any>({
        path: `/orders`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Orders
     * @name OrdersDetail
     * @summary Consulta um pedido pelo ID
     * @request GET:/orders/{id}
     */
    ordersDetail: (
      { id, ...query }: OrdersDetailParams,
      params: RequestParams = {},
    ) =>
      this.request<OrdersDetailData, any>({
        path: `/orders/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  payments = {
    /**
     * No description
     *
     * @tags Payments
     * @name PaymentsList
     * @summary Listar pagamentos
     * @request GET:/payments
     * @secure
     */
    paymentsList: (params: RequestParams = {}) =>
      this.request<PaymentsListData, ProblemBaseErr>({
        path: `/payments`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Confirma pagamento manual ou por conciliação externa
     *
     * @tags Payments
     * @name ConfirmCreate
     * @summary Confirmar pagamento
     * @request POST:/payments/confirm
     * @secure
     */
    confirmCreate: (
      body: PaymentsConfirmPaymentInput,
      params: RequestParams = {},
    ) =>
      this.request<ConfirmCreateData, ProblemBaseErr>({
        path: `/payments/confirm`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Anexa evidência/comprovante e hash de transação cripto
     *
     * @tags Payments
     * @name EvidenceCreate
     * @summary Adicionar comprovante
     * @request POST:/payments/{id}/evidence
     * @secure
     */
    evidenceCreate: (
      { id, ...query }: EvidenceCreateParams,
      body: PaymentsEvidenceInput,
      params: RequestParams = {},
    ) =>
      this.request<EvidenceCreateData, ProblemBaseErr>({
        path: `/payments/${id}/evidence`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  projects = {
    /**
     * No description
     *
     * @tags Projects
     * @name ProjectsList
     * @summary Listar projetos
     * @request GET:/projects
     * @secure
     */
    projectsList: (params: RequestParams = {}) =>
      this.request<ProjectsListData, ProblemBaseErr>({
        path: `/projects`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Projects
     * @name ProjectsCreate
     * @summary Criar projeto
     * @request POST:/projects
     * @secure
     */
    projectsCreate: (
      body: ProjectsCreateProjectInput,
      params: RequestParams = {},
    ) =>
      this.request<ProjectsCreateData, ProblemBaseErr>({
        path: `/projects`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Projects
     * @name MilestonesCreate
     * @summary Criar milestone
     * @request POST:/projects/{id}/milestones
     * @secure
     */
    milestonesCreate: (
      { id, ...query }: MilestonesCreateParams,
      body: ProjectsCreateMilestoneInput,
      params: RequestParams = {},
    ) =>
      this.request<MilestonesCreateData, ProblemBaseErr>({
        path: `/projects/${id}/milestones`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Projects
     * @name WorklogsCreate
     * @summary Criar worklog
     * @request POST:/projects/{id}/worklogs
     * @secure
     */
    worklogsCreate: (
      { id, ...query }: WorklogsCreateParams,
      body: ProjectsCreateWorklogInput,
      params: RequestParams = {},
    ) =>
      this.request<WorklogsCreateData, ProblemBaseErr>({
        path: `/projects/${id}/worklogs`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  webhooks = {
    /**
     * @description Recebe eventos de cobrança/pagamento do Asaas
     *
     * @tags Webhooks
     * @name AsaasCreate
     * @summary Webhook Asaas
     * @request POST:/webhooks/asaas
     */
    asaasCreate: (body: Record<string, any>, params: RequestParams = {}) =>
      this.request<AsaasCreateData, Record<string, any>>({
        path: `/webhooks/asaas`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Recebe eventos de assinatura digital do Clicksign
     *
     * @tags Webhooks
     * @name ClicksignCreate
     * @summary Webhook Clicksign
     * @request POST:/webhooks/clicksign
     */
    clicksignCreate: (body: Record<string, any>, params: RequestParams = {}) =>
      this.request<ClicksignCreateData, Record<string, any>>({
        path: `/webhooks/clicksign`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
