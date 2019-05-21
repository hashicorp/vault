/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './icon.md';
import icons from '../node_modules/@hashicorp/structure-icons/dist/index.js';
import { withKnobs, select } from '@storybook/addon-knobs';

storiesOf('Icon/', module)
  .addParameters({ options: { showPanel: true} })
  .addDecorator(withKnobs())
  .add(
    'Icon',
    () => ({
      template: hbs`
      <h5 class="title is-5">Icons from HashiCorp Structure</h5>
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
                <Icon @glyph={{type}} @size={{size}} />
              </td>
            </tr>
          {{/each}}
        </tbody>
      </table>
      `,
      context: {
        types: icons,
        size: select('Size', ['s', 'm', 'l', 'xl', 'xxl'], 'm'),
      },
    }),
    { notes }
  );
