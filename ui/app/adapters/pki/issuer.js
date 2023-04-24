/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { allSettled, all } from 'rsvp';
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

  async getIssuerMetaData(store, type, query, res) {
    const { backend, shouldShowIssuerMetaData } = query;

    const records = await allSettled(
      res.data.keys.map((id) => {
        return this.queryRecord(store, type, { id, backend });
      })
    );

    const issuerInfo = records.map((record) => {
      if (record.state === 'rejected') {
        return {};
      }
      if (record.state === 'fulfilled') {
        return record.value;
      }
    });

    // isRootRecords is a function that maps through each issuer meta data and uses the verifyCertificates function
    // to check if a certificate is a root. This is done in the adapter instead of the serializer because verifyCertificates
    // is an asynchronous function and we need to wait until every verifyCertificates function call is fulfilled.
    const isRootRecords = await all(
      issuerInfo.map((record) => {
        const { certificate } = record.data;
        return verifyCertificates(certificate, certificate);
      })
    );

    res.data.keys = issuerInfo.map((record, idx) => {
      return {
        ...record.data,
        shouldShowIssuerMetaData,
        isRoot: isRootRecords[idx],
      };
    });

    return res;
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
        const issuerMetaData = this.getIssuerMetaData(store, type, query, res);
        res = issuerMetaData;
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
