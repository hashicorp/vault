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
      <table class="table">
        <thead>
          <tr>
            <th>Glyph title</th>
            <th>Glyph</th>
          </tr>
        </thead>
        <tbody>
          {{#each types as |type|}}
            <tr>
              <td>
                <h5>{{humanize type}}</h5>
              </td>
              <td>
                <ICon @glyph={{type}} />
              </td>
            </tr>
          {{/each}}
        </tbody>
      </table>
      `,
      context: {
        types: GLYPHS_WITH_SVG_TAG,
      },
    }),
    { notes }
  );
