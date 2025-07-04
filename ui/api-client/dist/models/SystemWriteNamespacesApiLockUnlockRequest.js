"use strict";
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
Object.defineProperty(exports, "__esModule", { value: true });
exports.instanceOfSystemWriteNamespacesApiLockUnlockRequest = instanceOfSystemWriteNamespacesApiLockUnlockRequest;
exports.SystemWriteNamespacesApiLockUnlockRequestFromJSON = SystemWriteNamespacesApiLockUnlockRequestFromJSON;
exports.SystemWriteNamespacesApiLockUnlockRequestFromJSONTyped = SystemWriteNamespacesApiLockUnlockRequestFromJSONTyped;
exports.SystemWriteNamespacesApiLockUnlockRequestToJSON = SystemWriteNamespacesApiLockUnlockRequestToJSON;
exports.SystemWriteNamespacesApiLockUnlockRequestToJSONTyped = SystemWriteNamespacesApiLockUnlockRequestToJSONTyped;
/**
 * Check if a given object implements the SystemWriteNamespacesApiLockUnlockRequest interface.
 */
function instanceOfSystemWriteNamespacesApiLockUnlockRequest(value) {
    return true;
}
function SystemWriteNamespacesApiLockUnlockRequestFromJSON(json) {
    return SystemWriteNamespacesApiLockUnlockRequestFromJSONTyped(json, false);
}
function SystemWriteNamespacesApiLockUnlockRequestFromJSONTyped(json, ignoreDiscriminator) {
    if (json == null) {
        return json;
    }
    return {
        'unlockKey': json['unlock_key'] == null ? undefined : json['unlock_key'],
    };
}
function SystemWriteNamespacesApiLockUnlockRequestToJSON(json) {
    return SystemWriteNamespacesApiLockUnlockRequestToJSONTyped(json, false);
}
function SystemWriteNamespacesApiLockUnlockRequestToJSONTyped(value, ignoreDiscriminator = false) {
    if (value == null) {
        return value;
    }
    return {
        'unlock_key': value['unlockKey'],
    };
}
