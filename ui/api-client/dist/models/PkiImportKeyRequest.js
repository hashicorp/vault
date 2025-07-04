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
exports.instanceOfPkiImportKeyRequest = instanceOfPkiImportKeyRequest;
exports.PkiImportKeyRequestFromJSON = PkiImportKeyRequestFromJSON;
exports.PkiImportKeyRequestFromJSONTyped = PkiImportKeyRequestFromJSONTyped;
exports.PkiImportKeyRequestToJSON = PkiImportKeyRequestToJSON;
exports.PkiImportKeyRequestToJSONTyped = PkiImportKeyRequestToJSONTyped;
/**
 * Check if a given object implements the PkiImportKeyRequest interface.
 */
function instanceOfPkiImportKeyRequest(value) {
    return true;
}
function PkiImportKeyRequestFromJSON(json) {
    return PkiImportKeyRequestFromJSONTyped(json, false);
}
function PkiImportKeyRequestFromJSONTyped(json, ignoreDiscriminator) {
    if (json == null) {
        return json;
    }
    return {
        'keyName': json['key_name'] == null ? undefined : json['key_name'],
        'pemBundle': json['pem_bundle'] == null ? undefined : json['pem_bundle'],
    };
}
function PkiImportKeyRequestToJSON(json) {
    return PkiImportKeyRequestToJSONTyped(json, false);
}
function PkiImportKeyRequestToJSONTyped(value, ignoreDiscriminator = false) {
    if (value == null) {
        return value;
    }
    return {
        'key_name': value['keyName'],
        'pem_bundle': value['pemBundle'],
    };
}
