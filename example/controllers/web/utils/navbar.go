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

type Navbar struct {
	Navs []*Nav
	User string
}

func newMenu(id, name, url string, subItem []*Item) *Menu {
	return &Menu{
		Item:     newItem(id, name, url),
		SubItems: subItem,
	}
}

func newItem(id, name, url string) *Item {
	return &Item{
		Id:     id,
		Name:   name,
		Url:    url,
		Active: false,
	}
}

func (n Navbar) add(nav *Nav) Navbar {
	n.Navs = append(n.Navs, nav)
	return n
}

func createNavbar() Navbar {

	n := Navbar{}
	n = n.add(newNav("section1", "Раздел 1", "", SubMenus{
		newMenu("section_1_1", "Раздел 1.1", "", SubItems{
			newItem("section1/1_1/document", "Документ", "/section1/1_1/document"),
			newItem("section1/1_1/list", "Список", "/section1/1_1/list"),
		}),
		newMenu("section_1_2", "Раздел 1.2", "/section1/1_2", nil),
	})).
		add(newNav("ucce", "UCCE", "", SubMenus{newMenu("ucce/rules", "Правила", "/ucce/rules", nil)})).
		add(newNav("asterisk", "Asterisk", "", SubMenus{newMenu("asterisk/rules", "Правила", "/ucce/rules", nil)})).
		add(newNav("callfinder", "Call Finder", "", SubMenus{
			newMenu("callfinder/calls", "Поиск звонков", "/callfinder/calls", nil),
			newMenu("callfinder/cucm", "CUCM", "/callfinder/cucm", nil),
			newMenu("callfinder/cube", "CUBE", "/callfinder/cube", nil),
			newMenu("callfinder/asterisk", "Asterisk", "/callfinder/asterisk", nil),
		})).
		add(newNav("settings", "Настройки", "", SubMenus{newMenu("settings/rules", "Правила", "/settings/rules", nil)})).
		add(newNav("about", "О программе", "/about", nil))

	return n
}

func newNav(id, name string, url string, subMenus []*Menu) *Nav {
	return &Nav{
		Item:     newItem(id, name, url),
		SubMenus: subMenus,
	}
}

func NewNavbar(id string, user string) (n Navbar) {
	n = createNavbar()
	n.User = user
	for _, nav := range n.Navs {
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
	return n
}
