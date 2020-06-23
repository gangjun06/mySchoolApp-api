const DB = require("./../helpers/db");

module.exports = {
    categories: async () => {
        const result = await DB("post_category").select("*")
        return result
    },
    addCategory: async (body) => {
        const result = await DB("post_category").insert(body)
        return result[0]
    },
    addPost: async (body) => {
        const result = await DB("post").insert({ ...body, deleted: 0, timestamp: DB.fn.now() })
        return result[0]
    },
    posts: async (start) => {
        const result = await DB("post")
            .join("user", "post.author", '=', 'user.id')
            .join("post_category", "post.category", '=', "post_category.id")
            .select("post.id", "user.name", "post_category.text", "post.title", "post.maintext", "post.anon", "post.deleted", "post.timestamp")
            .orderBy("post.id", "desc")
            .limit(30)
            // .paginate(20, 1, true)
        return result
    },
    addComment: async (body) => {
        const result = await DB("post_comment").insert({ ...body, timestamp: DB.fn.now() }).catch(e => {
            throw new Error(e)
        })
        return result[0]
    },
    comments: async (id) => {
        const result = await DB("post")
            .where('post.id', '=', id)
            .join("post_comment", 'post.id', '=', 'post_comment.post_id')
            .join("user", "user.id", '=', 'post_comment.author')
            .select("post_comment.id", 'post_comment.parent', 'post_comment.timestamp', 'post_comment.maintext', 'user.name')
            .orderBy("timestamp", "desc")
        return result
    }
}