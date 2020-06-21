const Post = require('../models/post')

module.exports = {
    getCategory: async (req, res, next) => {
        let data;
        try {
            data = await Post.categories()
        } catch (e) {
            return res.status(500).json({ error: e })
        }
        res.status(200).json({ list: data })
    },
    addCategory: async (req, res, next) => {
        if (req.user.status !== 5) {
            res.status(401).json({ error: "You Don't Have a permission" })
        }
        try {
            await Post.addCategory(req.body)
        } catch (e) {
            return res.status(500).json({ error: e })
        }
        res.status(200).json({})
    },
    addPost: async (req, res, next) => {
        if (req.user.status !== 5 && req.user.status !== 1) {
            return res.status(401).json({ error: "You Don't Have a permission" })
        }
        req.body.author = req.user.id
        try {
            await Post.addPost(req.body)
        } catch (e) {
            return res.status(500).json({ error: e })
        }
        res.status(200).json({})
    },
    getPosts: async (req, res, next) => {
        let list
        try {
            list = await Post.posts(req.query.start)
        } catch (e) {
            return res.status(500).json({ error: e })
        }
        res.status(200).json({ list })
    },
    addComment: async (req, res, next) => {
        if (req.user.status !== 5 && req.user.status !== 1) {
            return res.status(401).json({ error: "You Don't Have a permission" })
        }
        req.body.author = req.user.id
        try {
            await Post.addComment(req.body)
        } catch (e) {
            return res.status(500).json({ error: e })
        }
        res.status(200).json({})
    },
    getComments: async (req, res, next) => {
        let list;
        try {
            list = await Post.comments(req.query.id)
        } catch (e) {
            return res.status(500).json({ error: e })
        }
        res.status(200).json({ list })
    }
}