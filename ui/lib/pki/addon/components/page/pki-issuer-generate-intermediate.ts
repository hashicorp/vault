import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import type PkiActionModel from 'vault/vault/models/pki/action';

interface Args {
  model: PkiActionModel;
}

export default class PagePkiIssuerGenerateIntermediateComponent extends Component<Args> {
  @tracked title = 'Generate intermediate CSR';
}
