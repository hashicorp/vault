import Mirage from 'ember-cli-mirage';
import { faker } from 'ember-cli-mirage';

export default function() {
  // These comments are here to help you get started. Feel free to delete them.

  /*
    Config (with defaults).

    Note: these only affect routes defined *after* them!
  */

  // this.urlPrefix = '';    // make this `http://localhost:8080`, for example, if your API is on a different server
  this.namespace = '/v1'; // make this `/api`, for example, if your API is namespaced
  // this.timing = 400;      // delay for each request, automatically set to 0 during testing

  /*
    Shorthand cheatsheet:

    this.get('/posts');
    this.post('/posts');
    this.get('/posts/:id');
    this.put('/posts/:id'); // or this.patch
    this.del('/posts/:id');

    http://www.ember-cli-mirage.com/docs/v0.2.x/shorthands/

  */
  this.post('/sys/replication/primary/enable', schema => {
    var cluster = schema.clusters.first();
    cluster.update('mode', 'primary');
    return new Mirage.Response(204);
  });

  // primary_cluster_addr=(opt)

  this.post('/sys/replication/primary/demote', schema => {
    var cluster = schema.clusters.first();
    cluster.update('mode', 'secondary');
    return new Mirage.Response(204);
  });

  this.post('/sys/replication/primary/disable', schema => {
    var cluster = schema.clusters.first();
    cluster.update('mode', 'disabled');
    return new Mirage.Response(204);
  });
  this.post('/sys/replication/primary/secondary-token', (schema, request) => {
    //id=(req) ttl=(opt) (sudo)
    var params = JSON.parse(request.requestBody);
    var cluster = schema.clusters.first();

    if (!params.id) {
      return new Mirage.Response(400, {}, { errors: ['id must be specified'] });
    } else {
      var newSecondaries = (cluster.attrs.known_secondaries || []).slice();
      newSecondaries.push(params.id);
      cluster.update('known_secondaries', newSecondaries);
      return new Mirage.Response(200, {}, { token: faker.random.uuid() });
    }
  });

  this.post('/sys/replication/primary/revoke-secondary', (schema, request) => {
    var params = JSON.parse(request.requestBody);
    var cluster = schema.clusters.first();

    if (!params.id) {
      return new Mirage.Response(400, {}, { errors: ['id must be specified'] });
    } else {
      var newSecondaries = cluster.attrs.known_secondaries.without(params.id);
      cluster.update('known_secondaries', newSecondaries);
      return new Mirage.Response(204);
    }
  });

  this.post('/sys/replication/secondary/enable', (schema, request) => {
    //token=(req)
    var params = JSON.parse(request.requestBody);
    var cluster = schema.clusters.first();

    if (!params.token) {
      return new Mirage.Response(400, {}, { errors: ['token must be specified'] });
    } else {
      cluster.update('mode', 'secondary');
      return new Mirage.Response(204);
    }
  });

  this.post('/sys/replication/secondary/promote', schema => {
    var cluster = schema.clusters.first();
    cluster.update('mode', 'primary');
    return new Mirage.Response(204);
  });

  //primary_cluster_addr=(opt)
  this.post('/sys/replication/secondary/disable', schema => {
    var cluster = schema.clusters.first();
    cluster.update('mode', 'disabled');
    return new Mirage.Response(204);
  });

  this.post('/sys/replication/secondary/update-primary', (schema, request) => {
    //token=(req)
    var params = JSON.parse(request.requestBody);

    if (!params.token) {
      return new Mirage.Response(400, {}, { errors: ['token must be specified'] });
    } else {
      return new Mirage.Response(204);
    }
  });

  this.post('/sys/replication/recover', () => {
    return new Mirage.Response(204);
  });
  this.post('/sys/replication/reindex', () => {
    return new Mirage.Response(204);
  });
  //(sudo)
  this.get('/sys/replication/status', schema => {
    let model = schema.clusters.first();
    return new Mirage.Response(200, {}, model);
  }); //(unauthenticated)

  // enable and auth method
  this.post('/sys/auth/:path', (schema, request) => {
    const { path } = JSON.parse(request.requestBody);
    schema.authMethods.create({
      path,
    });
    return new Mirage.Response(204);
  });

  // TODO making this the default is probably not desired, but there's not an
  // easy way to do overrides currently - should this maybe just live in the
  // relevant test with pretender stubs?
  this.get('/sys/mounts', () => {
    return new Mirage.Response(403, {});
  });

  this.passthrough();
}
