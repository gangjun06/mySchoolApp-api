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
  validateQuery: (schema) => {
    return (req, res, next) => {
      const result = Joi.validate(req.query, schema);
      if (result.error) {
        return res.status(400).json(result.error);
      }

      if (!req.error) {
        req.value = {};
      }
      req.value["body"] = result.value;
      next();
    }

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
    schedule: Joi.object().keys({
      text: Joi.string().required()
    }),
    scheduleRequest: Joi.object().keys({
      id: Joi.number().required(),
      day: Joi.number().required()
    }),
    scheduleItem: Joi.object().keys({
      schedule_id: Joi.number().required(),
      day: Joi.number().required(),
      order: Joi.number().required(),
      subject: Joi.string().required(),
      teacher: Joi.string().required(),
      description: Joi.string()
    }),
    category: Joi.object().keys({
      text: Joi.string().required(),
      admin_only: Joi.number().required()
    }),
    post: Joi.object().keys({
      title: Joi.string().required(),
      maintext: Joi.string(),
      category: Joi.number().required(),
      anon: Joi.number().required(),
      only_mygrade: Joi.number().required()
    })
  },
};
