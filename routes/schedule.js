const router = require("express-promise-router")();
const { validateBody, validateQuery, schemas } = require("../helpers/routeHelpers");
const UserController = require("../controllers/schedule");
const passport = require("passport");
const passportConf = require("../helpers/passport");

const passportJWT = passport.authenticate("jwt", { session: false });

router.route("/").get(UserController.getData);
router.route("/detail").get(validateQuery(schemas.scheduleRequest), UserController.detail);

router.route("/").post(validateBody(schemas.schedule), passportJWT, UserController.add)
router.route("/detail").post(validateBody(schemas.scheduleItem), passportJWT, UserController.addItem);

module.exports = router;
