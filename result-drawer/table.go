package result_drawer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	tacticus "tacticus-tg-bot/tacticus-wrapper"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/draw"
)

// Генерация картинки таблицы
func DrawTable(players []tacticus.PlayerData, output string) error {
	// Размер изображения
	width := 500
	height := 30 + len(players)*25 // заголовок + строки

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// фон белый
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Настройка цвета текста
	black := image.NewUniform(color.Black)

	// Подгружаем шрифт
	fontBytes, err := os.ReadFile("fonts/Arial.ttf") // или любой .ttf
	if err != nil {
		return err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(14)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(black)

	pt := freetype.Pt(10, 20) // начальная позиция текста

	// Заголовок
	_, _ = c.DrawString("Имя", pt)
	_, _ = c.DrawString("Осталось", freetype.Pt(250, 20))

	// Игроки
	y := 45
	for _, p := range players {
		seconds := p.Progress.GuildRaid.BombTokens.NextTokenInSeconds
		duration := time.Duration(seconds) * time.Second
		var timeStr string
		if duration <= 0 {
			timeStr = "Готова"
		} else {
			timeStr = fmt.Sprintf("%02d:%02d:%02d",
				int(duration.Hours()),
				int(duration.Minutes())%60,
				int(duration.Seconds())%60)
		}

		_, _ = c.DrawString(p.Details.Name, freetype.Pt(10, y))
		_, _ = c.DrawString(timeStr, freetype.Pt(250, y))
		y += 25
	}

	// Сохраняем PNG
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return png.Encode(outFile, img)
}
