/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './alert-banner.md';

storiesOf('AlertBanner/', module).add(
  'AlertBanner',
  () => ({
    template: hbs`
      <h1>Warning</h1>
      <AlertBanner @type="warning" @message={{message}}/>
      <h1>Info</h1>
      <AlertBanner @type="info" @message={{message}}/>
      <h1>Danger</h1>
      <AlertBanner @type="danger" @message={{message}}/>
      <h1>Success</h1>
      <AlertBanner @type="success" @message={{message}}/>
    `,
    context: {
      message: 'Here is a message.',
    },
  }),
  { notes }
);
