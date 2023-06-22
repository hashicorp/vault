import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type LdapConfigModel from 'vault/models/ldap/config';
import type SecretEngineModel from 'vault/models/secret-engine';
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports
import type { Breadcrumb } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  configModel: LdapConfigModel;
  configError: AdapterError;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
}

interface Field {
  label: string;
  value: any; // eslint-disable-line @typescript-eslint/no-explicit-any
  formatTtl?: boolean;
}

export default class LdapConfigurationPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  get defaultFields(): Array<Field> {
    const model = this.args.configModel;
    const keys = [
      'binddn',
      'url',
      'schema',
      'password_policy',
      'userdn',
      'userattr',
      'connection_timeout',
      'request_timeout',
    ];
    return model.allFields.reduce<Array<Field>>((filtered, field) => {
      if (keys.includes(field.name)) {
        const label =
          {
            schema: 'Schema',
            password_policy: 'Password Policy',
          }[field.name] || field.options.label;
        filtered.splice(keys.indexOf(field.name), 0, {
          label,
          value: model[field.name as keyof typeof model],
          formatTtl: field.name.includes('timeout'),
        });
      }
      return filtered;
    }, []);
  }

  get connectionFields(): Array<Field> {
    const model = this.args.configModel;
    const keys = ['certificate', 'starttls', 'insecure_tls', 'client_tls_cert', 'client_tls_key'];
    return model.allFields.reduce<Array<Field>>((filtered, field) => {
      if (keys.includes(field.name)) {
        filtered.splice(keys.indexOf(field.name), 0, {
          label: field.options.label,
          value: model[field.name as keyof typeof model],
        });
      }
      return filtered;
    }, []);
  }

  @task
  @waitFor
  *rotateRoot() {
    try {
      yield this.args.configModel.rotateRoot();
      this.flashMessages.success('Root password successfully rotated.');
    } catch (error) {
      this.flashMessages.danger(`Error rotating root password \n ${errorMessage(error)}`);
    }
  }
}
