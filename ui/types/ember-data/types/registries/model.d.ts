/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';
import KvSecretDataModel from 'vault/models/kv/data';
import KvSecretMetadataModel from 'vault/models/kv/metadata';
import PkiActionModel from 'vault/models/pki/action';
import PkiCertificateGenerateModel from 'vault/models/pki/certificate/generate';
import PkiConfigAcmeModel from 'vault/models/pki/config/acme';
import PkiConfigClusterModel from 'vault/models/pki/config/cluster';
import PkiConfigCrlModel from 'vault/models/pki/config/crl';
import PkiConfigUrlsModel from 'vault/models/pki/config/urls';
import ClientsActivityModel from 'vault/models/clients/activity';
import ClientsConfigModel from 'vault/models/clients/config';
import ClientsVersionHistoryModel from 'vault/models/clients/version-history';
import CaConfigModel from 'vault/models/ssh/ca-config';

declare module 'ember-data/types/registries/model' {
  export default interface ModelRegistry {
    'pki/action': PkiActionModel;
    'pki/certificate/generate': PkiCertificateGenerateModel;
    'pki/config/acme': PkiConfigAcmeModel;
    'pki/config/cluster': PkiConfigClusterModel;
    'pki/config/crl': PkiConfigCrlModel;
    'pki/config/urls': PkiConfigUrlModel;
    'kv/data': KvSecretDataModel;
    'kv/metadata': KvSecretMetadataModel;
    'clients/activity': ClientsActivityModel;
    'clients/config': ClientsConfigModel;
    'clients/version-history': ClientsVersionHistoryModel;
    'ssh/ca-config': CaConfigModel;
    // Catchall for any other models
    [key: string]: any;
  }
}
