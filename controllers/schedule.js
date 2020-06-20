const Schedule = require('../models/schedule')

module.exports = {
    getData: async (req, res, next) => {
        let data;
        try {
            data = await Schedule.schedules()
        } catch (e) {
            return res.status(305).json({ success: false, error: e })
        }
        res.status(200).json({ data })
    },
    detail: async (req, res, next) => {
        let data;
        try {
            data = await Schedule.schedule(req.query)
        } catch (e) {
            return res.status(305).json({ success: false, error: e })
        }
        res.status(200).json({ list: data })

    },
    add: async (req, res, next) => {
        if (req.user.status !== 5) {
            res.status(401).json({ error: "You Don't Have a permission" })
        }

        try {
            await Schedule.addSchedule(req.body)
        } catch (e) {
            return res.status(305).json({ error: e })
        }

        res.status(200).json({})

    },
    delete: async (req, res, next) => {

    },
    addItem: async (req, res, next) => {
        if (req.user.status !== 5) {
            res.status(401).json({ error: "You Don't Have a permission" })
        }
        try {
            await Schedule.addScheduleItem(req.body)
        } catch (e) {
            return res.status(305).json({ error: e })
        }

        res.status(200).json({})
    },

}