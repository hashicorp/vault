import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './ttl-form.md';

storiesOf('TtlForm', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `TtlForm`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Ttl Form</h5>
      <TtlForm/>
    `,
      context: {},
    }),
    { notes }
  );
