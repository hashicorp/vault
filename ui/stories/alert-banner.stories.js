/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './alert-banner.md';

const TYPES = ['warning', 'info', 'danger', 'success'];

storiesOf('AlertBanner/', module)
  .addParameters({ options: { showPanel: false } })
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
        types: TYPES,
        message: 'Here is a message.',
      },
    }),
    { notes }
  );
