package utils

type BreadCrumb struct {
	Name string
	Url  string
	Page bool
}

type BreadCrumbs []BreadCrumb

func NewBreadCrumb(name string, page ...bool) (b BreadCrumb) {
	b = BreadCrumb{
		Name: name,
	}
	if len(page) > 0 {
		b.Page = page[0]
	}
	return
}

func DefaultBreadCrumbs(names ...string) (b BreadCrumbs) {
	for i, name := range names {
		page := false
		if i == len(names)-1 {
			page = true
		}
		b = append(b, NewBreadCrumb(name, page))
	}
	return
}

func NewBreadCrumbs(breadCrumb ...BreadCrumb) (b BreadCrumbs) {
	b = append(b, breadCrumb...)
	return
}
