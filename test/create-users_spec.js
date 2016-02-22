var frisby = require('frisby');
var URL = 'http://localhost:8000/users';

frisby.create('Create user Bob')
  .post(URL, {
    name: "Bob"
  }, {
    json: true
  })
  .expectStatus(201)
  .expectHeaderContains('content-type', 'application/json')
  .expectJSONTypes({
    id:   Number,
    name: String,
    type: String
  })
  .expectJSON({
    name: "Bob",
    type: "user"
  })
  .afterJSON(function () {
    frisby.create('Create user Jack')
      .post(URL, {
        name: "Jack"
      }, {
        json: true
      })
      .expectStatus(201)
      .expectHeaderContains('content-type', 'application/json')
      .expectJSONTypes({
        id:   Number,
        name: String,
        type: String
      })
      .expectJSON({
        name: "Jack",
        type: "user"
      })
      .afterJSON(function () {
        frisby.create('Create user Rose')
          .post(URL, {
            name: "Rose"
          }, {
            json: true
          })
          .expectStatus(201)
          .expectHeaderContains('content-type', 'application/json')
          .expectJSONTypes({
            id:   Number,
            name: String,
            type: String
          })
          .expectJSON({
            name: "Rose",
            type: "user"
          })
          .toss();
      })
      .toss();
  })
  .toss();
