package main

import (
	"crawler-andlabs/crawler"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/hpcloud/tail"
	"io"
	"log"
	"os"
	"time"
)
var outputEntry *ui.MultilineEntry
var (
	logFileName = "temp_output.log"
)

func logConfig() {
	logFile, logErr := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", logFileName)
		os.Exit(1)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//清空文件内容
	f, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
}

func main() {
	logConfig()

	err := ui.Main(func() {
		// 生成：文本框
		paramsEntry := ui.NewMultilineEntry()
		defaultParams := "-rulePageUrl=https://www.tupianzj.com/mingxing/xiezhen/20130730/3492_[99,100].html\n-regularImgUrl=https://img.lianzhixiu.com/uploads/.*?.jpg"

		paramsEntry.SetText(defaultParams)
		// 生成：输出控制台
		outputEntry = ui.NewNonWrappingMultilineEntry()
		outputEntry.SetReadOnly(true)
		defaultOutputText := `说明：
	-c    自定义代码执行
	-r string
		referer
	-regularImgUrl string
		正则图片url，如：https://img.lianzhixiu.com/uploads/allimg/.*?.jpg
	-ruleImgUrl string
		规则图片url，如：https://img-pre.ivsky.com/img/tupian/pre/202107/15/wenquan-[1,3].jpg 或者https://img-pre.ivsky.com/img/tupian/pre/202107/15/wenquan-[001,003].jpg
	-rulePageUrl string
		规则页面url，如：https://www.tupianzj.com/meinv/20201102/219671_[1,2].html，需要配合regularImgUrl使用
	-ua
		是否设置user-agent`
		outputEntry.SetText(defaultOutputText)
		// 生成：按钮
		button := ui.NewButton(`点击爬取`)
		// 设置：按钮点击事件
		button.OnClicked(func(*ui.Button) {
			go crawler.StartByAndlabs(paramsEntry.Text(), outputEntry)

			//outputEntry.SetText(`你好，` + paramsEntry.Text() + `！`)
		})
		// 生成：垂直容器
		box := ui.NewVerticalBox()
		// 往 垂直容器 中添加 控件
		box.Append(ui.NewLabel(`执行参数：`), false)
		box.Append(paramsEntry, true)
		box.Append(button, false)
		box.Append(outputEntry, true)

		// 生成：窗口（标题，宽度，高度，是否有 菜单 控件）
		window := ui.NewWindow(`crawler`, 1000, 600, false)
		// 窗口容器绑定
		window.SetChild(box)

		// 设置：窗口关闭时
		window.OnClosing(func(*ui.Window) bool {
			// 窗体关闭
			ui.Quit()
			return true
		})
		// 窗体显示
		window.Show()
		go outputEntryRender()
	})
	if err != nil {
		panic(err)
	}
}

func outputEntryRender() {
	go func() {
		fmt.Println("进来-----------")
		config := tail.Config{
			ReOpen:    true,                                 // 重新打开
			Follow:    true,                                 // 是否跟随
			Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
			MustExist: false,                                // 文件不存在不报错
			Poll:      true,
		}
		tails, err := tail.TailFile(logFileName, config)
		if err != nil {
			fmt.Println("tail file failed, err:", err)
			return
		}
		var (
			line *tail.Line
			ok   bool
			loopCount int
		)
		for {
			line, ok = <-tails.Lines
			if !ok {
				fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
				time.Sleep(time.Second)
				continue
			}
			loopCount++
			if loopCount > 300 {
				outputEntry.SetText("")
				loopCount = 0
			}
			outputEntry.Append(line.Text+"\n")
			//fmt.Println("line:", line.Text)
		}

		return
	}()

}