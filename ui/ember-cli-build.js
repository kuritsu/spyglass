'use strict';

const EmberApp = require('ember-cli/lib/broccoli/ember-app');

module.exports = function (defaults) {
  const app = new EmberApp(defaults, {
    'ember-bootstrap': {
      bootstrapVersion: 5,
      importBootstrapCSS: true,
    },
  });

  app.import('node_modules/bootstrap/dist/css/bootstrap.css');

  return app.toTree();
};
