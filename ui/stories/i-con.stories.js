/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './i-con.md';
import { GLYPHS_WITH_SVG_TAG } from '../app/components/i-con.js';

storiesOf('ICon/', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    'ICon',
    () => ({
      template: hbs`
        {{#each types as |type|}}
          <h5 class="title is-5">{{humanize type}}</h5>
          <ICon @glyph={{type}} />
          <br />
        {{/each}}
      `,
      context: {
        types: GLYPHS_WITH_SVG_TAG,
      },
    }),
    { notes }
  );
