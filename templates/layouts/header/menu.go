package header

// MenuItem представляет элемент меню
type MenuItem struct {
	URL  string
	Text string
}

// GetDefaultMenuItems возвращает стандартный набор элементов меню
func GetDefaultMenuItems() []MenuItem {
	return []MenuItem{
		{URL: "/", Text: "Главная"},
		{URL: "/about", Text: "О проекте"},
		{URL: "/contact", Text: "Контакты"},
		{URL: "/users", Text: "Пользователи"},
	}
}
