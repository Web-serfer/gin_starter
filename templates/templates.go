package templates

import (
	pages "gin-starter/templates/pages"
	layouts "gin-starter/templates/layouts"
	header "gin-starter/templates/layouts/header"
	footer "gin-starter/templates/layouts/footer"
	"github.com/a-h/templ"
)

// GetDefaultMenuItems возвращает стандартный набор элементов меню
func GetDefaultMenuItems() []header.MenuItem {
	return header.GetDefaultMenuItems()
}

// Обертки для шаблонов страниц
func IndexPage(canonicalURL string, menuItems []header.MenuItem) templ.Component {
	return pages.IndexPage(canonicalURL, menuItems)
}

func AboutPage(canonicalURL string, menuItems []header.MenuItem) templ.Component {
	return pages.AboutPage(canonicalURL, menuItems)
}

func ContactPage(canonicalURL string, menuItems []header.MenuItem) templ.Component {
	return pages.ContactPage(canonicalURL, menuItems)
}

func UsersPage(canonicalURL string, menuItems []header.MenuItem) templ.Component {
	return pages.UsersPage(canonicalURL, menuItems)
}

// Обертки для шаблонов макетов
func Layout(title string, description string, canonicalURL string, menuItems []header.MenuItem, body templ.Component) templ.Component {
	return layouts.Layout(title, description, canonicalURL, menuItems, body)
}

func Header(menuItems []header.MenuItem) templ.Component {
	return header.Header(menuItems)
}

func Footer(menuItems []header.MenuItem, currentYear int) templ.Component {
	return footer.Footer(menuItems, currentYear)
}

func Logo() templ.Component {
	return header.Logo()
}