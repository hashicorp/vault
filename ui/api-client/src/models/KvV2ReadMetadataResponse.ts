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
 * @interface KvV2ReadMetadataResponse
 */
export interface KvV2ReadMetadataResponse {
    /**
     * 
     * @type {boolean}
     * @memberof KvV2ReadMetadataResponse
     */
    casRequired?: boolean;
    /**
     * 
     * @type {Date}
     * @memberof KvV2ReadMetadataResponse
     */
    createdTime?: Date;
    /**
     * 
     * @type {number}
     * @memberof KvV2ReadMetadataResponse
     */
    currentVersion?: number;
    /**
     * User-provided key-value pairs that are used to describe arbitrary and version-agnostic information about a secret.
     * @type {object}
     * @memberof KvV2ReadMetadataResponse
     */
    customMetadata?: object;
    /**
     * The length of time before a version is deleted.
     * @type {string}
     * @memberof KvV2ReadMetadataResponse
     */
    deleteVersionAfter?: string;
    /**
     * The number of versions to keep
     * @type {number}
     * @memberof KvV2ReadMetadataResponse
     */
    maxVersions?: number;
    /**
     * 
     * @type {number}
     * @memberof KvV2ReadMetadataResponse
     */
    oldestVersion?: number;
    /**
     * 
     * @type {Date}
     * @memberof KvV2ReadMetadataResponse
     */
    updatedTime?: Date;
    /**
     * 
     * @type {object}
     * @memberof KvV2ReadMetadataResponse
     */
    versions?: object;
}

/**
 * Check if a given object implements the KvV2ReadMetadataResponse interface.
 */
export function instanceOfKvV2ReadMetadataResponse(value: object): value is KvV2ReadMetadataResponse {
    return true;
}

export function KvV2ReadMetadataResponseFromJSON(json: any): KvV2ReadMetadataResponse {
    return KvV2ReadMetadataResponseFromJSONTyped(json, false);
}

export function KvV2ReadMetadataResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): KvV2ReadMetadataResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'casRequired': json['cas_required'] == null ? undefined : json['cas_required'],
        'createdTime': json['created_time'] == null ? undefined : (new Date(json['created_time'])),
        'currentVersion': json['current_version'] == null ? undefined : json['current_version'],
        'customMetadata': json['custom_metadata'] == null ? undefined : json['custom_metadata'],
        'deleteVersionAfter': json['delete_version_after'] == null ? undefined : json['delete_version_after'],
        'maxVersions': json['max_versions'] == null ? undefined : json['max_versions'],
        'oldestVersion': json['oldest_version'] == null ? undefined : json['oldest_version'],
        'updatedTime': json['updated_time'] == null ? undefined : (new Date(json['updated_time'])),
        'versions': json['versions'] == null ? undefined : json['versions'],
    };
}

export function KvV2ReadMetadataResponseToJSON(json: any): KvV2ReadMetadataResponse {
    return KvV2ReadMetadataResponseToJSONTyped(json, false);
}

export function KvV2ReadMetadataResponseToJSONTyped(value?: KvV2ReadMetadataResponse | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'cas_required': value['casRequired'],
        'created_time': value['createdTime'] == null ? undefined : ((value['createdTime']).toISOString()),
        'current_version': value['currentVersion'],
        'custom_metadata': value['customMetadata'],
        'delete_version_after': value['deleteVersionAfter'],
        'max_versions': value['maxVersions'],
        'oldest_version': value['oldestVersion'],
        'updated_time': value['updatedTime'] == null ? undefined : ((value['updatedTime']).toISOString()),
        'versions': value['versions'],
    };
}

