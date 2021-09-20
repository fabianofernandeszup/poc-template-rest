const run = require("./formula/formula")

const language = process.env.RIT_LANGUAGE_LIST
const service = process.env.RIT_SERVICE_LIST
const currentPwd = process.env.CURRENT_PWD

run(language, service, currentPwd)
