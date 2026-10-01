// Quillnote: a small notes app used as the proof fixture of the Fyne mockup recipe.
// All data is invented and fixed, so every build looks the same.
package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type note struct {
	Title, Body, Tag string
	Pinned           bool
}

var notes = []note{
	{"Harbour trip checklist", "Tide table printed\nSpare mooring line\nThermos and two mugs\nCharge the handheld VHF", "Travel", true},
	{"Sourdough schedule", "Feed starter at 19:00\nMix dough at 08:00\nFold every 45 minutes, three times\nBake at 240 C for 25 minutes", "Kitchen", false},
	{"Reading list", "The Overstory\nBraiding Sweetgrass\nA Pattern Language\nThe Dispossessed", "Books", false},
	{"Garden beds", "Bed A: beans and marigold\nBed B: carrots, then kale\nBed C: left fallow this year", "Garden", false},
	{"Quarterly goals", "Finish the kayak rack\nLearn 20 chess openings\nWalk the coast path in sections", "Plans", true},
}

type ui struct {
	win    fyne.Window
	list   *widget.List
	title  *widget.Entry
	body   *widget.Entry
	tag    *widget.Select
	pinned *widget.Check
	status *widget.Label
	tabs   *container.AppTabs
	picked int
}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Quillnote")
	w.Resize(fyne.NewSize(1100, 700))
	u := &ui{win: w, picked: 0}
	u.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Notes", theme.DocumentIcon(), u.notesTab()),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), u.settingsTab()),
		container.NewTabItemWithIcon("About", theme.InfoIcon(), u.aboutTab()),
	)
	u.tabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(u.tabs)
	u.apply(initialState())
	announceFirstFrame()
	w.ShowAndRun()
}

// apply puts the window in a named state: notes, editor, settings, about or dialog.
func (u *ui) apply(state string) {
	switch state {
	case "settings":
		u.tabs.SelectIndex(1)
	case "about":
		u.tabs.SelectIndex(2)
	case "editor":
		u.list.Select(2)
	case "dialog":
		u.list.Select(1)
		u.confirmDelete()
	}
}

func (u *ui) notesTab() fyne.CanvasObject {
	u.list = widget.NewList(
		func() int { return len(notes) },
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("tag"),
				container.NewVBox(widget.NewLabelWithStyle("title", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), widget.NewLabel("preview")))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			// Border keeps the centre first, then top, bottom, left, right
			text := row.Objects[0].(*fyne.Container)
			text.Objects[0].(*widget.Label).SetText(notes[i].Title)
			first, _, _ := strings.Cut(notes[i].Body, "\n")
			text.Objects[1].(*widget.Label).SetText(first)
			row.Objects[2].(*widget.Label).SetText(notes[i].Tag)
			icon := theme.DocumentIcon()
			if notes[i].Pinned {
				icon = theme.ConfirmIcon()
			}
			row.Objects[1].(*widget.Icon).SetResource(icon)
		},
	)
	u.list.OnSelected = func(i widget.ListItemID) { u.show(i) }

	u.title = widget.NewEntry()
	u.body = widget.NewMultiLineEntry()
	u.body.Wrapping = fyne.TextWrapWord
	u.tag = widget.NewSelect([]string{"Travel", "Kitchen", "Books", "Garden", "Plans"}, nil)
	u.pinned = widget.NewCheck("Pinned to top", nil)
	u.status = widget.NewLabel("")

	save := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() { u.save() })
	save.Importance = widget.HighImportance
	del := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() { u.confirmDelete() })
	form := widget.NewForm(
		widget.NewFormItem("Title", u.title),
		widget.NewFormItem("Tag", u.tag),
		widget.NewFormItem("", u.pinned),
	)
	editor := container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("Edit note", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), form),
		container.NewVBox(widget.NewSeparator(), container.NewHBox(u.status, layout.NewSpacer(), del, save)),
		nil, nil, u.body)

	toolbar := container.NewHBox(
		widget.NewButtonWithIcon("New note", theme.ContentAddIcon(), func() { u.newNote() }),
		layout.NewSpacer(),
		widget.NewLabel(fmt.Sprintf("%d notes", len(notes))),
	)
	left := container.NewBorder(toolbar, nil, nil, nil, u.list)
	split := container.NewHSplit(left, editor)
	split.SetOffset(0.38)
	u.list.Select(0)
	return split
}

func (u *ui) show(i int) {
	u.picked = i
	n := notes[i]
	u.title.SetText(n.Title)
	u.body.SetText(n.Body)
	u.tag.SetSelected(n.Tag)
	u.pinned.SetChecked(n.Pinned)
	u.status.SetText("Last saved 14:02")
}

func (u *ui) save() {
	notes[u.picked] = note{u.title.Text, u.body.Text, u.tag.Selected, u.pinned.Checked}
	u.list.Refresh()
	u.status.SetText("Saved")
}

func (u *ui) confirmDelete() {
	name := notes[u.picked].Title
	dialog.ShowConfirm("Delete note", fmt.Sprintf("Delete \"%s\"? This cannot be undone.", name),
		func(ok bool) {
			if !ok {
				return
			}
			notes = append(notes[:u.picked], notes[u.picked+1:]...)
			u.list.UnselectAll()
			u.list.Refresh()
			u.list.Select(0)
		}, u.win)
}

func (u *ui) newNote() {
	title := widget.NewEntry()
	title.SetPlaceHolder("Title")
	dialog.ShowForm("New note", "Create", "Cancel", []*widget.FormItem{widget.NewFormItem("Title", title)},
		func(ok bool) {
			if !ok || title.Text == "" {
				return
			}
			notes = append(notes, note{title.Text, "", "Plans", false})
			u.list.Refresh()
			u.list.Select(len(notes) - 1)
		}, u.win)
}

func (u *ui) settingsTab() fyne.CanvasObject {
	size := widget.NewSlider(10, 24)
	size.SetValue(14)
	sizeLabel := widget.NewLabel("14 pt")
	size.OnChanged = func(v float64) { sizeLabel.SetText(fmt.Sprintf("%.0f pt", v)) }
	sync := widget.NewRadioGroup([]string{"Every 5 minutes", "Every hour", "Manual only"}, nil)
	sync.SetSelected("Every hour")
	return container.NewVScroll(container.NewVBox(
		widget.NewCard("Editor", "How notes look while you write", widget.NewForm(
			widget.NewFormItem("Font size", container.NewBorder(nil, nil, nil, sizeLabel, size)),
			widget.NewFormItem("Default tag", widget.NewSelect([]string{"Travel", "Kitchen", "Books"}, nil)),
			widget.NewFormItem("Spell check", widget.NewCheck("Check as I type", nil)),
		)),
		widget.NewCard("Backup", "Where a copy of your notes is kept", widget.NewForm(
			widget.NewFormItem("Folder", widget.NewEntry()),
			widget.NewFormItem("Schedule", sync),
			widget.NewFormItem("Progress", widget.NewProgressBar()),
		)),
	))
}

func (u *ui) aboutTab() fyne.CanvasObject {
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Quillnote 1.4", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("A tiny notes app drawn by Fyne.", fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewHyperlink("Project page", nil),
	))
}
