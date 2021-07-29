package gui

import (
	"io/fs"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/wsndshx/LocalizeEpub/core"
)

var myWindow fyne.Window
var anniu bool = true
var button *widget.Button
var Terminal *widget.TextGrid

func init() {
	myApp := app.New()
	// 设置界面主题
	myApp.Settings().SetTheme(&myTheme{})
	// 设置标题
	myWindow = myApp.NewWindow("LocalizeEpub")
	// 设置窗口大小
	myWindow.Resize(fyne.NewSize(800, 400))
}

func Start() {
	core.LogInfo.Println("咕噜咕噜~")
	core.LogInfo.Println("项目发布地址: github.com/wsndshx/LocalizeEpub")
	myWindow.SetContent(container.New(layout.NewGridLayoutWithRows(3), file(), start(), terminal()))
	// 刷新终端
	go func() {
		for {
			log, err := fs.ReadFile(core.LogFS, "log")
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			Terminal.SetText(string(log))
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
	myWindow.ShowAndRun()
}

func terminal() *fyne.Container {
	// 终端区域
	Terminal = widget.NewTextGrid()
	box := container.NewScroll(container.New(layout.NewHBoxLayout(), container.New(layout.NewVBoxLayout(), Terminal)))

	// 卡片
	care := widget.NewCard("", "", box)

	content := container.New(layout.NewPaddedLayout(), care)
	return content

}

func start() *fyne.Container {
	// 选择转换的模式
	model := widget.NewSelect([]string{
		"简体化",
		"繁體化",
		"中国化",
		"香港化",
		"台灣化",
		"维基简体化",
		"維基繁體化"},
		func(s string) {
			core.LogInfo.Printf("你选择了模式[%s]\n", s)
		})

	// 卡片
	care := widget.NewCard("Model", "", model)

	// 按钮
	button = widget.NewButton("转换", func() {
		core.LogInfo.Println("开始转换喵~")
		go func() {
			model.Disable()
			button.Disable()
			err := core.Start(model.Selected)
			if err != nil {
				core.LogError.Println(err)
				dialog.ShowError(err, myWindow)
				model.Enable()
				button.Enable()
				return
			}
			model.Enable()
			button.Enable()
		}()
	})
	// care.Resize(care.MinSize())
	miao := container.NewHBox(care, container.New(layout.NewCenterLayout(), button))
	miao.Resize(miao.MinSize())
	content := container.New(layout.NewCenterLayout(), miao)
	return content
}

func file() *fyne.Container {
	return container.New(layout.NewGridLayoutWithColumns(2), input(), output())
}

func input() *fyne.Container {
	// 放置一个文本框, 用于显示目前选定的文件目录
	input_path := widget.NewEntry()
	// 设置闲置文本
	input_path.SetPlaceHolder("Nya~")
	// 禁止用户手动输入
	input_path.Disable()

	// 用于打开一个窗口, 用于选择需要转换的文件
	fileOpen := dialog.NewFileOpen(func(path fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if path != nil {
			// 修改上面那个文本框的内容
			input_path.SetText(path.URI().Path())
			// 修改变量
			core.Input = path.URI().Path()
		}
	}, myWindow)
	// 限制只能选择.epub文件
	fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".epub"}))

	name := widget.NewLabelWithStyle("打开EPUB文件", 1, fyne.TextStyle{})
	content := container.NewVBox(name, input_path, widget.NewButton("打开", func() {
		fileOpen.Show()
	}))

	return content
}

func output() *fyne.Container {
	// 放置一个文本框, 用于显示目前选定的文件目录
	output_path := widget.NewEntry()
	// 设置闲置文本
	output_path.SetPlaceHolder("Nya~")
	// 禁止用户手动输入
	output_path.Disable()

	// 打开一个对话框, 选择保存的位置
	fileSave := dialog.NewFileSave(func(path fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if path != nil {
			// 修改上面那个文本框的内容
			output_path.SetText(path.URI().Path())
			// 修改变量
			core.Output = path.URI().Path()
		}
	}, myWindow)

	name := widget.NewLabelWithStyle("输出的EPUB文件", 1, fyne.TextStyle{})
	content := container.NewVBox(name, output_path, widget.NewButton("打开", func() {
		fileSave.Show()
	}))
	return content
}
