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

export type GinH = Record<string, any>;

export interface TypesAuthorizedResponseType {
  issue?: string;
  redirect_url?: string;
  status?: string;
}

export interface TypesCheckCredentialType {
  password?: string;
  secret?: string;
  username?: string;
}

export interface TypesCredentialType {
  password?: string;
  username?: string;
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
 * @version 1.0.0
 * @license Apache License Version 2.0 (https://github.com/mrzack99s/coco-captive-portal)
 * @baseUrl /api
 * @contact
 *
 * This is a COCO Captive Portal
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  v1 = {
    /**
     * @description Get admin signed
     *
     * @tags Authentication
     * @name AdmSigned
     * @summary Get admin signed
     * @request GET:/v1/adm-signed
     */
    admSigned: (params: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/v1/adm-signed`,
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
     * @request POST:/v1/authentication
     */
    authentication: (params: TypesCheckCredentialType, requestParams: RequestParams = {}) =>
      this.request<TypesAuthorizedResponseType, GinH>({
        path: `/v1/authentication`,
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
     * @request POST:/v1/check-is-administrator
     * @secure
     */
    checkIsAdministrator: (params: TypesCredentialType, requestParams: RequestParams = {}) =>
      this.request<GinH, GinH>({
        path: `/v1/check-is-administrator`,
        method: "POST",
        body: params,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description Get all session
     *
     * @tags Operator
     * @name GetAllSession
     * @summary Get all session
     * @request GET:/v1/get-all-session
     * @secure
     */
    getAllSession: (params: RequestParams = {}) =>
      this.request<TypesSessionType[], GinH>({
        path: `/v1/get-all-session`,
        method: "GET",
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Get html properties
     *
     * @tags HTML
     * @name HtmlProperties
     * @summary Get html properties
     * @request GET:/v1/html-properties
     */
    htmlProperties: (params: RequestParams = {}) =>
      this.request<TypesHTMLType, GinH>({
        path: `/v1/html-properties`,
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
     * @request GET:/v1/initialize
     */
    initialize: (params: RequestParams = {}) =>
      this.request<TypesInitializedType, GinH>({
        path: `/v1/initialize`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description To kick via ip address
     *
     * @tags Operator
     * @name KickViaIpAddress
     * @summary To kick via ip address
     * @request PUT:/v1/kick-ip-address
     * @secure
     */
    kickViaIpAddress: (params: TypesSessionType, requestParams: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/v1/kick-ip-address`,
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
     * @request PUT:/v1/kick-username
     * @secure
     */
    kickViaUsername: (params: TypesSessionType, requestParams: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/v1/kick-username`,
        method: "PUT",
        body: params,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...requestParams,
      }),

    /**
     * @description Revoke administrator
     *
     * @tags Operator
     * @name RevokeAdministrator
     * @summary Revoke administrator
     * @request GET:/v1/revoke-administrator
     * @secure
     */
    revokeAdministrator: (params: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/v1/revoke-administrator`,
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
     * @request GET:/v1/sign-out
     */
    signOut: (params: RequestParams = {}) =>
      this.request<string, GinH>({
        path: `/v1/sign-out`,
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
     * @request GET:/v1/signed
     */
    signed: (params: RequestParams = {}) =>
      this.request<TypesSessionType, GinH>({
        path: `/v1/signed`,
        method: "GET",
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
}
