/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { all } from 'rsvp';
import { verifyCertificates } from 'vault/utils/parse-pki-cert';

export default class PkiIssuerAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _getBackend(snapshot) {
    const { record, adapterOptions } = snapshot;
    return adapterOptions?.mount || record.backend;
  }

  optionsForQuery(id) {
    const data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  }

  urlForQuery(backend, id) {
    const baseUrl = `${this.buildURL()}/${encodePath(backend)}`;
    if (id) {
      return `${baseUrl}/issuer/${encodePath(id)}`;
    } else {
      return `${baseUrl}/issuers`;
    }
  }

  async getIssuerMetadata(store, type, query, response) {
    const { backend } = query;
    const newKeyInfo = await all(
      response.data.keys.map((id) => {
        const keyInfo = response.data.key_info[id];
        return this.queryRecord(store, type, { id, backend })
          .then(async (resp) => {
            const { certificate } = resp.data;
            const isRoot = await verifyCertificates(certificate, certificate);

            return { [id]: { ...keyInfo, isRoot } };
          })
          .catch(() => {
            return { [id]: { ...keyInfo } };
          });
      })
    );

    response.data.key_info = newKeyInfo;
    return response;
  }

  updateRecord(store, type, snapshot) {
    const { issuerId } = snapshot.record;
    const backend = this._getBackend(snapshot);
    const data = this.serialize(snapshot);
    const url = this.urlForQuery(backend, issuerId);
    return this.ajax(url, 'POST', { data });
  }

  query(store, type, query) {
    const { backend, shouldShowIssuerMetaData } = query;
    const url = this.urlForQuery(backend);

    return this.ajax(url, 'GET', this.optionsForQuery()).then(async (res) => {
      // To show issuer meta data tags, we have a flag called shouldShowIssuerMetaData and only want to
      // grab each issuer data only if there are less than 10 issuers to avoid making too many requests
      if (shouldShowIssuerMetaData && res.data.keys.length <= 10) {
        const issuerMetaData = await this.getIssuerMetadata(store, type, query, res);
        return issuerMetaData;
      }
      return res;
    });
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;

    return this.ajax(`${this.urlForQuery(backend, id)}`, 'GET', this.optionsForQuery(id));
  }

  deleteAllIssuers(backend) {
    const deleteAllIssuersAndKeysUrl = `${this.buildURL()}/${encodePath(backend)}/root`;

    return this.ajax(deleteAllIssuersAndKeysUrl, 'DELETE');
  }
}
