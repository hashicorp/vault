import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './icon.md';
import icons from '../../../node_modules/@hashicorp/structure-icons/dist/index.js';
import { withKnobs, select } from '@storybook/addon-knobs';
import { structureIconMap, localIconMap } from '../icon-mappings';

storiesOf('Icon', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    'Icon',
    () => ({
      template: hbs`
      <h5 class="title is-5">HashiCorp Flight Icons</h5>
      <a href="https://flight-hashicorp.vercel.app/">https://flight-hashicorp.vercel.app/</a>

      <h5 class="title is-5 has-top-margin-l">
        HashiCorp Structure Icons with Flight Mappings
      </h5>
      <table class="table">
        <thead>
          <tr>
            <th>Structure Icon Name</th>
            <th>Structure Glyph</th>
            <th>Flight Icon Name</th>
            <th>Flight Glyph></th>
          </tr>
        </thead>
        <tbody>
          {{#each types as |type|}}
            <tr>
              <td>
                {{type}}
              </td>
              <td>
                <Icon @name={{type}} @size={{size}} />
              </td>
              {{#let (get structureIconMap type) as |flightIcon|}}
                <td>
                  {{#if flightIcon}}
                    {{flightIcon}}
                  {{else}}
                    &mdash;
                  {{/if}}
                </td>
                <td>
                  {{#if flightIcon}}
                    <Icon @name={{flightIcon}} @size={{size}} />
                  {{else}}
                    &mdash;
                  {{/if}}
                </td>
              {{/let}}
            </tr>
          {{/each}}
        </tbody>
      </table>

      <h5 class="title is-5 has-top-margin-l">
        Local Icons with Flight Mappings
      </h5>
      <table class="table">
        <thead>
          <tr>
            <th>Local Icon Name</th>
            <th>Glyph</th>
            <th>Flight Icon Name</th>
            <th>Flight Glyph</th>
          </tr>
        </thead>
        <tbody>
          {{#each-in localIconMap as |localIcon flightIcon|}}
            <tr>
              <td>
                {{localIcon}}
              </td>
              <td>
                <Icon @name={{localIcon}} @size={{size}} />
              </td>
              <td>
                {{#if flightIcon}}
                  {{flightIcon}}
                {{else}}
                  &mdash;
                {{/if}}
              </td>
              <td>
                {{#if flightIcon}}
                  <Icon @name={{flightIcon}} @size={{size}} />
                {{else}}
                  &mdash;
                {{/if}}
              </td>
            </tr>
          {{/each-in}}
        </tbody>
      </table>
      `,
      context: {
        types: icons,
        structureIconMap,
        localIconMap,
        size: select('Size', ['16', '24'], '16'),
      },
    }),
    { notes }
  );
