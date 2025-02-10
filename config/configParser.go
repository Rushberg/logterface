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
			hm.AddHandler(lh)
			handlersMap[hc.Id] = lh
		case "Counter":
			name, _ := hc.Params["name"].(string)
			ch := handlers.NewCounterHandler(name, hc.RegEx)
			hm.AddHandler(ch)
			handlersMap[hc.Id] = ch
		case "Graph":
			name, _ := hc.Params["name"].(string)
			length, _ := hc.Params["length"].(float64)
			height, _ := hc.Params["height"].(float64)
			gh := handlers.NewGraphHandler(name, hc.RegEx, int(length), int(height))
			hm.AddHandler(gh)
			handlersMap[hc.Id] = gh
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
		case "Chunk":
			cl := layouts.NewChunkLayout()
			for _, hlc := range lc.Handlers {
				cl.AddHandler(handlersMap[hlc.Id])
			}
			lm.AddLayout(cl)
		}
	}
	return &hm, &lm, config.Refresh

}
