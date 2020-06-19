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
      console.log(error)
      return res.status(304).json({ error });
    }

    // res.json({ user: user[0] });
    const token = signToken(user);
    return res.status(200).json({ token });
  },
  signIn: async (req, res, next) => {
    const token = signToken(req.user.id);
    res.status(200).json({ token });
  },
  // facebookOAuth: (req, res, next) => {
  //   const token = signToken(req.user.f_id);
  //   res.cookie("access_token", token, {
  //     httpOnly: true,
  //   });
  //   res.status(200).json({ success: true });
  // },
  // linkFacebookOAuth: (req, res, next) => {
  //   res.json({
  //     success: true,
  //     message: "Successfully unlinked account from Facebook",
  //   });
  // },
  // unlinkFacebookOAuth: (req, res, next) => {},
  checkAuth: async (req, res, next) => {
    res.status(200).json({ success: true });
  },
};
