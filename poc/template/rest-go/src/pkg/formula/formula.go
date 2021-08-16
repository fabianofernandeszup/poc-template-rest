// This is the formula implementation class.
// Where you will code your methods and manipulate the inputs to perform the specific operation you wish to automate.

package formula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"

	"gopkg.in/yaml.v2"
)

const (
	dynamicPattern    = `\{{([^{}]+)}}`
	expressionPattern = `(\[\[.*\]\])`
)

type (
	FormulaYAML struct {
		Name               string    `yaml:"name"`
		Description        string    `yaml:"description"`
		Template           string    `yaml:"template"`
		TemplateRelease    string    `yaml:"templateRelease"`
		DockerImageBuilder string    `yaml:"dockerImageBuilder"`
		Plugin             string    `yaml:"plugin"`
		Inputs             []Input   `yaml:"inputs"`
		Execution          Execution `yaml:"execution"`
	}

	Input struct {
		Label   string   `yaml:"label"`
		Name    string   `yaml:"name"`
		Default string   `yaml:"default"`
		Type    string   `yaml:"type"`
		Items   []string `yaml:"items,omitempty"`
	}

	Execution struct {
		Workflow string `yaml:"workflow"`
		Steps    []Step `yaml:"steps"`
	}

	Step struct {
		Name    string                 `yaml:"name"`
		Method  string                 `yaml:"method"`
		URL     string                 `yaml:"url"`
		Headers map[string]string      `yaml:"headers,omitempty"`
		Data    map[string]interface{} `yaml:"data,omitempty"`
		Output  map[string]interface{} `yaml:"output"`
	}
)

func check(e error, m string) {
	if e != nil {
		if m != "" {
			fmt.Println(m)
		}
		panic(e)
	}
}

func TypeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func Verbose(v bool, i interface{}, m string) {
	if v {
		if m != "" {
			fmt.Println("\n" + m)
		}
		fmt.Println(i)
	}
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func FormulaYAMLValid(filePath string) bool {
	if !Exists(filePath) {
		return true
	}
	f, _ := os.Open(filePath)
	fi, err := f.Stat()
	if err != nil {
		return true
	}

	if fi.Size() == 0 {
		return true
	}

	return false
}

func Run() {
	if !FormulaYAMLValid("formula.yml") {
		yamlFile, err := ioutil.ReadFile("formula.yml")
		y := FormulaYAML{}
		err = yaml.Unmarshal([]byte(yamlFile), &y)
		check(err, "ERROR: formula.yml file incorrect format")

		verbose, _ := strconv.ParseBool(os.Getenv("RIT_VERBOSE"))
		sm := make(map[string]string)

		fmt.Println("\033[0;34m\n🤖 START WORKFLOW EXECUTION ⚙️\033[0m")

		for _, step := range y.Execution.Steps {

			fmt.Println("\033[1;34m\n🤖 " + strings.ToUpper(step.Name) + " STEP\033[0m")

			if step.URL != "" {
				step.URL = CheckDynamicValues(step.URL, sm).(string)
				Verbose(verbose, step.URL, "\033[1mURL 👀\033[0m")
			}

			if len(step.Headers) != 0 {
				step.Headers = CheckDynamicValues(step.Headers, sm).(map[string]string)
				Verbose(verbose, step.Headers, "\033[1mHEADERS 👀\033[0m")
			}

			if len(step.Data) != 0 {
				step.Data = CheckDynamicValues(step.Data, sm).(map[string]interface{})
				Verbose(verbose, step.Data, "\033[1mDATA 👀\033[0m")
			}

			response := CallAPI(step.Method, step.URL, step.Headers, step.Data, step.Output)
			sm[step.Name] = response
		}
		fmt.Println("\033[0;34m\n🤖 END OF WORKFLOW EXECUTION 🚀\033[0m")
	} else {
		fmt.Println("ERROR: Invalid formula.yml file.")
	}
}

func CallAPI(method string, url string, headers map[string]string, data map[string]interface{}, output map[string]interface{}) string {
	client := &http.Client{}

	var jsonStr []byte
	if len(data) != 0 {
		j, err := json.Marshal(data)
		check(err, "ERROR: step.data incorrect format")
		jsonStr = j
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(jsonStr))
	check(err, "ERROR: Building API request")

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	check(err, "ERROR: API call unexpected error")

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	check(err, "ERROR: Couldn't read API response")

	if resp.StatusCode >= 300 {
		fmt.Println("🔴 API Response StatusCode:", resp.StatusCode)
		fmt.Println("ERROR: Please, check the URL and datas sent (with the 'Verbose' input) or if the service is available.")
		defer func() {
			if panicMessage := recover(); panicMessage != nil {
				//DoNothing
			}
		}()
		panic("")
	} else {
		fmt.Println("\n🟢 API Response StatusCode:", resp.StatusCode)
		if len(bodyBytes) > 0 {
			jsonParsed, err := gabs.ParseJSON(bodyBytes)
			check(err, "ERROR: Couldn't store API response")
			CheckOutput(bodyBytes, output)
			return jsonParsed.String()
		}
	}

	return ""
}

func CheckOutput(bodyBytes []byte, output map[string]interface{}) {
	for k, v := range output {
		if k == "format" {
			if v == "table" {
				// Can't use https://github.com/gosuri/uitable or https://github.com/lensesio/tableprinter without response struct
			}
			if v == "json" {
				fmt.Println(string(bodyBytes))
			}

			if v == "template" {
				output["template"] = GetOutputValues(bodyBytes, output["template"])
				fmt.Println("\nOUTPUT:")
				fmt.Println(output["template"])
			}
		}
	}
}

func GetOutputValues(bodyBytes []byte, ot interface{}) interface{} {
	var i interface{}
	var value string

	if TypeOf(ot) == "map[interface {}]interface {}" {
		i := make(map[string]interface{})
		for k, v := range ot.(map[interface{}]interface{}) {
			i[k.(string)] = GetOutputValues(bodyBytes, v)
		}
		return i
	}

	if TypeOf(ot) == "string" {
		path := ot.(string)
		jsonParsed, _ := gabs.ParseJSON(bodyBytes)
		value = jsonParsed.Path(path).String()
		value = strings.Replace(value, "\"", "", -1)
		return value
	}

	return i
}

func CheckDynamicValues(data interface{}, sm map[string]string) interface{} {
	var d interface{}

	if TypeOf(data) == "map[string]interface {}" {
		d, _ := data.(map[string]interface{})
		for k, _ := range d {
			d[k] = CheckDynamicValues(d[k], sm)
		}
		d = ConvertStructToMapStringInterface(d)
		return d
	}

	if TypeOf(data) == "map[string]string" {
		d, _ := data.(map[string]string)
		for k, _ := range d {
			if ContainsDynamicValues(d[k]) {
				d[k] = UpdateDynamicValues(d[k], sm)
			}
		}
		return d
	}

	if TypeOf(data) == "string" {
		d, _ := data.(string)
		if ContainsDynamicValues(d) {
			d = UpdateDynamicValues(d, sm)
		}
		return d
	}

	fmt.Println("\nData Type not support:", TypeOf(data))

	return d
}

func ContainsDynamicValues(field string) bool {
	re := regexp.MustCompile(dynamicPattern)
	match := re.FindAllStringSubmatch(field, 10)
	if len(match) > 0 {
		return true
	} else {
		return false
	}
}

func ConvertStructToMapStringInterface(s map[string]interface{}) map[string]interface{} {
	sv := reflect.ValueOf(s)
	msi := make(map[string]interface{})
	if sv.Kind() == reflect.Map {
		for _, key := range sv.MapKeys() {
			strct := sv.MapIndex(key)
			// fmt.Println("Key:", key.String())
			// fmt.Println("Interface Structure:", strct.Interface())
			msi[key.String()] = strct.Interface()
		}
	}
	return msi
}

func UpdateDynamicValues(param string, sm map[string]string) string {
	re := regexp.MustCompile(dynamicPattern)
	match := re.FindAllStringSubmatch(param, 10)
	if len(match) != 0 {
		for i := 0; i < len(match); i++ {
			key := match[i][0]
			value := match[i][1]
			newValue := ConvertDynamicValues(value, sm)
			param = strings.Replace(param, key, newValue, -1)
		}
	}
	return param
}

func ConvertDynamicValues(value string, sm map[string]string) string {
	sv := strings.Split(value, ".")
	if len(sv) != 0 {
		switch vt := sv[0]; vt {
		case "inputs":
			localVariableName := sv[1]
			value = os.Getenv(strings.ToUpper(localVariableName))
		case "steps":
			for k, _ := range sm {
				if k == sv[1] {
					prefix := "steps." + sv[1] + "."
					path := strings.Replace(value, prefix, "", -1)
					jsonParsedObj, _ := gabs.ParseJSON([]byte(sm[k]))
					value = jsonParsedObj.Path(path).String()
					value = strings.Replace(value, "\"", "", -1)
				}
			}
		default:
			fmt.Println("ERROR: Splitting Dynamic value")
		}
	} else {
		fmt.Println("ERROR: Splitting Dynamic value")
	}
	return value
}

func ReadOutputValue(path string, json string) string {
	jsonParsedObj, _ := gabs.ParseJSON([]byte(json))
	value := jsonParsedObj.Path(path).String()
	value = strings.Replace(value, "\"", "", -1)
	return value
}
