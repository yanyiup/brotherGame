package view

import (
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var (
	Width     = 4  // 横向格子数量
	Height    = 4  // 纵向格子数量
	ScoreSize = 20 // 得分数字大小
	NumSize   = 40 // 数字初始大小
)

var SCORE int
var HIGHESTSCORE int

// 游戏模块
type GameModule struct {
	Window       fyne.Window
	GameCanvas   fyne.CanvasObject
	BlockCanvas  *fyne.Container
	Start        *BlockButton
	Afresh       *BlockButton
	HighestScore *BlockButton
	Score        *BlockButton
	Buttons      map[string]*widget.Button
	Blocks       [][] *BlockButton
	CurrentNum   [][] int
	OldNum       [][] int
	Win          bool
	Lose         bool
}

func (g *GameModule) LoadGame(win fyne.Window) {
	g.Window = win
	g.Buttons = make(map[string]*widget.Button, 0)
	g.Blocks = make([][]*BlockButton, 0)
	g.CurrentNum = make([][]int, 0)

	for i := 0; i < Height; i++ {
		block := make([]*BlockButton, 0)
		current := make([]int, 0)
		for j := 0; j < Width; j++ {
			b := NewBlock(0, color.White, NumSize, nil)
			block = append(block, b)
			current = append(current, 0)
		}
		g.Blocks = append(g.Blocks, block)
		g.CurrentNum = append(g.CurrentNum, current)
	}

	g.BlockCanvas = fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	for i := 0; i < Height; i++ {
		h := fyne.NewContainerWithLayout(layout.NewGridLayout(Width))
		for j := 0; j < Width; j++ {
			h.AddObject(g.Blocks[i][j])
		}
		g.BlockCanvas.AddObject(h)
	}

	// 随机生成数字
	g.Random(2)

	// 测试数据
	// test(g)

	top := widget.NewVBox(
		layout.NewSpacer(),
		layout.NewSpacer(),
		makeCell(10, 10),
		widget.NewLabelWithStyle("Welcome to the digital world", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		makeCell(10, 10),
		layout.NewSpacer(),
	)

	padding := makeCell(30, 30)
	borderLayout := layout.NewBorderLayout(top, padding, padding, padding)
	board01 := fyne.NewContainerWithLayout(borderLayout, top, padding, padding, padding, g.BlockCanvas)

	// 得分展示相关
	g.Score = NewScoreBlock("0", color.Black, ScoreSize, nil)
	g.HighestScore = NewScoreBlock("0", color.Black, ScoreSize, nil)

	g.Start = NewMenuBlock("新游戏", color.White, ScoreSize, g.Restart)
	g.Afresh = NewMenuBlock("回退", color.White, ScoreSize, g.Back)

	title := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		NewBlockWithlable("当前得分", color.White, ScoreSize, nil),
		NewBlockWithlable("最高得分", color.White, ScoreSize, nil),
	)

	score := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		g.Score,
		g.HighestScore,
	)
	start := fyne.NewContainerWithLayout(layout.NewGridLayout(3),
		g.Start,
		widget.NewLabel(""),
		g.Afresh,
	)

	control := fyne.NewContainerWithLayout(layout.NewGridLayout(1), title, score, makeCell(10, 10), start)

	bor := layout.NewBorderLayout(padding, padding, padding, padding)
	operate01 := fyne.NewContainerWithLayout(bor, padding, padding, padding, padding, control)

	// 方向键相关
	up := NewBlockButtonWithIcon(theme.MoveUpIcon(), g.Up)
	down := NewBlockButtonWithIcon(theme.MoveDownIcon(), g.Down)
	left := NewBlockButtonWithIcon(theme.NavigateBackIcon(), g.Left)
	right := NewBlockButtonWithIcon(theme.NavigateNextIcon(), g.Right)

	direction := fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			widget.NewLabel(""),
			up,
			widget.NewLabel(""),
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			left,
			down,
			right,
		),
	)

	operate02 := fyne.NewContainerWithLayout(bor, padding, padding, padding, padding, direction)

	// 自定义平分模块
	board02 := &widget.SplitContainer{
		Offset:     0.6,
		Horizontal: false,
		Leading:    operate01,
		Trailing:   operate02,
	}
	board02.ExtendBaseWidget(board02)

	g.Keyboard()
	g.GameCanvas = widget.NewHSplitContainer(board01, board02)

}

