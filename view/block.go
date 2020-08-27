package view

import (
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const buttonTapDuration = 250

type blockButtonRenderer struct {
	icon    *canvas.Image
	label   *canvas.Text
	button  *BlockButton
	objects []fyne.CanvasObject
	layout  fyne.Layout
}

func (b *blockButtonRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*2, theme.Padding()*2)
}

// MinSize计算按钮的最小大小
// 这是基于包含的文本、任何设置的图标和标准
// 添加的填充量.

func (b *blockButtonRenderer) MinSize() (size fyne.Size) {
	labelSize := b.label.MinSize()
	size = labelSize.Add(b.padding())
	return
}

// 布局按钮小部件的组件
func (b *blockButtonRenderer) Layout(size fyne.Size) {
	hasIcon := b.icon != nil
	hasLabel := b.label.Text != ""
	if !hasIcon && !hasLabel {
		// Nothing to layout
		return
	}
	iconSize := size.Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	labelSize := b.label.MinSize()
	padding := b.padding()
	if hasLabel {
		if hasIcon {
			// Both
			var objects []fyne.CanvasObject
			if b.button.IconPlacement == ButtonIconLeadingText {
				objects = append(objects, b.icon, b.label)
			} else {
				objects = append(objects, b.label, b.icon)
			}
			b.icon.SetMinSize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
			min := b.layout.MinSize(objects)
			b.layout.Layout(objects, min)
			pos := alignedPosition(b.button.Alignment, padding, min, size)
			b.label.Move(b.label.Position().Add(pos))
			b.icon.Move(b.icon.Position().Add(pos))
		} else {
			// Label Only
			b.label.Move(alignedPosition(b.button.Alignment, padding, labelSize, size))
			b.label.Resize(labelSize)
		}
	} else {
		// Icon Only
		b.icon.Move(alignedPosition(b.button.Alignment, padding, iconSize, size))
		b.icon.Resize(iconSize)
	}
}

func alignedPosition(align ButtonAlign, padding, objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case ButtonAlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case ButtonAlignLeading:
		pos.X = padding.Width / 2
	case ButtonAlignTrailing:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
}

// 更新此按钮以匹配当前主题
func (b *blockButtonRenderer) applyTheme() {
	b.label.TextSize = theme.TextSize()
	b.label.Color = theme.TextColor()
	if b.button.Disabled() {
		b.label.Color = theme.DisabledTextColor()
	}
}

// 设置按钮颜色
func (b *blockButtonRenderer) BackgroundColor() color.Color {

	if b.button.ScoreBlock {
		return color.RGBA{128, 128, 128, 1}
	}

	if b.button.MenuBlock {
		switch {
		case b.button.Disabled():
			return theme.DisabledButtonColor()
		case b.button.Style == PrimaryButton:
			return theme.PrimaryColor()
		case b.button.hovered, b.button.tapped:
			return theme.HoverColor()
		default:
			return color.RGBA{151, 87, 66, 1}
		}
	}

	if b.icon != nil {
		switch {
		case b.button.Disabled():
			return theme.DisabledButtonColor()
		case b.button.Style == PrimaryButton:
			return theme.PrimaryColor()
		case b.button.hovered, b.button.tapped:
			return theme.HoverColor()
		default:
			return theme.ButtonColor()
		}
	} else {
		switch b.label.Text {
		case "2":
			return color.RGBA{128, 128, 128, 1}
		case "4":
			return color.RGBA{110, 120, 120, 1}
		case "8":
			return color.RGBA{149, 127, 77, 1}
		case "16":
			return color.RGBA{151, 87, 66, 1}
		case "32":
			return color.RGBA{146, 63, 14, 1}
		case "64":
			return color.RGBA{102, 35, 13, 1}
		case "128":
			return color.RGBA{144, 140, 20, 1}
		case "256":
			return color.RGBA{172, 158, 22, 1}
		case "512":
			return color.RGBA{172, 156, 2, 1}
		case "1024":
			return color.RGBA{178, 141, 53, 1}
		case "2048":
			return color.RGBA{183, 135, 24, 1}
		default:
			return color.RGBA{50, 50, 50, 1}
		}
	}
}

func (b *blockButtonRenderer) Refresh() {
	//b.applyTheme()
	b.label.Text = b.button.Text
	b.label.Color = b.button.Color
	b.label.TextSize = b.button.TextSize

	if b.button.Icon != nil && b.button.Visible() {
		if b.icon == nil {
			b.icon = canvas.NewImageFromResource(b.button.Icon)
			b.icon.FillMode = canvas.ImageFillContain
			b.objects = append(b.objects, b.icon)
		} else {
			if b.button.Disabled() {
				// 如果图标已更改，请创建新的禁用版本
				// if we could be sure that button.Icon is only ever set through the button.SetIcon method, we could remove this
				if !strings.HasSuffix(b.button.disabledIcon.Name(), b.button.Icon.Name()) {
					b.icon.Resource = theme.NewDisabledResource(b.button.Icon)
				} else {
					b.icon.Resource = b.button.disabledIcon
				}
			} else {
				b.icon.Resource = b.button.Icon
			}
		}
		b.icon.Show()
	} else if b.icon != nil {
		b.icon.Hide()
	}

	b.Layout(b.button.Size())
	canvas.Refresh(b.button)
}

func (b *blockButtonRenderer) Objects() []fyne.CanvasObject {
	return b.objects
}

func (b *blockButtonRenderer) Destroy() {
}

