package view

import (
	"brotherGame/static/image"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// 程序主视图
type MainView struct {
	Window   fyne.Window
	MainMenu *fyne.MainMenu
	File     *fyne.Menu
	Help     *fyne.Menu
	Settings *fyne.MenuItem
}

func NewMainView() *MainView {
	m := new(MainView)
	return m
}

func (m *MainView) LoadUI(app fyne.App) {
	m.Settings = fyne.NewMenuItem("Settings", m.SetUp)
	m.File = fyne.NewMenu("文件", fyne.NewMenuItemSeparator(), m.Settings)
	m.Help = fyne.NewMenu("帮助", fyne.NewMenuItem("Help", m.HelpShow))
	m.MainMenu = fyne.NewMainMenu(m.File, m.Help)

	m.Window = app.NewWindow("老弟游戏")
	m.Window.SetIcon(image.GameLogo())
	m.Window.Resize(fyne.NewSize(800, 480))
	m.Window.SetMainMenu(m.MainMenu)
	m.Window.SetMaster()

	game := new(GameModule)
	game.LoadGame(m.Window)

	tabs := widget.NewTabContainer(widget.NewTabItemWithIcon("Game", theme.DocumentCreateIcon(), game.GameCanvas))
	tabs.SetTabLocation(widget.TabLocationLeading)
	tabs.SelectTabIndex(app.Preferences().Int("currentTab"))

	m.Window.SetContent(tabs)
	m.Window.ShowAndRun()
	app.Preferences().SetInt("currentTab", tabs.CurrentTabIndex())
}

func (m *MainView) HelpShow() {
	dialog.ShowInformation("", "垃圾！这都不会玩吗？\n(((┏(;￣▽￣)┛装完逼就跑", m.Window)
}
func (m *MainView) SetUp() {
	dialog.ShowInformation("", "敬请期待", m.Window)
}
