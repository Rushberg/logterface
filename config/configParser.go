package config

import (
	"encoding/json"
	"fmt"
	"logterface/handlers"
	"logterface/layouts"
	"logterface/utils"
	"os"
)

type LogHandlerConfig struct {
	Type   string                 `json:"type"`
	Id     string                 `json:"id"`
	RegEx  string                 `json:"regex"`
	Params map[string]interface{} `json:"params"`
}

type LogHandlerLayoutConfig struct {
	Id     string                 `json:"id"`
	Params map[string]interface{} `json:"params"`
}

type LayoutConfig struct {
	Type     string                   `json:"type"`
	Params   map[string]interface{}   `json:"params"`
	Handlers []LogHandlerLayoutConfig `json:"handlers"`
}

type Config struct {
	Refresh  int                `json:"refresh_mills"`
	Handlers []LogHandlerConfig `json:"handlers"`
	Layouts  []LayoutConfig     `json:"layouts"`
}

func ParseConfig(filePath string) (*handlers.HandlerManager, *layouts.LayoutManager, int) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, nil, 0
	}

	// Step 2: Parse the JSON data into a struct
	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, nil, 0
	}

	hm := handlers.NewHandlerManager()
	handlersMap := map[string]handlers.LogHandler{}
	for _, hc := range config.Handlers {
		switch hc.Type {
		case "Numbers":
			methodName, _ := hc.Params["method"].(string)
			method, _ := handlers.MethodFromString(utils.Capitalize(methodName))
			name, _ := hc.Params["name"].(string)
			lh := handlers.NewNumbersHandler(name, hc.RegEx, method)
			thm, exists := hc.Params["thresholdMethod"]
			if exists {
				lh.ThresholdMethod, _ = handlers.ThresholdMethodFromString(utils.Capitalize(thm.(string)))
				lh.Threshold, _ = hc.Params["threshold"].(float64)
			}

			hm.AddHandler(lh)
			handlersMap[hc.Id] = lh
		case "Graph":
			name, _ := hc.Params["name"].(string)
			width, _ := hc.Params["width"].(float64)
			height, _ := hc.Params["height"].(float64)
			gh := handlers.NewGraphHandler(name, hc.RegEx, int(width), int(height))
			hm.AddHandler(gh)
			handlersMap[hc.Id] = gh
		case "Progress":
			width, _ := hc.Params["width"].(float64)
			name, _ := hc.Params["name"].(string)
			ph := handlers.NewProgressHandler(name, hc.RegEx, int(width))
			dt, exists := hc.Params["defaultTotal"]
			if exists {
				ph.DefaultTotalValue = dt.(float64)
			}
			rt, exists := hc.Params["regexTotal"]
			if exists {
				ph.RegexTotalValue = rt.(string)
			}
			hm.AddHandler(ph)
			handlersMap[hc.Id] = ph
		case "Filter":
			name, _ := hc.Params["name"].(string)
			width, _ := hc.Params["width"].(float64)
			length, _ := hc.Params["length"].(float64)
			fh := handlers.NewFilterHandler(name, hc.RegEx, int(length), int(width))
			hm.AddHandler(fh)
			handlersMap[hc.Id] = fh
		}
	}
	lm := layouts.NewLayoutManager()

	for _, lc := range config.Layouts {
		switch lc.Type {
		case "Pipe":
			lm.AddPipe(&hm)
		case "Line":
			width, _ := lc.Params["width"].(float64)
			ll := layouts.NewLineLayout(int(width))
			for _, hlc := range lc.Handlers {
				ll.AddHandler(handlersMap[hlc.Id])
			}
			lm.AddLayout(ll)
		}
	}
	return &hm, &lm, config.Refresh

}
