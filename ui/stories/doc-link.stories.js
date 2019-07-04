/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './doc-link.md';


storiesOf('DocLink/', module)
  .addParameters({ options: { showPanel: true } })
  .add(`DocLink`, () => ({
    template: hbs`
      <h5 class="title is-5">Doc Link</h5>
      <DocLink @path="/docs/secrets/kv/kv-v2.html">Learn about KV v2</DocLink>
    `
  }),
  {notes}
);
