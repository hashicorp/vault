import { create } from '@storybook/theming';

// Fonts and colors are pulled from _colors.scss and _bulma_variables.scss.

const uiGray300 = '#BAC1CC';
const uiGray900 = '#1f2124';
const blue500 = '#1563ff';

export default create({
  base: 'light',

  colorPrimary: uiGray900,
  colorSecondary: blue500,

  // UI
  appBorderColor: uiGray300,

  // Typography
  fontBase: 'system-ui, -apple-system, BlinkMacSystemFont, sans-serif',
  fontCode: '"SFMono-Regular", Consolas, monospace',

  // Text colors
  textColor: uiGray900,

  // Toolbar default and active colors
  barTextColor: uiGray300,
  barSelectedColor: 'white',
  barBg: uiGray900,

  brandTitle: 'Vault UI Storybook',
  brandUrl: 'https://www.vaultproject.io/',
});
