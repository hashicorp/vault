/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';

storiesOf('AlertBanner', module)
  .add('warning', () => ({
    template: hbs`
      <AlertBanner @type="warning" @message={{message}} />
    `,
    context: {
      message: "Oops, don't do that again!",
    },
  }))
  .add('info', () => ({
    template: hbs`
      <AlertBanner @type="info" @message={{message}} />
    `,
    context: {
      message: "It's dangerous to go alone. Take this.",
    },
  }));
