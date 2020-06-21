const router = require("express-promise-router")();
const { validateBody, validateQuery, schemas } = require("../helpers/routeHelpers");
const PostController = require("../controllers/posts");
const passport = require("passport");
const passportConf = require("../helpers/passport");

const passportSingIn = passport.authenticate("local", { session: false });
const passportJWT = passport.authenticate("jwt", { session: false });

router.route('/').get(PostController.getPosts)
router.route("/").post(validateBody(schemas.post), passportJWT, PostController.addPost)

router.route("/category").get(PostController.getCategory)
router.route("/category").post(validateBody(schemas.category), passportJWT, PostController.addCategory);

router.route("/comment").get(validateQuery(schemas.commentQuery), PostController.getComments)
router.route("/comment").post(validateBody(schemas.comment), passportJWT, PostController.addComment)

module.exports = router;
