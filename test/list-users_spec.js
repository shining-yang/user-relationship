var frisby = require('frisby');
var URL = 'http://localhost:8000/users';

frisby.create('List all users')
  .get(URL)
  .expectStatus(200)
  .expectHeaderContains('content-type', 'application/json')
  .inspectJSON()
  .toss();
