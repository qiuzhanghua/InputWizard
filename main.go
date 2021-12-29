package main

import (
	"InputWizard/utils"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/qiuzhanghua/go-input"
	"github.com/qiuzhanghua/i10n"
	"github.com/robertkrimen/otto"
	"log"
	"strconv"
	"strings"
)

type Step struct {
	XMLName    xml.Name    `xml:"step"`
	Id         int         `xml:"id,attr"` // id first start should be 0
	Name       string      `xml:"name,attr"`
	ShowMsg    string      `xml:"show-msg"` // 显示的已经国际化的内容
	Default    string      `xml:"default"`
	Required   bool        `xml:"required"`
	Masked     bool        `xml:"masked"`  // 如果是密码
	Options    []string    `xml:"options"` // 用户可选的选项
	OptionFunc string      `xml:"option-func"`
	Collected  interface{} `xml:"collected"`  // 获取的输入
	CollectTo  string      `xml:"collect-to"` // 搜集到Map的那个key下
	NextId     int         `xml:"next-id"`    // if NextID is -1, Wizard will end.
	NextJs     string      `xml:"next-js"`    // javascript to cal NextId, prefer
}

type Wizard struct {
	XMLName xml.Name `xml:"wizard"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Lang    string   `xml:"lang,attr"` // such as en-US, zh-CN, zh etc
	Step    []Step   `xml:"step"`
}

func GetDirs() []string {
	return []string{"Go", "Node", "JavaScript"}
}

func init() {
	utils.StubStorage["GetDirs"] = GetDirs
	//ret, err := utils.Call("GetDirs")
	//fmt.Println(ret, err)
	//_ = i10n.SetDefaultLang("en-US")
	_ = i10n.SetDefaultLang("zh-CN")
	for _, name := range AssetNames() {
		if strings.HasPrefix(name, "locales") && strings.HasSuffix(name, ".properties") {
			buffer, err := Asset(name)
			if err != nil {
				log.Fatal(err)
			}
			p, err := properties.Load(buffer, properties.UTF8)
			if err != nil {
				log.Fatal(err)
			}
			tag := i10n.ParseTagWithDefault(name)
			i10n.AddTagMap(tag, p.Map())
		}
	}
	// 使用i10的方法
	input.T = i10n.T

}

func main() {
	w := Wizard{}
	prefix := "wizards/tdg."
	suffix := ".xml"
	var languages []string
	for _, lang := range AssetNames() {
		if strings.HasPrefix(lang, prefix) {
			index := strings.LastIndex(lang, suffix)
			if index > 0 {
				s := lang[len(prefix):index]
				languages = append(languages, s)
			}
		}
	}
	lang := i10n.Nearest(i10n.GetDefaultTag().String(), languages)
	buffer, err := Asset(fmt.Sprintf("%s%s%s", prefix, lang, suffix))
	//buffer, err := ioutil.ReadFile("tdg.zh.xml")
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
			if len(step.CollectTo) == 0 {
				fmt.Println(step.ShowMsg)
			} else if len(step.Options) > 0 || len(step.OptionFunc) > 0 {
				if len(step.OptionFunc) > 0 {
					ret, err := utils.Call(step.OptionFunc)
					checkError(err)
					step.Options = ret.([]string)
					if len(step.Options) == 0 {
						// TODO 按照业务调整
						err := errors.New("can't find dirs")
						checkError(err)
					}
				}
				lang, err := ui.Select(step.ShowMsg, step.Options, &input.Options{
					Default: step.Default,
					Loop:    true,
				})
				checkError(err)
				step.Collected = lang
			} else {
				name, err := ui.Ask(step.ShowMsg, &input.Options{
					Default:  step.Default,
					Mask:     step.Masked,
					Required: step.Required,
					Loop:     true,
				})
				checkError(err)
				step.Collected = name
			}

			if step.Required && len(step.CollectTo) > 0 {
				result[step.CollectTo] = step.Collected
			}

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

	fmt.Println("result Map:", result)
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
