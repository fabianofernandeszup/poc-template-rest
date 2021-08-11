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
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	dynamicPattern    = `\${{([^{}]+)}}`
	expressionPattern = `(\$\[\[.*\]\])`
)

type (
	FormulaYAML struct {
		Name               string  `yaml:"name"`
		Description        string  `yaml:"description"`
		Template           string  `yaml:"template"`
		TemplateRelease    string  `yaml:"templateRelease"`
		DockerImageBuilder string  `yaml:"dockerImageBuilder"`
		Plugin             string  `yaml:"plugin"`
		Inputs             []Input `yaml:"inputs"`
		Steps              []Step  `yaml:"steps"`
	}

	Input struct {
		Label   string   `yaml:"label"`
		Name    string   `yaml:"name"`
		Default string   `yaml:"default"`
		Type    string   `yaml:"type"`
		Items   []string `yaml:"items,omitempty"`
	}

	Step struct {
		Name    string            `yaml:"name"`
		Method  string            `yaml:"method"`
		URL     string            `yaml:"url"`
		Output  map[string]string `yaml:"output"`
		Headers map[string]string `yaml:"headers,omitempty"`
		Data    map[string]string `yaml:"data,omitempty"`
	}
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Run() {
	if Exists("formula.yml") {
		yamlFile, err := ioutil.ReadFile("formula.yml")
		y := FormulaYAML{}
		err = yaml.Unmarshal([]byte(yamlFile), &y)
		check(err)
		sm := make(map[string]map[string]string)
		for _, step := range y.Steps {
			fmt.Println("\n- - - - - - - - - - Step " + step.Name + " - - - - - - - - - - ")
			// fmt.Println("\nName:", step.Name)
			// fmt.Println("Method:", step.Method)
			// fmt.Println("Url:", step.URL)
			// fmt.Println("Headers:", step.Headers)
			// fmt.Println("Data:", step.Data)
			// fmt.Println("Output:", step.Output)
			step.URL = convertURLDynamicValues(step.URL, sm)
			if len(step.Data) != 0 {
				step.Data = convertMapDynamicValues(step.Data, sm, "data")
			}
			if len(step.Headers) != 0 {
				step.Headers = convertMapDynamicValues(step.Headers, sm, "headers")
			}
			response := consumeAPI(step.Method, step.URL, step.Headers, step.Data, step.Output)
			sm[step.Name] = response
		}
	} else {
		fmt.Println("ERROR: formula.yml file not found")
	}
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func consumeAPI(method string, url string, headers map[string]string, data map[string]string, output map[string]string) map[string]string {
	client := &http.Client{}

	var jsonStr []byte
	if len(data) != 0 {
		j, err := json.Marshal(data)
		check(err)
		jsonStr = j
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(jsonStr))
	check(err)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	check(err)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	check(err)

	fmt.Println("\nRequest StatusCode:", resp.StatusCode)

	for k, v := range output {
		if k == "format" {
			if v == "table" {
				// Can't use https://github.com/gosuri/uitable or https://github.com/lensesio/tableprinter without response struct
				// fmt.Println(string(bodyBytes))
			}
			if v == "json" {
				fmt.Println(string(bodyBytes))
			}
		}
	}

	m := make(map[string]string)
	err = json.Unmarshal(bodyBytes, &m)
	if err != nil {
		return nil
	}
	return m
}

func convertMapDynamicValues(data map[string]string, sm map[string]map[string]string, mapType string) map[string]string {
	for k, v := range data {
		if isExpression(data[k], sm) {
			data[k] = updateDynamicValues(data[k], sm)
			re2 := regexp.MustCompile(`\$\[\[([^{}]+)\]\]`) // Inside Expression Pattern
			match2 := re2.FindStringSubmatch(data[k])
			data[k] = match2[1]
		} else {
			if strings.Contains(v, "${{") {
				v = strings.Replace(v, "${{", "", -1)
				v = strings.Replace(v, "}}", "", -1)
				data[k] = convertDynamicValues(v, sm)
			}
		}
	}
	if len(data) != 0 {
		if mapType == "data" {
			fmt.Println("\nUpdated DATA:", data)
		}
		if mapType == "headers" {
			fmt.Println("\nUpdated HEADERS:", data)
		}
	}
	return data
}

func convertURLDynamicValues(url string, sm map[string]map[string]string) string {
	if strings.Contains(url, "${{") {
		url = updateDynamicValues(url, sm)
		fmt.Println("\nUpdated URL:", url)
	}
	return url
}

func updateDynamicValues(param string, sm map[string]map[string]string) string {
	re := regexp.MustCompile(dynamicPattern)
	match := re.FindAllStringSubmatch(param, 10)
	//fmt.Println("\nMatch:", match)
	if len(match) != 0 {
		for i := 0; i < len(match); i++ {
			key := match[i][0]
			value := match[i][1]
			newValue := convertDynamicValues(value, sm)
			param = strings.Replace(param, key, newValue, -1)
		}
	}
	return param
}

func isExpression(field string, sm map[string]map[string]string) bool {
	//fmt.Println("\nField:", field)
	re1 := regexp.MustCompile(expressionPattern)
	match := re1.FindStringSubmatch(field)
	if len(match) == 0 {
		return false
	} else {
		return true
	}
}

func convertDynamicValues(value string, sm map[string]map[string]string) string {
	sv := strings.Split(value, ".")
	//fmt.Println("\nSplitted Values:", sv)
	if len(sv) != 0 {
		if sv[0] == "inputs" {
			localVariableName := sv[1]
			value = os.Getenv(strings.ToUpper(localVariableName))
			//fmt.Println("Step Data New Value (Input):", value)
		}
		if sv[0] == "steps" {
			// TODO: Manage to get values on deepest level (currently only work for ${{step.1.2}}). e.g: ${{step.1.2.3.4}}
			for k, _ := range sm {
				if k == sv[1] {
					ms := sm[k]
					value = ms[sv[2]]
					//fmt.Println("Step Data New Value (Step):", value)
				}
			}
		}
	} else {
		fmt.Println("ERROR: Splitting value")
	}
	return value
}
