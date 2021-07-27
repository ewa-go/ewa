package utils

type Nav struct {
	Item     *Item
	SubMenus SubMenus
}

type SubMenus []*Menu

type Item struct {
	Id     string
	Name   string
	Url    string
	Active bool
}

type Menu struct {
	*Item
	SubItems SubItems
}

type SubItems []*Item

type NavBar struct {
	Navs []*Nav
	User string
}

var navbar = NewNavBar()

func NewMenu(id, name, url string, subItem []*Item) *Menu {
	return &Menu{
		Item:     NewItem(id, name, url),
		SubItems: subItem,
	}
}

func NewItem(id, name, url string) *Item {
	return &Item{
		Id:     id,
		Name:   name,
		Url:    url,
		Active: false,
	}
}

func (n NavBar) add(nav *Nav) NavBar {
	n.Navs = append(n.Navs, nav)
	return n
}

func NewNavBar() (n NavBar) {

	n = NavBar{}
	n = n.add(NewNav("section1", "Раздел 1", "", SubMenus{
		NewMenu("section_1_1", "Раздел 1.1", "/section1/1_1", nil),
		NewMenu("section_1_2", "Раздел 1.2", "/section1/1_2", nil),
	})).
		add(NewNav("ucce", "UCCE", "", SubMenus{NewMenu("ucce/rules", "Правила", "/ucce/rules", nil)})).
		add(NewNav("asterisk", "Asterisk", "", SubMenus{NewMenu("asterisk/rules", "Правила", "/ucce/rules", nil)})).
		add(NewNav("callfinder", "Call Finder", "", SubMenus{
			NewMenu("callfinder/calls", "Поиск звонков", "/callfinder/calls", nil),
			NewMenu("callfinder/cucm", "CUCM", "/callfinder/cucm", nil),
			NewMenu("callfinder/cube", "CUBE", "/callfinder/cube", nil),
			NewMenu("callfinder/asterisk", "Asterisk", "/callfinder/asterisk", nil),
		})).
		add(NewNav("settings", "Настройки", "", SubMenus{NewMenu("settings/rules", "Правила", "/settings/rules", nil)})).
		add(NewNav("about", "О программе", "/about", nil))

	return n
}

func NewNav(id, name string, url string, subMenus []*Menu) *Nav {
	return &Nav{
		Item:     NewItem(id, name, url),
		SubMenus: subMenus,
	}
}

func GetNavBar(id string, user string) NavBar {
	navbar.User = user
	for _, nav := range navbar.Navs {
		nav.Item.Active = false
		if nav.Item.Id == id {
			nav.Item.Active = true
		}
		if nav.SubMenus != nil {
			for _, sub := range nav.SubMenus {
				sub.Active = false
				if sub.Id == id {
					sub.Active = true
				}
			}
		}
	}
	return navbar
}
