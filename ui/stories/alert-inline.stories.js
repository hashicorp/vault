/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './alert-inline.md';
import { MESSAGE_TYPES } from '../app/helpers/message-types.js';

storiesOf('Alerts/AlertInline/', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    'AlertInline',
    () => ({
      template: hbs`
      {{#each types as |type|}}
        <h5 class="title is-5">{{humanize type}}</h5>
        <AlertInline @type={{type}} @message={{message}}/>
      {{/each}}
    `,
      context: {
        types: Object.keys(MESSAGE_TYPES),
        message: 'Here is a message.',
      },
    }),
    { notes }
  );
