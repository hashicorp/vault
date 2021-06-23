import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './vault-logo-spinner.md';

storiesOf('VaultLogoSpinner', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `VaultLogoSpinner`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Vault Logo Spinner</h5>
        <VaultLogoSpinner/>
    `,
      context: {},
    }),
    { notes }
  );
