package search

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"go-micro.dev/v4/logger"
	"strings"
	"text/template"
	"unicode/utf8"
)

//go:embed templates
var templates embed.FS

var parsedTemplates *template.Template

func init() {
	parsedTemplates = template.Must(template.ParseFS(templates, "templates/*.txt"))
}

func formatGenres(genres []string) string {
	result := ""
	for _, g := range genres {
		result += strings.TrimSpace(strings.ToLower(g)) + ", "
	}
	if len(result) > 2 {
		result = result[0 : len(result)-2]
	}
	return result
}

func formatDescription(d string) string {
	const maxLength = 350
	if utf8.RuneCountInString(d) <= maxLength {
		return d
	}

	cnt := 0
	found := false
	split := strings.FieldsFunc(d, func(r rune) bool {
		cnt++
		if cnt > maxLength && r == ' ' && !found {
			found = true
			return true
		}
		return false
	})
	return split[0] + "..."
}

func (s *searchCommand) formatMovieMessage(mov *rms_library.FoundMovie) *communication.BotMessage {
	m := &communication.BotMessage{}
	if mov.Info.Poster != "" {
		m.Attachment = &communication.Attachment{
			Type:     communication.Attachment_PhotoURL,
			MimeType: "",
			Content:  []byte(mov.Info.Poster),
		}
	}

	m.Buttons = append(m.Buttons, &communication.Button{Title: "Скачать", Command: "/download auto " + mov.Id})
	m.Buttons = append(m.Buttons, &communication.Button{Title: "Скачать быстрее", Command: "/download faster " + mov.Id})
	m.Buttons = append(m.Buttons, &communication.Button{Title: "Выбрать раздачу", Command: "/download select " + mov.Id})
	m.Buttons = append(m.Buttons, &communication.Button{Title: "Файл", Command: "/download file " + mov.Id})

	m.KeyboardStyle = communication.KeyboardStyle_Message

	var ui struct {
		Title       string
		Year        uint32
		Rating      string
		Genres      string
		Description string
	}
	ui.Title = mov.Info.Title
	ui.Year = mov.Info.Year
	ui.Rating = fmt.Sprintf("%.1f", mov.Info.Rating)
	ui.Genres = formatGenres(mov.Info.Genres)
	ui.Description = formatDescription(mov.Info.Description)

	var buf bytes.Buffer
	if err := parsedTemplates.ExecuteTemplate(&buf, "movie", &ui); err != nil {
		s.l.Logf(logger.ErrorLevel, "execute template failed: %s", err)
	}
	m.Text = buf.String()
	return m
}
