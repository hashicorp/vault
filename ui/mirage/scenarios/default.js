export default function(server) {
  /*
    Seed your development database using your factories.
    This data will not be loaded in your tests.

    Make sure to define a factory for each model you want to create.
  */

  server.schema.clusters.create({
    name: 'vault',
    id: '1',
    mode: 'disabled',
  });
}
