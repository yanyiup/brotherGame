package main

import (
	"brotherGame/view"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"os"
)

func main() {
	//设置字体环境变量
	os.Setenv("FYNE_FONT", "./static/emoji.ttf")
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	v := view.NewMainView()
	v.LoadUI(a)
}