// 按钮小部件有一个文本标签，并在单击时触发事件函数
type BlockButton struct {
	widget.DisableableWidget

	Text     string
	TextSize int
	Color    color.Color
	Style    ButtonStyle
	Icon     fyne.Resource

	disabledIcon  fyne.Resource
	Alignment     ButtonAlign
	IconPlacement ButtonIconPlacement

	OnTapped   func() `json:"-"`
	HideShadow bool

	hovered, tapped bool

	ScoreBlock bool
	MenuBlock  bool
}

// 按钮样式决定按钮的行为和呈现
type ButtonStyle int

const (
	// 是标准按钮样式
	DefaultButton ButtonStyle = iota
	// 应该对用户更加突出
	PrimaryButton
)

// 表示按钮的水平对齐方式.
type ButtonAlign int

const (
	// 按钮对齐中心将图标和文本居中对齐.
	ButtonAlignCenter ButtonAlign = iota
	// 按钮导航将图标和文本与前沿对齐.
	ButtonAlignLeading
	// 将图标和文本与后缘对齐.
	ButtonAlignTrailing
)

// 表示按钮中图标和文本的顺序.
type ButtonIconPlacement int

const (
	// 按钮标题文本将图标与文本的前缘对齐.
	ButtonIconLeadingText ButtonIconPlacement = iota
	// 将图标与文本后缘对齐.
	ButtonIconTrailingText
)

// 当捕获指针点击事件并触发任何tap处理程序时，调用Tapped
func (b *BlockButton) Tapped(*fyne.PointEvent) {
	b.tapped = true
	defer func() { // TODO move to a real animation
		time.Sleep(time.Millisecond * buttonTapDuration)
		b.tapped = false
		b.Refresh()
	}()
	b.Refresh()

	if b.OnTapped != nil && !b.Disabled() {
		b.OnTapped()
	}
}

func (b *BlockButton) Tappedo(event *fyne.KeyEvent) {
	b.tapped = true
	defer func() { // TODO move to a real animation
		time.Sleep(time.Millisecond * buttonTapDuration)
		b.tapped = false
		b.Refresh()
	}()
	b.Refresh()

	if b.OnTapped != nil && !b.Disabled() {
		b.OnTapped()
	}
}

// 当桌面指针进入小部件时调用MouseIn
func (b *BlockButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

// 当桌面指针退出小部件时调用MouseOut
func (b *BlockButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

// 当桌面指针悬停在小部件上时调用MouseMoved
func (b *BlockButton) MouseMoved(*desktop.MouseEvent) {
}

// MinSize返回此小部件不应缩小到以下的大小
func (b *BlockButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

// CreateRenderer是Fyne的私有方法，它将这个小部件链接到它的呈现器
func (b *BlockButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	var icon *canvas.Image
	if b.Icon != nil {
		icon = canvas.NewImageFromResource(b.Icon)
		icon.FillMode = canvas.ImageFillContain
	}

	text := canvas.NewText(b.Text, theme.TextColor())
	text.TextStyle.Bold = true
	text.Color = b.Color
	text.TextSize = b.TextSize

	objects := []fyne.CanvasObject{
		text,
	}

	if icon != nil {
		objects = append(objects, icon)
	}

	return &blockButtonRenderer{icon, text, b, objects, layout.NewHBoxLayout()}
}

// 允许更改按钮标签
func (b *BlockButton) SetText(text string) {
	b.Text = text

	b.Refresh()
}

// 返回这个小部件的游标类型
func (b *BlockButton) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// 更新标签上的图标-pass nil以隐藏图标
func (b *BlockButton) SetIcon(icon fyne.Resource) {
	b.Icon = icon

	if icon != nil {
		b.disabledIcon = theme.NewDisabledResource(icon)
	} else {
		b.disabledIcon = nil
	}
	b.Refresh()
}

func (b *BlockButton) SetBlock(num int) {
	text := ""
	size := 0

	t := strconv.Itoa(num)
	if t != "0" {
		text = t
	}
	if num < 10 {
		size = 40
	} else if num >= 10 && num < 100 {
		size = 38
	} else if num >= 100 && num < 1000 {
		size = 36
	} else {
		size = 33
	}

	if num == 2 || num == 4 {
		b.Color = color.Black
	} else {
		b.Color = color.White
	}
	b.Text = text
	b.TextSize = size

}

// 处理程序创建一个新的按钮小部件
func NewBlock(num int, color color.Color, size int, tapped func()) *BlockButton {
	text := ""
	t := strconv.Itoa(num)
	if t != "0" {
		text = t
	}

	button := &BlockButton{
		Text:     text,
		OnTapped: tapped,
		Color:    color,
		TextSize: size,
	}

	button.ExtendBaseWidget(button)
	return button
}

func NewScoreBlock(label string, color color.Color, size int, tapped func()) *BlockButton {
	button := &BlockButton{
		Text:       label,
		OnTapped:   tapped,
		Color:      color,
		TextSize:   size,
		ScoreBlock: true,
	}

	button.ExtendBaseWidget(button)
	return button
}

func NewMenuBlock(label string, color color.Color, size int, tapped func()) *BlockButton {
	button := &BlockButton{
		Text:      label,
		OnTapped:  tapped,
		Color:     color,
		TextSize:  size,
		MenuBlock: true,
	}
	button.ExtendBaseWidget(button)
	return button
}

func NewBlockButtonWithIcon(icon fyne.Resource, tapped func()) *BlockButton {
	button := &BlockButton{
		OnTapped: tapped,
		Icon:     icon,
	}

	button.ExtendBaseWidget(button)
	return button
}

func NewBlockWithlable(label string, color color.Color, size int, tapped func()) *BlockButton {
	button := &BlockButton{
		Text:     label,
		OnTapped: tapped,
		Color:    color,
		TextSize: size,
	}

	button.ExtendBaseWidget(button)
	return button
}
