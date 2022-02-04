!(function () {
  'use strict'

  const el = document.createElement('div');
  el.innerHTML =
    '<div class="ie-warning">' +
    '  <p class="ie-warning-description">' +
    '    Internet Explorer is no longer supported.' +
    '    <a href="https://support.hashicorp.com/hc/en-us/articles/4416485547795">' +
    '      Learn more.' +
    '    </a>' +
    '  </p>' +
    '</div>' +
    '<style>' +
    '  .ie-warning {' +
    '    background-color: #FCF0F2;' +
    '    border-bottom: 1px solid #FFD4D6;' +
    '    color: #BA2226;' +
    '    text-align: center;' +
    '    font-family: "Segoe UI", sans-serif;' +
    '    font-weight: bold;' +
    '  }' +
    '  .ie-warning-description {' +
    '    padding: 16px 0;' +
    '    margin: 0;' +
    '    color: #BA2226;' +
    '  }' +
    '</style>';

  document.body.insertBefore(el, document.body.childNodes[0]);
})();
