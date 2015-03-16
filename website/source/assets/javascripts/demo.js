window.Demo = Ember.Application.create({
  rootElement: '#demo-app',
});

Demo.deferReadiness();

if (document.getElementById('demo-app')) {
  Demo.advanceReadiness();
}
