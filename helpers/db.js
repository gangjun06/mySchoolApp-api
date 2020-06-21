const { db } = require("./../config");

const knex = require("knex")({
  client: "mysql",
  connection: {
    host: db.url,
    user: db.user,
    password: db.password,
    database: db.database,
  },
});
const setupPaginator = require('knex-paginator');
setupPaginator(knex);

module.exports = knex
