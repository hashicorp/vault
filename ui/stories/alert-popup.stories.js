/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './alert-popup.md';
import { MESSAGE_TYPES } from '../app/helpers/message-types.js';

storiesOf('Alerts/AlertPopup/', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    `AlertPopup`,
    () => ({
      template: hbs`
      {{#each types as |type|}}
        <h5 class="title is-5">{{humanize type}}</h5>
        <AlertPopup
          @type={{message-types type}}
          @message={{message}}
          @close={{close}}/>
      {{/each}}
    `,
      context: {
        close: () => {
          console.log('closing!');
        },
        types: Object.keys(MESSAGE_TYPES),
        message: 'Hello!',
      },
    }),
    { notes }
  );
