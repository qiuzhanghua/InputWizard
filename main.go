package main

import (
	"encoding/xml"
	"fmt"
	"github.com/qiuzhanghua/go-input"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"strconv"
)

type Step struct {
	XMLName   xml.Name    `xml:"step"`
	Id        int         `xml:"id,attr"` // id first start should be 0
	Name      string      `xml:"name,attr"`
	ShowMsg   string      `xml:"show-msg"` // 显示的已经国际化的内容
	Default   string      `xml:"default"`
	Required  bool        `xml:"required"`
	Options   []string    `xml:"options"`    // 用户可选的选项
	Collected interface{} `xml:"collected"`  // 获取的输入
	CollectTo string      `xml:"collect-to"` // 搜集到Map的那个key下
	NextId    int         `xml:"next-id"`    // if NextID is -1, Wizard will end.
	NextJs    string      `xml:"next-js"`    // javascript to cal NextId, prefer
}

type Wizard struct {
	XMLName xml.Name `xml:"wizard"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Lang    string   `xml:"lang,attr"` // such as en-US, zh-CN, zh etc
	Step    []Step   `xml:"step"`
}

func main() {
	w := Wizard{}
	buffer, err := ioutil.ReadFile("tdg.zh.xml")
	checkError(err)
	err = xml.Unmarshal(buffer, &w)
	checkError(err)

	ui := input.DefaultUI()
	nextId := 0
	vm := otto.New()
	result := make(map[string]interface{})

	for nextId >= 0 {
		index := -1
		for i := range w.Step {
			if w.Step[i].Id == nextId {
				index = i
				break
			}
		}
		if index == -1 {
			nextId = -1
		} else {
			step := w.Step[index]
			if len(step.Options) > 0 {
				lang, err := ui.Select(step.ShowMsg, step.Options, &input.Options{
					Default: step.Default,
					Loop:    true,
				})
				checkError(err)
				step.Collected = lang
			} else {
				// TODO add input and other
			}

			if step.Required {
				result[step.CollectTo] = step.Collected
			}

			// 按照当前的xml，选择Java会死循环，选择Golang or Rust会成功
			if len(step.NextJs) > 6 {
				vm.Set("id", step.Id)
				vm.Set("option", step.Collected)
				val, err := vm.Run(step.NextJs)
				checkError(err)
				x, err := strconv.Atoi(val.String())
				checkError(err)
				nextId = x
			} else {
				nextId = step.NextId
			}
		}
	}

	//存储的结果
	fmt.Println(result)
}

func showSampleXm() {
	welcome := Step{
		Id:      0,
		Name:    "welcome",
		ShowMsg: "欢迎使用TDP",
		NextId:  10,
	}
	first := Step{
		Id:        10,
		Name:      "select first",
		ShowMsg:   "请选择你喜欢的语言",
		Options:   []string{"Java", "Go", "Rust"},
		Required:  true,
		CollectTo: "lang",
		NextId:    -1,
		NextJs:    `if (option == "Java") nextId = 20; else nextId = -1;`,
	}
	second := Step{
		Id:      20,
		Name:    "show second",
		ShowMsg: "流程结束",
		NextId:  -1,
	}

	wizard := Wizard{
		Name: "cmd",
		Lang: "zh",
		Step: []Step{welcome, first, second},
	}

	out, _ := xml.MarshalIndent(&wizard, " ", "  ")
	fmt.Println(string(out))

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
