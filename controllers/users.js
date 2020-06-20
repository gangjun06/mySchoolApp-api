const User = require("./../models/user");
const JWT = require("jsonwebtoken");
const config = require("./../config");

const signToken = (id) => {
  return JWT.sign(
    {
      iss: "gangjun",
      sub: id,
      iat: new Date().getTime(),
      exp: new Date().setDate(new Date().getDate() + 1),
    },
    config.JWT_SECRET
  );
};

module.exports = {
  signUp: async (req, res, next) => {
    const body = req.value.body;
    let user;

    try {
      user = await User.createUser(body);
    } catch (error) {
      return res.status(305).json({ error });
    }

    // res.json({ user: user[0] });
    const token = signToken(user);
    return res.status(200).json({ token });
  },
  signIn: async (req, res, next) => {
    const token = signToken(req.user.id);
    res.status(200).json({ token });
  },
  checkAuth: async (req, res, next) => {
    res.status(200).json({ success: true });
  },
};
