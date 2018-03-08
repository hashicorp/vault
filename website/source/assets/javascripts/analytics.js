document.addEventListener('DOMContentLoaded', function() {
  track('.downloads .download .details li a', function(el) {
    var m = el.href.match(/vault_(.*?)_(.*?)_(.*?)\.zip/)
    return {
      event: 'Download',
      category: 'Button',
      label: 'Vault | v' + m[1] + ' | ' + m[2] + ' | ' + m[3],
      version: m[1],
      os: m[2],
      architecture: m[3],
      product: 'vault'
    }
  })
})
