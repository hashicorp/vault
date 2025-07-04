/* tslint:disable */
/* eslint-disable */
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

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface KvV2PatchResponse
 */
export interface KvV2PatchResponse {
    /**
     * 
     * @type {Date}
     * @memberof KvV2PatchResponse
     */
    createdTime?: Date;
    /**
     * 
     * @type {object}
     * @memberof KvV2PatchResponse
     */
    customMetadata?: object;
    /**
     * 
     * @type {string}
     * @memberof KvV2PatchResponse
     */
    deletionTime?: string;
    /**
     * 
     * @type {boolean}
     * @memberof KvV2PatchResponse
     */
    destroyed?: boolean;
    /**
     * 
     * @type {number}
     * @memberof KvV2PatchResponse
     */
    version?: number;
}

/**
 * Check if a given object implements the KvV2PatchResponse interface.
 */
export function instanceOfKvV2PatchResponse(value: object): value is KvV2PatchResponse {
    return true;
}

export function KvV2PatchResponseFromJSON(json: any): KvV2PatchResponse {
    return KvV2PatchResponseFromJSONTyped(json, false);
}

export function KvV2PatchResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): KvV2PatchResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'createdTime': json['created_time'] == null ? undefined : (new Date(json['created_time'])),
        'customMetadata': json['custom_metadata'] == null ? undefined : json['custom_metadata'],
        'deletionTime': json['deletion_time'] == null ? undefined : json['deletion_time'],
        'destroyed': json['destroyed'] == null ? undefined : json['destroyed'],
        'version': json['version'] == null ? undefined : json['version'],
    };
}

export function KvV2PatchResponseToJSON(json: any): KvV2PatchResponse {
    return KvV2PatchResponseToJSONTyped(json, false);
}

export function KvV2PatchResponseToJSONTyped(value?: KvV2PatchResponse | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'created_time': value['createdTime'] == null ? undefined : ((value['createdTime']).toISOString()),
        'custom_metadata': value['customMetadata'],
        'deletion_time': value['deletionTime'],
        'destroyed': value['destroyed'],
        'version': value['version'],
    };
}

