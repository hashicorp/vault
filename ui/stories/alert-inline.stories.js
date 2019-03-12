/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './alert-inline.md';

const TYPES = ['warning', 'info', 'danger', 'success'];

storiesOf('AlertInline/', module)
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
        types: TYPES,
        message: 'Here is a message.',
      },
    }),
    { notes }
  );
