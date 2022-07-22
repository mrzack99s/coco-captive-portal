/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface AuthenticationLDAPEndpointType {
  domain_names?: string[];
  hostname?: string;
  port?: number;
  single_domain?: boolean;
  tls_enable?: boolean;
}

export interface AuthenticationRadiusEndpointType {
  hostname?: string;
  port?: number;
  secret?: string;
}

export type GinH = Record<string, any>;

export interface TypesAuthorizedResponseType {
  issue?: string;
  redirect_url?: string;
  status?: string;
}

export interface TypesCaptivePortalConfigFundamentalType {
  domain_names?: string[];
  html?: TypesHTMLType;
  mode?: string;
  single_domain?: boolean;
}

export interface TypesCheckCredentialType {
  password?: string;
  secret?: string;
  username?: string;
}

export interface TypesConfigType {
  administrator?: TypesCredentialType;
  allow_endpoints?: TypesEndpointType[];
  bypass_networks?: string[];
  ddos_prevention?: boolean;
  domain_names?: { auth_domain_name?: string; operator_domain_name?: string };
  egress_interface?: string;
  external_portal_url?: string;
  fqdn_blocklist?: string[];
  html?: TypesHTMLType;
  ldap?: AuthenticationLDAPEndpointType;
  max_concurrent_session?: number;
  radius?: AuthenticationRadiusEndpointType;
  redirect_url?: string;
  secure_interface?: string;
  session_idle?: number;
}

export interface TypesCredentialType {
  password?: string;
  username?: string;
}

export interface TypesEndpointType {
  hostname?: string;
  port?: number;
}

export interface TypesExtendConfigType {
  administrator?: TypesCredentialType;
  allow_endpoints?: TypesEndpointType[];
  bypass_networks?: string[];
  ddos_prevention?: boolean;
  domain_names?: { auth_domain_name?: string; operator_domain_name?: string };
  egress_interface?: string;
  external_portal_url?: string;
  fqdn_blocklist?: string[];
  html?: TypesHTMLType;
  ldap?: AuthenticationLDAPEndpointType;
  max_concurrent_session?: number;
  radius?: AuthenticationRadiusEndpointType;
  redirect_url?: string;
  secure_interface?: string;
  session_idle?: number;
  status?: { egress_ip_address?: string; secure_ip_address?: string };
}

export interface TypesHTMLType {
  background_file_name?: string;
  default_language?: string;
  en_sub_title?: string;
  en_title_name?: string;
  logo_file_name?: string;
  th_sub_title?: string;
  th_title_name?: string;
}

export interface TypesInitializedType {
  ip_address?: string;
  secret?: string;
}

