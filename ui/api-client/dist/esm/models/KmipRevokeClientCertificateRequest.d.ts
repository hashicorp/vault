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
 * @interface KmipRevokeClientCertificateRequest
 */
export interface KmipRevokeClientCertificateRequest {
    /**
     * PEM-encoded certificate from which to extract serial number.
     * @type {string}
     * @memberof KmipRevokeClientCertificateRequest
     */
    certificate?: string;
    /**
     * Serial number of the certificate.
     * @type {string}
     * @memberof KmipRevokeClientCertificateRequest
     */
    serialNumber?: string;
}
/**
 * Check if a given object implements the KmipRevokeClientCertificateRequest interface.
 */
export declare function instanceOfKmipRevokeClientCertificateRequest(value: object): value is KmipRevokeClientCertificateRequest;
export declare function KmipRevokeClientCertificateRequestFromJSON(json: any): KmipRevokeClientCertificateRequest;
export declare function KmipRevokeClientCertificateRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): KmipRevokeClientCertificateRequest;
export declare function KmipRevokeClientCertificateRequestToJSON(json: any): KmipRevokeClientCertificateRequest;
export declare function KmipRevokeClientCertificateRequestToJSONTyped(value?: KmipRevokeClientCertificateRequest | null, ignoreDiscriminator?: boolean): any;
