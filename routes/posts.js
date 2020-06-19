const router = require("express-promise-router")();
const { validateBody, schemas } = require("../helpers/routeHelpers");
const UserController = require("../controllers/users");
const passport = require("passport");
const passportConf = require("../helpers/passport");

const passportSingIn = passport.authenticate("local", { session: false });
const passportJWT = passport.authenticate("jwt", { session: false });

router.route("/write").post(validateBody(), passportJWT, UserController.checkAuth);

module.exports = router;