// 键盘输入
func (g *GameModule) Keyboard() {
	g.Window.Canvas().SetOnTypedRune(g.typedRune)
	g.Window.Canvas().SetOnTypedKey(g.typedKey)
}

func (g *GameModule) typedRune(r rune) {

}

func (g *GameModule) typedKey(ev *fyne.KeyEvent) {
	if !g.Win || !g.Lose {
		switch ev.Name {
		case fyne.KeyUp:
			g.Up()
		case fyne.KeyDown:
			g.Down()
		case fyne.KeyLeft:
			g.Left()
		case fyne.KeyRight:
			g.Right()
		case fyne.KeySpace:
			g.Back()
		case fyne.KeyF12:
			test(g)
		}
	}
}

// 判断输赢
func (g *GameModule) WinOrLose() {
	space := false
	win := false
	alike := false

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			if g.CurrentNum [i][j] == 0 {
				space = true
			}
			if g.CurrentNum [i][j] >= 2048 {
				win = true
			}
		}
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			if j+1 < Width {
				if g.CurrentNum [i][j] == g.CurrentNum [i][j+1] {
					alike = true
				}
			}
		}
	}

	for j := 0; j < Width; j++ {
		for i := 0; i < Height; i++ {
			if i+1 < Height {
				if g.CurrentNum [i][j] == g.CurrentNum [i+1][j] {
					alike = true
				}
			}
		}
	}

	if win {
		cnf := dialog.NewConfirm("(づ｡◕ᴗᴗ◕｡)づ", "恭喜你赢得了全世界！", g.Callback, g.Window)
		cnf.SetDismissText("回味一下")
		cnf.SetConfirmText("再来一次")
		cnf.Show()
		g.Win = true
	}

	if !space && !alike {
		dialog.ShowInformation("ε(┬┬﹏┬┬)3", "你输了！再接再厉", g.Window)
		g.Lose = true
	}

}

func (g *GameModule) Callback(ok bool) {
	if ok {
		g.Restart()
	}
}

// 从新开始
func (g *GameModule) Restart() {
	g.Win = false
	g.Lose = false
	g.OldNum = make([][]int, 0)
	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.CurrentNum [i][j] = 0
		}
	}
	g.Random(2)
	if HIGHESTSCORE < SCORE {
		HIGHESTSCORE = SCORE
		g.HighestScore.Text = strconv.Itoa(HIGHESTSCORE)
		g.HighestScore.Refresh()
	}
	SCORE = 0
	g.Score.Text = strconv.Itoa(SCORE)
	g.Score.Refresh()
}

// 回退
func (g *GameModule) Back() {
	if g.Win || g.Lose {
		return
	}
	if len(g.OldNum) > 0 {
		g.CurrentNum = g.OldNum
		g.OldNum = make([][]int, 0)
		for i := 0; i < len(g.CurrentNum); i++ {
			for j := 0; j < len(g.CurrentNum[i]); j++ {
				g.Blocks [i][j].SetBlock(g.CurrentNum[i][j])
			}
		}
		g.BlockCanvas.Refresh()
	} else {
		dialog.ShowInformation("<(▰˘◡˘▰)>", "当前状态下无法回退", g.Window)
	}
}

