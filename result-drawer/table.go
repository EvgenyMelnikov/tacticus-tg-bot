package result_drawer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strconv"
	tacticus "tacticus-tg-bot/tacticus-wrapper"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/draw"
)

// Основные цвета
var (
	ColorHeader  = color.RGBA{R: 235, G: 225, B: 245, A: 255} // светло-фиолетовый для заголовка
	ColorGray1   = color.RGBA{R: 240, G: 240, B: 240, A: 255} // зебра светлая
	ColorGray2   = color.RGBA{R: 255, G: 255, B: 255, A: 255} // зебра темная
	ColorGreen1  = color.RGBA{R: 220, G: 245, B: 220, A: 255} // готово светло-зеленый
	ColorGreen2  = color.RGBA{R: 210, G: 235, B: 210, A: 255} // готово чуть темнее
	ColorRed1    = color.RGBA{R: 255, G: 200, B: 200, A: 255} // красная зебра светлая
	ColorRed2    = color.RGBA{R: 255, G: 150, B: 150, A: 255} // красная зебра темнее
	ColorYellow1 = color.RGBA{R: 255, G: 255, B: 200, A: 255} // желтая зебра светлая
	ColorYellow2 = color.RGBA{R: 255, G: 255, B: 150, A: 255} // желтая зебра темнее
)

const (
	RowHeight          = 30  // высота строки
	Padding            = 8   // отступ текста от левого края ячейки
	ColNameWidth       = 200 // ширина колонки "Имя"
	BombColTimeWidth   = 150 // ширина колонки "Осталось"
	TokenColCountWidth = 70  // ширина колонки "Токены"
	TokenColTimeWidth  = 150 // ширина колонки "До следующего"
)

func DrawImageWithTables(bombPlayers []tacticus.PlayerData, output string) error {

	gap := 20

	bombTableWidth := ColNameWidth + BombColTimeWidth
	tokenTableWidth := ColNameWidth + TokenColCountWidth + TokenColTimeWidth

	rows := len(bombPlayers)
	height := RowHeight*(rows+1) + 10
	width := bombTableWidth + gap + tokenTableWidth

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// ---- Шрифт ----
	fontBytes, err := os.ReadFile("fonts/Arial.ttf")
	if err != nil {
		return err
	}
	fontParsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fontParsed)
	c.SetFontSize(14)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	// ---- Таблицы ----
	drawBombTable(img, c, bombPlayers, 0, 0)
	drawTokenTable(img, c, bombPlayers, bombTableWidth+gap, 0, rows)

	// ---- Сохранение ----
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	return png.Encode(out, img)
}

func drawBombTable(
	img *image.RGBA,
	c *freetype.Context,
	players []tacticus.PlayerData,
	offsetX, offsetY int,
) {
	width := ColNameWidth + BombColTimeWidth

	sort.Slice(players, func(i, j int) bool {
		return players[i].Progress.GuildRaid.BombTokens.NextTokenInSeconds <
			players[j].Progress.GuildRaid.BombTokens.NextTokenInSeconds
	})

	// отрисовка заголовка
	drawTableHeader(img, c, offsetX, offsetY, width,
		"Бомбы",
		[]string{"Имя", "Осталось"},
		[]int{ColNameWidth, BombColTimeWidth},
	)

	for i, p := range players {
		y1 := offsetY + (i+2)*RowHeight // +2, чтобы учесть две строки заголовка
		y2 := y1 + RowHeight

		seconds := p.Progress.GuildRaid.BombTokens.NextTokenInSeconds

		var bg *image.Uniform
		if seconds <= 0 {
			if i%2 == 0 {
				bg = image.NewUniform(ColorGreen1)
			} else {
				bg = image.NewUniform(ColorGreen2)
			}
		} else if i%2 == 0 {
			bg = image.NewUniform(color.White)
		} else {
			bg = image.NewUniform(ColorGray1)
		}

		draw.Draw(img, image.Rect(offsetX, y1, offsetX+width, y2), bg, image.Point{}, draw.Src)

		duration := time.Duration(seconds) * time.Second
		timeStr := "Готово"
		if duration > 0 {
			timeStr = fmt.Sprintf(
				"%02d:%02d:%02d",
				int(duration.Hours()),
				int(duration.Minutes())%60,
				int(duration.Seconds())%60,
			)
		}

		c.DrawString(p.Details.Name, freetype.Pt(offsetX+Padding, y2-10))
		c.DrawString(timeStr, freetype.Pt(offsetX+ColNameWidth+Padding, y2-10))
	}

	// рисуем сетку
	drawHorizontalLines(img, offsetX, offsetY, ColNameWidth+BombColTimeWidth, len(players), RowHeight, 2)
	drawVerticalLines(img, offsetX, offsetY, []int{ColNameWidth, BombColTimeWidth}, len(players), RowHeight, 2)
}

