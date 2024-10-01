/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { all } from 'rsvp';
import { verifyCertificates, parseCertificate } from 'vault/utils/parse-pki-cert';

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

  async getIssuerMetadata(store, type, query, response, id) {
    const keyInfo = response.data.key_info[id];
    try {
      const issuerRecord = await this.queryRecord(store, type, { id, backend: query.backend });
      const { data } = issuerRecord;
      const isRoot = await verifyCertificates(data.certificate, data.certificate);
      const parsedCertificate = parseCertificate(data.certificate);
      return {
        ...keyInfo,
        ...data,
        isRoot,
        parsedCertificate: { common_name: parsedCertificate.common_name },
      };
    } catch (e) {
      return { ...keyInfo, issuer_id: id };
    }
  }

  updateRecord(store, type, snapshot) {
    const { issuerId } = snapshot.record;
    const backend = this._getBackend(snapshot);
    const data = this.serialize(snapshot);
    const url = this.urlForQuery(backend, issuerId);
    return this.ajax(url, 'POST', { data });
  }

  query(store, type, query) {
    const { backend, isListView } = query;
    const url = this.urlForQuery(backend);

    return this.ajax(url, 'GET', this.optionsForQuery()).then(async (res) => {
      // To show issuer meta data tags, we have a flag called isListView and only want to
      // grab each issuer data only if there are less than 10 issuers to avoid making too many requests
      if (isListView && res.data.keys.length <= 10) {
        const keyInfoArray = await all(
          res.data.keys.map((id) => this.getIssuerMetadata(store, type, query, res, id))
        );
        const keyInfo = {};

        res.data.keys.forEach((issuerId) => {
          keyInfo[issuerId] = keyInfoArray.find((newKey) => newKey.issuer_id === issuerId);
        });

        res.data.key_info = keyInfo;

        return res;
      }

      return res;
    });
  }

  queryRecord(store, type, query) {
    const { backend, id } = query;

    return this.ajax(this.urlForQuery(backend, id), 'GET', this.optionsForQuery(id));
  }

  deleteAllIssuers(backend) {
    const deleteAllIssuersAndKeysUrl = `${this.buildURL()}/${encodePath(backend)}/root`;

    return this.ajax(deleteAllIssuersAndKeysUrl, 'DELETE');
  }
}
