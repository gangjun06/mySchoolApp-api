const DB = require("./../helpers/db");

module.exports = {
    schedules: async () => {
        const result = await DB("schedule").select("*").catch(error => {
            throw new Error(error)
        })
        return result
    },
    schedule: async ({ id, day }) => {
        const result = await DB("schedule")
            .join("schedule_item", 'schedule.id', '=', 'schedule_item.schedule_id')
            .where('schedule_item.day', '=', day)
            .select("schedule_item.subject", "schedule_item.teacher", "schedule_item.description", "schedule_item.order")
            .orderBy('schedule_item.order')
        return result
    },
    addSchedule: async (body) => {
        const result = await DB("schedule")
            .insert({
                text: body.text,
            })
            .catch((error) => {
                throw new Error(error);
            });
        return result[0];
    },
    addScheduleItem: async (body) => {
        const result = await DB("schedule_item").insert({
            ...body
        })
        return result[0]
    }
};
