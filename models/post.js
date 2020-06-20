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
    posts: async () => {
        const result = await DB("post")
            .join("user", "post.author", '=', 'user.id')
            .join("post_category", "post.category", '=', "post_category.id")
            .select("user.name", "post_category.text", "post.title", "post.maintext", "post.anon", "post.only_mygrade", "post.deleted")
        return result
    }

}