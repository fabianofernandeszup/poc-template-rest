const clc = require('cli-color')
const axios = require('axios')
const YAML = require('yaml')
const fs = require('fs')

const formula = async () => {
    console.log("\n")

    const formulaConfigYAML = 'formula.yml'
    const formulaConfigJSON = 'config.json'
    let formula = {}
    if(fs.existsSync(formulaConfigYAML)) {
        const file = fs.readFileSync(formulaConfigYAML, 'utf8')
        formula = YAML.parse(file)
    } else if (fs.existsSync(formulaConfigJSON)) {
        const file = fs.readFileSync(formulaConfigJSON, 'utf8')
        formula = JSON.parse(file)
    } else {
        console.log('Files formula.yml or config.json not found.')
        process.exit(1)
    }

    const inputsArray = formula.inputs.map(input => {
        let prop = input.name.toLowerCase();
        let propEnv = input.name.toUpperCase();
        let obj = new Object();
        obj[prop] = process.env[propEnv];
        return obj;
    })

    let inputs = {}
    for (input of inputsArray) {
        Object.assign(inputs, input)
    }

    let variables = {
        'inputs': inputs
    }

    let steps = {}
    for (step of formula.steps) {
        let method = 'GET'
        if (step.method) {
            method = step.method
        }
        console.log(clc.xterm(8)('Step: ') + clc.green(step.name) + clc.xterm(8)(' - Method: ') + clc.yellow(method))

        if (inputs.rit_verbose == 'true') {
            console.log(clc.xterm(8)('URL template:'), clc.xterm(235)(step?.url))
        }
        let urlReplaced = replaceVars(step.url, variables)
        if (inputs.rit_verbose == 'true') {
            console.log(clc.xterm(8)('URL replaced:'), clc.xterm(235)(urlReplaced))
        }

        let dataReplaced = []
        if (step?.data) {
            if (inputs.rit_verbose == 'true') {
                console.log(clc.xterm(8)('Data template:'), (step.data))
            }
            dataReplaced = (Object.keys(step.data).length) ? replaceVarsArray(step.data, variables) : []
            if (inputs.rit_verbose == 'true') {
                console.log(clc.xterm(8)('Data replaced:'), (dataReplaced))
            }
        }

        let headersReplaced = []
        if (step?.headers) {
            if (inputs.rit_verbose == 'true') {
                console.log(clc.xterm(8)('Headers template:'), (step.headers))
            }
            headersReplaced = (Object.keys(step.headers).length) ? replaceVarsArray(step.headers, variables) : []
            if (inputs.rit_verbose == 'true') {
                console.log(clc.xterm(8)('Headers replaced:'), (headersReplaced))
            }
        }

        const objAxios = {
            method: method,
            url: urlReplaced,
            headers: headersReplaced,
            data: dataReplaced
        }
        const response = await axios(objAxios).catch((err) => {
            // ToDo: Abort on Fail ?
            if (err?.response?.status) {
                console.log('Error', err.response.status, err.response.statusText)
            } else {
                console.log('Error', err)
            }
        })
        if (response?.data) {
            let outputDataPath = 'response.data'
            if (step?.output?.datapath) {
                outputDataPath += '.' + step.output.datapath
            }
            const outputData = eval(outputDataPath)
            if (step?.name) {
                steps[step.name] = outputData
                variables['steps'] = steps
            }
            if (step?.output?.format == 'json') {
                console.log(outputData)
            }
            if (step?.output?.format == 'table') {
                console.table(outputData)
            }
        }

        console.log(clc.xterm(208)("\n-------------------------------------------------------------------------------------------------\n"))
    }
}

const replaceVarsArray = (obj, variables) => {
    var output = {}
    for (i in obj) {
        if (Object.prototype.toString.apply(obj[i]) === '[object Object]') {
            output[i] = replaceVarsArray(obj[i], variables)
        } else {
            output[i] = replaceVars(obj[i], variables)
        }
    }
    return output
}

const replaceVars = (string, variables) => {
    return string.replace(/\${{([^{}]+)}}/g, (keyExpr, key) => {
        return eval(('variables.' + key.toLowerCase()).replace(/\./g, '?.')) || keyExpr
        }
    )
}

module.exports = formula
