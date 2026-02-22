package add

type action int

const (
	actionDownload action = iota
	actionOnline
	actionAdd
)

var actionStringValues = []string{
	"Скачать",
	"Смотреть онлайн",
	"Добавить в список",
}

var actionSuccessMessages = []string{
	"Скачивание началось",
	"Готово к просмотру онлайн",
	"Добавлено в список возможного просмотра",
}

func (a action) String() string {
	return actionStringValues[a]
}

func actionFromString(s string) (action, bool) {
	for a := range actionStringValues {
		if actionStringValues[a] == s {
			return action(a), true
		}
	}

	return actionDownload, false
}

func (a action) SuccessMessage() string {
	return actionSuccessMessages[a]
}
