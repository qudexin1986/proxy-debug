//Proxy Debug
//This simple program is for helping developers debug through http header.
//For more detail, see README.md

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//color config
var color map[string]interface{}

//Parse config file
func readConfig() {
	config, err := os.Open("config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	buf := bufio.NewReader(config)
	line, _ := buf.ReadString('\n')

	var jsonData interface{}
	err = json.Unmarshal([]byte(line), &jsonData)
	if err != nil {
		log.Fatalln(err)
	}
	var ok bool
	color, ok = jsonData.(map[string]interface{})
	if ok == false {
		log.Fatalln("Parse config file error, it must be a json string!")
	}
	for _, c := range color {
		if c.(float64) > 37 || c.(float64) < 30 {
			log.Fatalln("Config error!The valid value is 30-37.")
		}
	}
	item := [5]string{"url", "varName", "varType", "varValue", "group"}
	for _, i := range item {
		_, has := color[i]
		if has == false {
			log.Fatalln("Losing configuration:", i)
		}
	}
}

func main() {
	var port int = 8888

	if len(os.Args) == 1 {
		fmt.Println("Listening in default port:8888")
	} else if os.Args[1] == "--help" {
		fmt.Println("usage: proxy [-p port]")
		return
	} else if len(os.Args) != 3 || os.Args[1] != "-p" {
		log.Fatalln("Error arguments!Just support '-p port'.")
	} else {
		port, err := strconv.Atoi(os.Args[2])
		if err != nil && port > 65535 || port < 1024 {
			log.Fatalln("Error port, it should be 1024-65535, default is 8888.")
		}
	}

	readConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = ""
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer resp.Body.Close()

		//Get the debugging information form http header
		caterory := map[string]int{"I": 34, "W": 33, "E": 31}
		format := "\t \033[%dm-%s\033[%vm %s\033[%vm%s\033[%vm%s\n"
		debugItem := make(map[string]map[string]interface{})
		debugItemIndex := make([]string, 0, 5)
		var jsonData interface{}

		v, okDebugItem := resp.Header["Proxy_debug_item_count"]
		if okDebugItem {
			count, _ := strconv.Atoi(v[0])
			for i := 1; i <= count; i++ {
				index := "Proxy_debug_item_" + strconv.Itoa(i)
				vv, ok := resp.Header[index]
				if ok {
					err = json.Unmarshal([]byte(vv[0]), &jsonData)
					if err != nil {
						continue
					}
					data, ok := jsonData.(map[string]interface{})
					if ok == false {
						continue
					}
					debugItemIndex = append(debugItemIndex, index)
					debugItem[index] = data
				}
			}
		}

		debugGroup := make(map[string]interface{})
		debugGroupIndex := make([]string, 0, 5)

		v, okDebugGroup := resp.Header["Proxy_debug_group_count"]
		if okDebugGroup {
			count, _ := strconv.Atoi(v[0])
			for i := 1; i <= count; i++ {
				index := "Proxy_debug_group_" + strconv.Itoa(i)
				vv, ok := resp.Header[index]
				if ok {
					err = json.Unmarshal([]byte(vv[0]), &jsonData)
					if err != nil {
						continue
					}
					debugGroup[index] = jsonData
					debugGroupIndex = append(debugGroupIndex, index)
				}
			}
		}

		//response to browser
		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		for _, c := range resp.Cookies() {
			w.Header().Add("Set-Cookie", c.Raw)
		}
		w.WriteHeader(resp.StatusCode)
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Write(result)

		//output debugging information
		if okDebugItem {
			fmt.Printf("\033[%vm%v\n", color["url"], r.URL)

			var maxLenName, maxLenType int = 0, 0
			for _, vm := range debugItem {
				v := vm["name"].(string)
				if len(v) > maxLenName {
					maxLenName = len(v)
				}
				v = vm["type"].(string)
				if len(v) > maxLenType {
					maxLenType = len(v)
				}
			}

			for _, i := range debugItemIndex {
				n := debugItem[i]["name"].(string)
				t := debugItem[i]["type"].(string)
				c := debugItem[i]["category"].(string)
				fmt.Printf(
					format,
					caterory[c],
					c,
					color["varName"],
					n+strings.Repeat(" ", maxLenName-len(n)+1),
					color["varType"],
					t+strings.Repeat(" ", maxLenType-len(t)+1),
					color["varValue"],
					strings.Replace(fmt.Sprint(debugItem[i]["value"]), "map", "", 1))
			}
		}

		if okDebugGroup {
			if okDebugItem == false {
				fmt.Printf("\033[%vm%v\n", color["url"], r.URL)
			}
			maxLenName := make([]int, len(debugGroupIndex))
			maxLenType := make([]int, len(debugGroupIndex))
			k := 0
			for _, vm := range debugGroup {
				for _, vv := range vm.([]interface{}) {
					vk, ok := vv.(map[string]interface{})
					if ok == false {
						continue
					}
					v := vk["name"].(string)
					if len(v) > maxLenName[k] {
						maxLenName[k] = len(v)
					}
					v = vk["type"].(string)
					if len(v) > maxLenType[k] {
						maxLenType[k] = len(v)
					}
				}
				k++
			}

			k = 0
			for _, i := range debugGroupIndex {
				fmt.Printf("\t\033[%vm=Group %v=\n", color["group"], k+1)
				for _, v := range debugGroup[i].([]interface{}) {
					vk, ok := v.(map[string]interface{})
					if ok == false {
						continue
					}
					n := vk["name"].(string)
					t := vk["type"].(string)
					c := vk["category"].(string)
					fmt.Printf(
						format,
						caterory[c],
						c,
						color["varName"],
						n+strings.Repeat(" ", maxLenName[k]-len(n)+1),
						color["varType"],
						t+strings.Repeat(" ", maxLenType[k]-len(t)+1),
						color["varValue"],
						strings.Replace(fmt.Sprint(vk["value"]), "map", "", 1))
				}
				k++
				fmt.Printf("\t\033[%vm=GROUP=\n", color["group"])
			}
		}
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
