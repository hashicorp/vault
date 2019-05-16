/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, text } from '@storybook/addon-knobs';
import notes from './empty-state.md';

storiesOf('EmptyState/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs({escapeHTML: false}))
  .add(`EmptyState`, () => ({
    template: hbs`
      <h5 class="title is-5">Empty State</h5>
      <EmptyState @title={{title}} @message={{message}} />
    `,
    context: {
      title: text('Title', 'You don\'t have an secrets yet'),
      message: text('Message', 'An explanation of why you don\'t have any secrets but also you maybe want to create one.')
    },
  }),
  {notes}
  )
  .add(`EmptyState with content block`, () => ({
    template: hbs`
      <h5 class="title is-5">Empty State</h5>
      <EmptyState @title={{title}} @message={{message}}>
        <DocLink @path="/docs/secrets/kv/kv-v2.html">Learn about KV v2</DocLink>
      </EmptyState>
    `,
    context: {
      title: text('Title', 'You don\'t have an secrets yet'),
      message: text('Message', 'An explanation of why you don\'t have any secrets but also you maybe want to create one.')
    },
  }),
  {notes}
  );
