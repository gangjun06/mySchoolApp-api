const DB = require("./../helpers/db");
const bcrypt = require("bcryptjs");

module.exports = {
  db: DB("users"),
  isValidPassword: (password, input) => {
    try {
      return bcrypt.compare(input, password);
    } catch (error) {
      throw new Error(error);
    }
  },
  createUser: async (body) => {
    const salt = await bcrypt.genSalt(10);
    const passwordHash = await bcrypt.hash(body.password, salt);
    const grade = body.teacher ? null : { grade: body.grade, class: body.class, class_number: body.class_number }
    const result = await DB("user")
      .insert({
        name: body.name,
        password: passwordHash,
        teacher: body.teacher,
        status: 0,
        ...grade
      })
      .catch((error) => {
        throw new Error(error);
      });
    return result[0];
  },
  findUserById: async (id) => {
    const result = await DB("user").where("id", id);
    return result[0];
  },
  findUserByName: async (name) => {
    const result = await DB("user").where("name", name);
    return result[0];
  },
};
