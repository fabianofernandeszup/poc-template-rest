var fs = require('fs');
const path = require('path');

function Run(language, service, currentPwd) {
    const fileName = language + '_' + service + '.yaml'
    const originPath = path.join('templates', fileName);
    const targetPath = path.join(currentPwd, fileName);

    fs.copyFile(originPath, targetPath, (err) => {
        if (err) throw err;
        console.log('pipeline file was copied to destination');
    });
}

const formula = Run
module.exports = formula
