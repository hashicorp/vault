import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './alert-banner.md';
import { withKnobs, object } from '@storybook/addon-knobs';
import { MESSAGE_TYPES } from '../addon/helpers/message-types.js';

storiesOf('Alerts/AlertBanner/', module)
  .addParameters({ options: { showPanel: false } })
  .addDecorator(
    withKnobs({
      escapeHTML: false,
    })
  )
  .add(
    'AlertBanner',
    () => ({
      template: hbs`
      {{#each types as |type|}}
        <h5 class="title is-5">{{humanize type}}</h5>
        <AlertBanner @type={{type}} @message={{message}}/>
      {{/each}}
    `,
      context: {
        types: Object.keys(MESSAGE_TYPES),
        message: 'Here is a message.',
      },
    }),
    { notes }
  )
  .add(
    'AlertBanner with Progress Bar',
    () => ({
      template: hbs`
      <AlertBanner @type={{type}} @message={{message}} @progressBar={{progressBar}} />
    `,
      context: {
        type: 'info',
        message: 'Here is a message.',
        progressBar: object('Progress Bar', { value: 75, max: 100 }),
      },
    }),
    { notes }
  );
