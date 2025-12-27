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

func DrawImageWithTables(
	bombPlayers []tacticus.PlayerData,
	output string,
) error {

	rowHeight := 30
	gap := 20

	bombTableWidth := 420
	tokenTableWidth := 600

	rows := len(bombPlayers)
	height := rowHeight*(rows+1) + 10
	width := bombTableWidth + gap + tokenTableWidth

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

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
	defer out.Close()

	return png.Encode(out, img)
}

func drawBombTable(
	img *image.RGBA,
	c *freetype.Context,
	players []tacticus.PlayerData,
	offsetX, offsetY int,
) {
	rowHeight := 30
	colName := 260
	colTime := 160
	padding := 8
	width := colName + colTime

	grayBg := image.NewUniform(color.RGBA{240, 240, 240, 255})
	whiteBg := image.NewUniform(color.White)
	green1 := image.NewUniform(color.RGBA{220, 245, 220, 255})
	green2 := image.NewUniform(color.RGBA{210, 235, 210, 255})

	sort.Slice(players, func(i, j int) bool {
		return players[i].Progress.GuildRaid.BombTokens.NextTokenInSeconds <
			players[j].Progress.GuildRaid.BombTokens.NextTokenInSeconds
	})

	// отрисовка заголовка
	drawTableHeader(img, c, offsetX, offsetY, width,
		"💣 Bomb",
		[]string{"Имя", "Осталось"},
		[]int{colName, colTime},
	)

	for i, p := range players {
		y1 := offsetY + (i+2)*rowHeight // +2, чтобы учесть две строки заголовка
		y2 := y1 + rowHeight

		seconds := p.Progress.GuildRaid.BombTokens.NextTokenInSeconds

		var bg *image.Uniform
		if seconds <= 0 {
			if i%2 == 0 {
				bg = green1
			} else {
				bg = green2
			}
		} else if i%2 == 0 {
			bg = grayBg
		} else {
			bg = whiteBg
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

		c.DrawString(p.Details.Name, freetype.Pt(offsetX+padding, y2-10))
		c.DrawString(timeStr, freetype.Pt(offsetX+colName+padding, y2-10))
	}
}

func drawTokenTable(img *image.RGBA, c *freetype.Context, players []tacticus.PlayerData, offsetX, offsetY int, totalRows int) {
	rowHeight := 30
	colNameWidth := 240
	colCountWidth := 140
	colTimeWidth := 180
	totalWidth := colNameWidth + colCountWidth + colTimeWidth

	// заголовок
	drawTableHeader(img, c, offsetX, offsetY, totalWidth,
		"🪙 Token",
		[]string{"Имя", "Токены", "До следующего"},
		[]int{colNameWidth, colCountWidth, colTimeWidth},
	)

	whiteBg := image.NewUniform(color.White)
	grayBg := image.NewUniform(color.RGBA{245, 245, 245, 255})
	red1 := image.NewUniform(color.RGBA{255, 200, 200, 255})
	red2 := image.NewUniform(color.RGBA{255, 150, 150, 255})
	yellow1 := image.NewUniform(color.RGBA{255, 255, 200, 255})
	yellow2 := image.NewUniform(color.RGBA{255, 255, 150, 255})
	green1 := image.NewUniform(color.RGBA{220, 245, 220, 255})
	green2 := image.NewUniform(color.RGBA{210, 235, 210, 255})

	sort.Slice(players, func(i, j int) bool {
		tokenInfoI := players[i].Progress.GuildRaid.Tokens
		tokenInfoJ := players[j].Progress.GuildRaid.Tokens

		totalDurationI := (tokenInfoI.Max-tokenInfoI.Current)*tokenInfoI.RegenDelayInSeconds + tokenInfoI.NextTokenInSeconds
		totalDurationG := (tokenInfoJ.Max-tokenInfoJ.Current)*tokenInfoJ.RegenDelayInSeconds + tokenInfoJ.NextTokenInSeconds

		return totalDurationI < totalDurationG
	})

	for i, p := range players {
		y1 := offsetY + (i+2)*rowHeight
		y2 := y1 + rowHeight

		bg := whiteBg
		if i%2 == 0 {
			bg = grayBg
		}

		if p.Progress.GuildRaid.Tokens.Current == 3 {
			if i%2 == 0 {
				bg = red1
			} else {
				bg = red2
			}
		} else if p.Progress.GuildRaid.Tokens.Current == 2 {
			if i%2 == 0 {
				bg = yellow1
			} else {
				bg = yellow2
			}
		} else if p.Progress.GuildRaid.Tokens.Current == 1 {
			if i%2 == 0 {
				bg = green1
			} else {
				bg = green2
			}
		} else if i%2 == 0 {
			bg = grayBg
		} else {
			bg = whiteBg
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
		x += colNameWidth
		c.DrawString(strconv.Itoa(p.Progress.GuildRaid.Tokens.Current), freetype.Pt(x, y2-10))
		x += colCountWidth
		c.DrawString(timeStr, freetype.Pt(x, y2-10))
	}
}

func drawTableHeader(
	img *image.RGBA,
	c *freetype.Context,
	offsetX, offsetY, tableWidth int,
	title string,
	columnNames []string,
	colWidths []int,
) {
	rowHeight := 30
	headerBg := image.NewUniform(color.RGBA{235, 225, 245, 255})

	// фон заголовка (2 ряда: название + колонки)
	draw.Draw(img, image.Rect(offsetX, offsetY, offsetX+tableWidth, offsetY+rowHeight*2), headerBg, image.Point{}, draw.Src)

	// центрируем название таблицы
	textWidth := len([]rune(title)) * 8 // примерное измерение
	centerX := offsetX + (tableWidth-textWidth)/2
	c.DrawString(title, freetype.Pt(centerX, offsetY+rowHeight-10))

	// названия колонок
	x := offsetX + 8
	for i, name := range columnNames {
		c.DrawString(name, freetype.Pt(x, offsetY+rowHeight*2-10))
		x += colWidths[i]
	}
}