// 生成随机数
func (g *GameModule) Random(num int) {
	type coordinate struct {
		x   int
		y   int
		num int
	}

	s := make([]*coordinate, 0)

	for i := 0; i < len(g.CurrentNum); i++ {
		for j := 0; j < len(g.CurrentNum[i]); j++ {
			if g.CurrentNum[i][j] == 0 {
				coord := coordinate{
					x: j,
					y: i,
				}
				s = append(s, &coord)
			}

		}
	}

	if len(s) <= 0 {
		return
	}

	insert := make([]*coordinate, 0)
	rand.Seed(time.Now().UnixNano())
	if num == 2 {
		for {
			if len(insert) >= 2 {
				break
			}
			i := rand.Intn(len(s))
			s[i].num = 2
			insert = append(insert, s[i])
			s = append(s[0:i], s[i+1:]...)
		}

	} else {
		nums := []int{2, 2, 4}
		i := rand.Intn(len(s))
		n := rand.Intn(len(nums))
		s[i].num = nums[n]
		insert = append(insert, s[i])
	}

	for i := 0; i < len(insert); i++ {
		g.CurrentNum[insert[i].y][insert[i].x] = insert[i].num
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.Blocks [i][j].SetBlock(g.CurrentNum[i][j])
		}
	}
	g.BlockCanvas.Refresh()

}

func (g *GameModule) Up() {
	if g.Win || g.Lose {
		return
	}
	g.WinOrLose()
	if len(g.OldNum) == 0 {
		g.OldNum = make([][]int, Height)
		for i := range g.OldNum {
			g.OldNum[i] = make([]int, Width)
		}
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.OldNum[i][j] = g.CurrentNum[i][j]
		}
	}

	up(g)

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.Blocks [i][j].SetBlock(g.CurrentNum[i][j])
		}
	}

	g.Random(1)

	g.Score.Text = strconv.Itoa(SCORE)
	g.Score.Refresh()
	g.HighestScore.Text = strconv.Itoa(HIGHESTSCORE)
	g.HighestScore.Refresh()
	g.BlockCanvas.Refresh()
	g.WinOrLose()
}

func (g *GameModule) Down() {
	if g.Win || g.Lose {
		return
	}
	if len(g.OldNum) == 0 {
		g.OldNum = make([][]int, Height)
		for i := range g.OldNum {
			g.OldNum[i] = make([]int, Width)
		}
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.OldNum[i][j] = g.CurrentNum[i][j]
		}
	}

	right90(g)
	right90(g)
	up(g)
	left90(g)
	left90(g)

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.Blocks [i][j].SetBlock(g.CurrentNum[i][j])
		}
	}

	g.Random(1)
	g.Score.Text = strconv.Itoa(SCORE)
	g.Score.Refresh()
	g.HighestScore.Text = strconv.Itoa(HIGHESTSCORE)
	g.HighestScore.Refresh()
	g.BlockCanvas.Refresh()
	g.WinOrLose()
}

func (g *GameModule) Left() {
	if g.Win || g.Lose {
		return
	}
	if len(g.OldNum) == 0 {
		g.OldNum = make([][]int, Height)
		for i := range g.OldNum {
			g.OldNum[i] = make([]int, Width)
		}
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.OldNum[i][j] = g.CurrentNum[i][j]
		}
	}

	right90(g)
	up(g)
	left90(g)

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.Blocks [i][j].SetBlock(g.CurrentNum[i][j])
		}
	}

	g.Random(1)
	g.Score.Text = strconv.Itoa(SCORE)
	g.Score.Refresh()
	g.HighestScore.Text = strconv.Itoa(HIGHESTSCORE)
	g.HighestScore.Refresh()
	g.BlockCanvas.Refresh()
	g.WinOrLose()

}

func (g *GameModule) Right() {
	if g.Win || g.Lose {
		return
	}
	if len(g.OldNum) == 0 {
		g.OldNum = make([][]int, Height)
		for i := range g.OldNum {
			g.OldNum[i] = make([]int, Width)
		}
	}

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.OldNum[i][j] = g.CurrentNum[i][j]
		}
	}

	left90(g)
	up(g)
	right90(g)

	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			g.Blocks[i][j].SetBlock(g.CurrentNum[i][j])
		}
	}

	g.Random(1)
	g.Score.Text = strconv.Itoa(SCORE)
	g.Score.Refresh()
	g.HighestScore.Text = strconv.Itoa(HIGHESTSCORE)
	g.HighestScore.Refresh()
	g.BlockCanvas.Refresh()
	g.WinOrLose()
}

