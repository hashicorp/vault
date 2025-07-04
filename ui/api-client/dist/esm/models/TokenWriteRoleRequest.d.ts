/**
 * HashiCorp Vault API
 * HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.
 *
 * The version of the OpenAPI document: 1.21.0
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */
/**
 *
 * @export
 * @interface TokenWriteRoleRequest
 */
export interface TokenWriteRoleRequest {
    /**
     * String or JSON list of allowed entity aliases. If set, specifies the entity aliases which are allowed to be used during token generation. This field supports globbing.
     * @type {Array<string>}
     * @memberof TokenWriteRoleRequest
     */
    allowedEntityAliases?: Array<string>;
    /**
     * If set, tokens can be created with any subset of the policies in this list, rather than the normal semantics of tokens being a subset of the calling token's policies. The parameter is a comma-delimited string of policy names.
     * @type {Array<string>}
     * @memberof TokenWriteRoleRequest
     */
    allowedPolicies?: Array<string>;
    /**
     * If set, tokens can be created with any subset of glob matched policies in this list, rather than the normal semantics of tokens being a subset of the calling token's policies. The parameter is a comma-delimited string of policy name globs.
     * @type {Array<string>}
     * @memberof TokenWriteRoleRequest
     */
    allowedPoliciesGlob?: Array<string>;
    /**
     * Use 'token_bound_cidrs' instead.
     * @type {Array<string>}
     * @memberof TokenWriteRoleRequest
     * @deprecated
     */
    boundCidrs?: Array<string>;
    /**
     * If set, successful token creation via this role will require that no policies in the given list are requested. The parameter is a comma-delimited string of policy names.
     * @type {Array<string>}
     * @memberof TokenWriteRoleRequest
     */
    disallowedPolicies?: Array<string>;
    /**
     * If set, successful token creation via this role will require that no requested policies glob match any of policies in this list. The parameter is a comma-delimited string of policy name globs.
     * @type {Array<string>}
     * @memberof TokenWriteRoleRequest
     */
    disallowedPoliciesGlob?: Array<string>;
    /**
     * Use 'token_explicit_max_ttl' instead.
     * @type {string}
     * @memberof TokenWriteRoleRequest
     * @deprecated
     */
    explicitMaxTtl?: string;
    /**
     * If true, tokens created via this role will be orphan tokens (have no parent)
     * @type {boolean}
     * @memberof TokenWriteRoleRequest
     */
    orphan?: boolean;
    /**
     * If set, tokens created via this role will contain the given suffix as a part of their path. This can be used to assist use of the 'revoke-prefix' endpoint later on. The given suffix must match the regular expression.\w[\w-.]+\w
     * @type {string}
     * @memberof TokenWriteRoleRequest
     */
    pathSuffix?: string;
    /**
     * Use 'token_period' instead.
     * @type {string}
     * @memberof TokenWriteRoleRequest
     * @deprecated
     */
    period?: string;
    /**
     * Tokens created via this role will be renewable or not according to this value. Defaults to "true".
     * @type {boolean}
     * @memberof TokenWriteRoleRequest
     */
    renewable?: boolean;
    /**
     * Comma separated string or JSON list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.
     * @type {Array<string>}
     * @memberof TokenWriteRoleRequest
     */
    tokenBoundCidrs?: Array<string>;
    /**
     * If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.
     * @type {string}
     * @memberof TokenWriteRoleRequest
     */
    tokenExplicitMaxTtl?: string;
    /**
     * If true, the 'default' policy will not automatically be added to generated tokens
     * @type {boolean}
     * @memberof TokenWriteRoleRequest
     */
    tokenNoDefaultPolicy?: boolean;
    /**
     * The maximum number of times a token may be used, a value of zero means unlimited
     * @type {number}
     * @memberof TokenWriteRoleRequest
     */
    tokenNumUses?: number;
    /**
     * If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").
     * @type {string}
     * @memberof TokenWriteRoleRequest
     */
    tokenPeriod?: string;
    /**
     * The type of token to generate, service or batch
     * @type {string}
     * @memberof TokenWriteRoleRequest
     */
    tokenType?: string;
}
/**
 * Check if a given object implements the TokenWriteRoleRequest interface.
 */
export declare function instanceOfTokenWriteRoleRequest(value: object): value is TokenWriteRoleRequest;
export declare function TokenWriteRoleRequestFromJSON(json: any): TokenWriteRoleRequest;
export declare function TokenWriteRoleRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): TokenWriteRoleRequest;
export declare function TokenWriteRoleRequestToJSON(json: any): TokenWriteRoleRequest;
export declare function TokenWriteRoleRequestToJSONTyped(value?: TokenWriteRoleRequest | null, ignoreDiscriminator?: boolean): any;
