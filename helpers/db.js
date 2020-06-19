const { db } = require("./../config");

module.exports = require("knex")({
  client: "mysql",
  connection: {
    host: db.url,
    user: db.user,
    password: db.password,
    database: db.database,
  },
});