func drawTokenTable(img *image.RGBA, c *freetype.Context, players []tacticus.PlayerData, offsetX, offsetY int, totalRows int) {
	totalWidth := ColNameWidth + TokenColCountWidth + TokenColTimeWidth

	// заголовок
	drawTableHeader(img, c, offsetX, offsetY, totalWidth,
		"Токены",
		[]string{"Имя", "Токены", "До следующего"},
		[]int{ColNameWidth, TokenColCountWidth, TokenColTimeWidth},
	)

	sort.Slice(players, func(i, j int) bool {
		tokenInfoI := players[i].Progress.GuildRaid.Tokens
		tokenInfoJ := players[j].Progress.GuildRaid.Tokens

		totalDurationI := (tokenInfoI.Max-tokenInfoI.Current)*tokenInfoI.RegenDelayInSeconds + tokenInfoI.NextTokenInSeconds
		totalDurationG := (tokenInfoJ.Max-tokenInfoJ.Current)*tokenInfoJ.RegenDelayInSeconds + tokenInfoJ.NextTokenInSeconds

		return totalDurationI < totalDurationG
	})

	for i, p := range players {
		y1 := offsetY + (i+2)*RowHeight
		y2 := y1 + RowHeight

		bg := image.NewUniform(color.White)

		if p.Progress.GuildRaid.Tokens.Current == 3 {
			if i%2 == 0 {
				bg = image.NewUniform(ColorRed1)
			} else {
				bg = image.NewUniform(ColorRed2)
			}
		} else if p.Progress.GuildRaid.Tokens.Current == 2 {
			if i%2 == 0 {
				bg = image.NewUniform(ColorYellow1)
			} else {
				bg = image.NewUniform(ColorYellow2)
			}
		} else if p.Progress.GuildRaid.Tokens.Current == 1 {
			if i%2 == 0 {
				bg = image.NewUniform(ColorGreen1)
			} else {
				bg = image.NewUniform(ColorGreen2)
			}
		} else if i%2 == 0 {
			bg = image.NewUniform(color.White)
		} else {
			bg = image.NewUniform(ColorGray1)
		}

		draw.Draw(img, image.Rect(offsetX, y1, offsetX+totalWidth, y2), bg, image.Point{}, draw.Src)

		if i >= len(players) {
			continue
		}

		seconds := p.Progress.GuildRaid.Tokens.NextTokenInSeconds
		duration := time.Duration(seconds) * time.Second

		timeStr := "Готово"
		if duration > 0 {
			timeStr = fmt.Sprintf("%02d:%02d:%02d",
				int(duration.Hours()),
				int(duration.Minutes())%60,
				int(duration.Seconds())%60,
			)
		}

		x := offsetX + 8
		c.DrawString(p.Details.Name, freetype.Pt(x, y2-10))
		x += ColNameWidth
		c.DrawString(strconv.Itoa(p.Progress.GuildRaid.Tokens.Current), freetype.Pt(x, y2-10))
		x += TokenColCountWidth
		c.DrawString(timeStr, freetype.Pt(x, y2-10))
	}

	// рисуем сетку
	drawHorizontalLines(img, offsetX, offsetY, ColNameWidth+TokenColCountWidth+TokenColTimeWidth, len(players), RowHeight, 2)
	drawVerticalLines(img, offsetX, offsetY, []int{ColNameWidth, TokenColCountWidth, TokenColTimeWidth}, len(players), RowHeight, 2)
}

func drawTableHeader(
	img *image.RGBA,
	c *freetype.Context,
	offsetX, offsetY, tableWidth int,
	title string,
	columnNames []string,
	colWidths []int,
) {
	// фон заголовка (2 ряда: название + колонки)
	draw.Draw(img, image.Rect(offsetX, offsetY, offsetX+tableWidth, offsetY+RowHeight*2), image.NewUniform(ColorHeader), image.Point{}, draw.Src)

	// центрируем название таблицы
	textWidth := len([]rune(title)) * 8 // примерное измерение
	centerX := offsetX + (tableWidth-textWidth)/2
	c.DrawString(title, freetype.Pt(centerX, offsetY+RowHeight-10))

	// названия колонок
	x := offsetX + 8
	for i, name := range columnNames {
		c.DrawString(name, freetype.Pt(x, offsetY+RowHeight*2-10))
		x += colWidths[i]
	}
}

func drawHorizontalLines(img *image.RGBA, offsetX, offsetY, width, rows, rowHeight int, headerRows int) {
	for i := 0; i <= rows+headerRows; i++ {
		y := offsetY + i*rowHeight
		for x := offsetX; x < offsetX+width; x++ {
			img.Set(x, y, color.Black)
		}
	}
}

func drawVerticalLines(img *image.RGBA, offsetX, offsetY int, colWidths []int, rows, rowHeight, headerRows int) {
	for i := 0; i <= len(colWidths); i++ {
		x := offsetX
		for j := 0; j < i; j++ {
			x += colWidths[j]
		}

		startY := offsetY
		if i != 0 && i != len(colWidths) {
			startY += rowHeight
		}
		for y := startY; y <= offsetY+(rows+headerRows)*rowHeight; y++ {
			img.Set(x, y, color.Black)
		}
	}
}
