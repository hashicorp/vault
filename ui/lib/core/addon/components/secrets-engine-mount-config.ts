import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

import type SecretEngineModel from 'vault/models/secret-engine';

interface Args {
  model: SecretEngineModel;
}
interface Field {
  label: string;
  value: string | boolean;
}

export default class SecretsEngineMountConfigComponent extends Component<Args> {
  @tracked showConfig = false;

  get fields(): Array<Field> {
    const { model } = this.args;
    return [
      { label: 'Secret Engine Type', value: model.engineType },
      { label: 'Path', value: model.path },
      { label: 'Accessor', value: model.accessor },
      { label: 'Local', value: model.local },
      { label: 'Seal Wrap', value: model.sealWrap },
      { label: 'Default Lease TTL', value: model.config.defaultLeaseTtl },
      { label: 'Max Lease TTL', value: model.config.maxLeaseTtl },
    ];
  }
}
