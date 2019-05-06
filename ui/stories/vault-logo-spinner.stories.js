/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './vault-logo-spinner.md';


storiesOf('Loading/VaultLogoSpinner/', module)
  .addParameters({ options: { showPanel: true } })
  .add(`VaultLogoSpinner`, () => ({
    template: hbs`
        <h5 class="title is-5">Vault Logo Spinner</h5>
        <VaultLogoSpinner/>
    `,
    context: {},
  }),
  {notes}
);
