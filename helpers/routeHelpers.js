const Joi = require("joi");

module.exports = {
  validateBody: (schema) => {
    return (req, res, next) => {
      const result = Joi.validate(req.body, schema);
      if (result.error) {
        return res.status(400).json(result.error);
      }

      if (!req.error) {
        req.value = {};
      }
      req.value["body"] = result.value;
      next();
    };
  },
  schemas: {
    authSchema: Joi.object().keys({
      name: Joi.string().required(),
      password: Joi.string().required(),
    }),
    signUpSchema: Joi.object().keys({
      name: Joi.string().required(),
      password: Joi.string().required(),
      grade: Joi.number(),
      class: Joi.number(),
      class_number: Joi.number(),
      teacher: Joi.number().required()
    }),
  },
};