export interface TypesSessionType {
  ip_address?: string;
  issue?: string;
  last_seen?: string;
  session_uuid?: string;
}

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

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType | null) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "/api";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams);

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

  private encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  private addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  private addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) => (Array.isArray(query[key]) ? this.addArrayQueryParam(query, key) : this.addQueryParam(query, key)))
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
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
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  private mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
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

  private createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
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

    return this.customFetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      signal: cancelToken ? this.createAbortSignal(cancelToken) : void 0,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
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
 * @title COCO Captive Portal
 * @version 1
 * @license Apache License Version 2.0 (https://github.com/mrzack99s/coco-captive-portal)
 * @baseUrl /api
 * @contact
 *
 * This is a COCO Captive Portal
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  api = {
    /**
     * @description Get admin signed
     *
     * @tags Authentication
     * @name AdmSigned
     * @summary Get admin signed
     * @request GET:/api/adm-signed
     */
    admSigned: (params: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/adm-signed`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Check credential to get access
     *
     * @tags Authentication
     * @name Authentication
     * @summary Authentication
     * @request POST:/api/authentication
     */
    authentication: (params: TypesCheckCredentialType, requestParams: RequestParams = {}) =>
      this.request<TypesAuthorizedResponseType, GinH>({
        path: `/api/authentication`,
        method: "POST",
        body: params,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description Check is administrator
     *
     * @tags Operator
     * @name CheckIsAdministrator
     * @summary Check is administrator
     * @request POST:/api/check-is-administrator
     */
    checkIsAdministrator: (params: TypesCredentialType, requestParams: RequestParams = {}) =>
      this.request<GinH, GinH>({
        path: `/api/check-is-administrator`,
        method: "POST",
        body: params,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description Get config
     *
     * @tags Operator
     * @name GetConfig
     * @summary Get config
     * @request GET:/api/config
     * @secure
     */
    getConfig: (params: RequestParams = {}) =>
      this.request<TypesExtendConfigType, GinH>({
        path: `/api/config`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Set config
     *
     * @tags Operator
     * @name SetConfig
     * @summary Set config
     * @request PUT:/api/config
     * @secure
     */
    setConfig: (params: TypesConfigType, requestParams: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/config`,
        method: "PUT",
        body: params,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description Set config with restart system
     *
     * @tags Operator
     * @name SetConfigWithRestartSystem
     * @summary Set config with restart system
     * @request PUT:/api/config-with-restart-system
     * @secure
     */
    setConfigWithRestartSystem: (params: TypesConfigType, requestParams: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/config-with-restart-system`,
        method: "PUT",
        body: params,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description Count all session
     *
     * @tags Operator
     * @name CountAllSession
     * @summary Count all session
     * @request GET:/api/count-all-session
     * @secure
     */
    countAllSession: (params: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/count-all-session`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get all session
     *
     * @tags Operator
     * @name GetAllSession
     * @summary Get all session
     * @request GET:/api/get-all-session
     * @secure
     */
    getAllSession: (params: RequestParams = {}) =>
      this.request<TypesSessionType[], GinH>({
        path: `/api/get-all-session`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Captive portal config fundamental
     *
     * @tags Utils
     * @name GetCaptivePortalConfigFundamental
     * @summary Captive portal config fundamental
     * @request GET:/api/get-captive-portal-config-fundamental
     */
    getCaptivePortalConfigFundamental: (params: RequestParams = {}) =>
      this.request<TypesCaptivePortalConfigFundamentalType, GinH>({
        path: `/api/get-captive-portal-config-fundamental`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get initialize secret to get access
     *
     * @tags Authentication
     * @name Initialize
     * @summary Get initialize secret to get access
     * @request GET:/api/initialize
     */
    initialize: (params: RequestParams = {}) =>
      this.request<TypesInitializedType, GinH>({
        path: `/api/initialize`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Exist initialize secret
     *
     * @tags Authentication
     * @name IsExistInitializeSecret
     * @summary Exist initialize secret
     * @request POST:/api/is-exist-initialize-secret
     */
    isExistInitializeSecret: (params: TypesInitializedType, requestParams: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/is-exist-initialize-secret`,
        method: "POST",
        body: params,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description To kick via ip address
     *
     * @tags Operator
     * @name KickViaIpAddress
     * @summary To kick via ip address
     * @request PUT:/api/kick-ip-address
     * @secure
     */
    kickViaIpAddress: (params: TypesSessionType, requestParams: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/kick-ip-address`,
        method: "PUT",
        body: params,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description To kick via username
     *
     * @tags Operator
     * @name KickViaUsername
     * @summary To kick via username
     * @request PUT:/api/kick-username
     * @secure
     */
    kickViaUsername: (params: TypesSessionType, requestParams: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/kick-username`,
        method: "PUT",
        body: params,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description Network Interfaces Usage
     *
     * @tags Operator
     * @name NetInterfacesBytesUsage
     * @summary Network Interfaces Usage
     * @request GET:/api/net-intf-usage
     */
    netInterfacesBytesUsage: (params: RequestParams = {}) =>
      this.request<GinH, GinH>({
        path: `/api/net-intf-usage`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Revoke administrator
     *
     * @tags Operator
     * @name RevokeAdministrator
     * @summary Revoke administrator
     * @request GET:/api/revoke-administrator
     * @secure
     */
    revokeAdministrator: (params: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/revoke-administrator`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Sign out
     *
     * @tags Authentication
     * @name SignOut
     * @summary Sign out
     * @request GET:/api/sign-out
     */
    signOut: (params: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/api/sign-out`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get signed by ip address
     *
     * @tags Authentication
     * @name Signed
     * @summary Get signed by ip address
     * @request GET:/api/signed
     */
    signed: (params: RequestParams = {}) =>
      this.request<TypesSessionType, GinH>({
        path: `/api/signed`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
