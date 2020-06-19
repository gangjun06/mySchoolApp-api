const passport = require("passport");
const JwtStrategy = require("passport-jwt").Strategy;
const ExtractJwt = require("passport-jwt").ExtractJwt;
const LocalStrategy = require("passport-local").Strategy;
const config = require("./../config");
const User = require("./../models/user");

passport.use(
  new JwtStrategy(
    {
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken("Authorization"),
      secretOrKey: config.JWT_SECRET,
    },
    async (payLoad, done) => {
      try {
        const user = await User.findUserById(payLoad.sub);

        if (user == undefined) {
          return done(null, false);
        }

        done(null, user);
      } catch (error) {
        done(error, false);
      }
    }
  )
);

passport.use(
  new LocalStrategy(
    {
      usernameField: "name",
    },
    async (name, password, done) => {
      try {
        const user = await User.findUserByName(name);

        if (user == undefined) {
          return done(null, false);
        }

        const isMatch = await User.isValidPassword(user.password, password);
        if (!isMatch) {
          return done(null, false);
        }

        done(null, user);
      } catch (error) {
        done(error, false);
      }
    }
  )
);
