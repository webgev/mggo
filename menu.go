package mggo

type menuItem struct {
    ID, Title, Href, ParentID string
    Active                    bool
}

// Menu is items menu
type Menu struct {
    items []menuItem
}

// SetActivePage is set active flag by menu id
func (m *Menu) SetActivePage(id string) {
    for i := 0; i < len(m.items); i++ {
        m.items[i].Active = m.items[i].ID == id
    }
}

// Append is append menu item
func (m *Menu) Append(id, title, href string) {
    if m.items == nil {
        m.items = []menuItem{}
    }
    m.items = append(m.items, menuItem{ID: id, Title: title, Href: href})
}