func left90(g *GameModule) {
	temp := make([][]int, Height)

	for i := range temp {
		temp[i] = make([]int, Width)
	}
	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			temp[Height-j-1][i] = g.CurrentNum[i][j]
		}
	}

	g.CurrentNum = temp
	return

}

func right90(g *GameModule) {
	temp := make([][]int, Height)

	for i := range temp {
		temp[i] = make([]int, Width)
	}
	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			temp[j][Height-i-1] = g.CurrentNum[i][j]
		}
	}
	g.CurrentNum = temp
}

func up(g *GameModule) {
	for j := 0; j < Width; j++ {
		point := 0
		for point < Height-1 {
			find(g.CurrentNum, j, &point)
		}
	}

}

func find(currentNum [][]int, j int, point *int) {
	start := currentNum[*point][j]
	temp := false

	for i := *point; i < Height; i++ {
		if currentNum[i][j] != 0 {
			if start == 0 {
				currentNum[*point][j] = currentNum[i][j]
				currentNum[i][j] = 0
				temp = true
				break
			} else {
				if *point != i {
					if currentNum[*point][j] == currentNum[i][j] {
						currentNum[*point][j] += currentNum[i][j]
						//记录得分
						SCORE += 2 * currentNum[i][j]
						currentNum[i][j] = 0
						temp = true
						break
					} else {
						if *point+1 != i {
							currentNum[*point+1][j] = currentNum[i][j]
							currentNum[i][j] = 0
							temp = true
							break
						} else {
							temp = true
							break
						}
					}
				}
			}
		}

	}

	if !temp {
		*point = Height - 1
		return
	} else if start == 0 {
		return
	} else {
		*point++
		return
	}

}

// 制造填充区域
func makeCell(w, h int) fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.NRGBA{0x42, 0x42, 0x42, 0xff})
	rect.SetMinSize(fyne.NewSize(w, h))
	return rect
}

// 测试使用
func test(g *GameModule) {
	g.Blocks[0][0].Text = "2"
	g.CurrentNum[0][0] = 2
	g.Blocks[0][0].Color = color.Black
	g.Blocks[0][1].Text = "4"
	g.CurrentNum[0][1] = 4
	g.Blocks[0][1].Color = color.Black
	g.Blocks[0][2].Text = "8"
	g.CurrentNum[0][2] = 8
	g.Blocks[0][2].Color = color.White
	g.Blocks[0][3].Text = "16"
	g.CurrentNum[0][3] = 16
	g.Blocks[0][3].Color = color.White
	g.Blocks[0][3].TextSize = 38
	g.Blocks[1][0].Text = "32"
	g.CurrentNum[1][0] = 32
	g.Blocks[1][0].Color = color.White
	g.Blocks[1][0].TextSize = 38
	g.Blocks[1][1].Text = "64"
	g.CurrentNum[1][1] = 64
	g.Blocks[1][1].TextSize = 38
	g.Blocks[1][2].Text = "128"
	g.CurrentNum[1][2] = 128
	g.Blocks[1][2].TextSize = 36
	g.Blocks[1][3].Text = "256"
	g.CurrentNum[1][3] = 256
	g.Blocks[1][3].TextSize = 36
	g.Blocks[2][0].Text = "512"
	g.CurrentNum[2][0] = 512
	g.Blocks[2][0].TextSize = 36
	g.Blocks[2][1].Text = "1024"
	g.CurrentNum[2][1] = 1024
	g.Blocks[2][1].TextSize = 33
	g.Blocks[2][2].Text = "2048"
	g.CurrentNum[2][2] = 2048
	g.Blocks[2][2].TextSize = 33
	g.Blocks[2][3].Text = ""
	g.CurrentNum[2][3] = 0
	g.Blocks[3][0].Text = ""
	g.CurrentNum[3][0] = 0
	g.Blocks[3][1].Text = ""
	g.CurrentNum[3][1] = 0
	g.Blocks[3][2].Text = ""
	g.CurrentNum[3][2] = 0
	g.Blocks[3][3].Text = ""
	g.CurrentNum[3][3] = 0
	g.BlockCanvas.Refresh()
}
