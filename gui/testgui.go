package main

import (
	"log"

	"github.com/jroimartin/gocui"
	component "github.com/vdo/gocui-component"
)

var dView *gocui.View

type UI struct {
	*gocui.Gui
	Menu          *component.Form
	Secretform    *component.Form
	Sharelist     *component.Form
	currentSecret string
	shares        []string
	secretActive  bool
	listActive    bool
}

func (ui *UI) getShares(g *gocui.Gui, v *gocui.View) error {
	if !ui.Secretform.Validate() {
		return nil
	}
	ui.currentSecret = ui.Secretform.GetFieldText("Secret:")
	ui.listActive = true
	ui.Sharelist.Draw()
	// run sharelist !
	// The share list can be created with inputField, no editable, set text

	//ui.Secretinput.Close(g, v)
	return nil
}

func requireValidator(value string) bool {
	if value == "" {
		return false
	}
	return true
}

func main() {

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.BgColor = gocui.ColorBlack
	g.FgColor = gocui.ColorWhite
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.Mouse = true

	// terminalWidth, terminalHeight := g.Size()

	ui := &UI{
		g,
		component.NewForm(g, "menu", 2, 2, 0, 0),
		component.NewForm(g, "secret", 25, 2, 0, 0),
		component.NewForm(g, "shares", 2, 10, 0, 0),
		"",
		make([]string, 0, 24),
		false,
		false,
	}
	// Menu
	ui.Menu.AddButton("SPLIT", ui.secretInput)
	ui.Menu.AddButton("RECOVER", ui.shareInput)

	// Secret Input
	ui.Secretform.AddInputField("Secret:", 7, 52).SetMaskKeybinding(gocui.KeyCtrlA)
	ui.Secretform.Cursor = true
	ui.Secretform.AddButton("OK", ui.getShares)
	ui.Secretform.AddButton("Cancel", nil)

	ui.Sharelist.AddInputField("1:", 2, 80).SetText("0bQCkagO6d/Tt2Ays8H9qfj4KucUqpOG0CG+o7mJgTQ=").SetEditable(false)
	ui.Sharelist.AddInputField("2:", 2, 80).SetText("EVrWx0yY809A4zEfkxstogAO9fawh/OB8n4ZxEymQH0=").SetEditable(false)
	ui.Sharelist.AddInputField("3:", 2, 80).SetText("o9arB2ALaQ8Jtkj+E4Rff18RxdiIA8zQQSE9GnTCFvg=").SetEditable(false)
	ui.Sharelist.AddInputField("4:", 2, 80).SetText("vOtQ6AhdEnsZXJwYZz5VYJhDpL8beMWzj2G3JOOQFbY=").SetEditable(false)
	ui.Sharelist.AddInputField("5:", 2, 80).SetText("s9Ef/+FTg7KPlYDXBGTKwM3lBUydXJC9LvhTBIwtfSE=").SetEditable(false)
	ui.Sharelist.AddInputField("6:", 2, 80).SetText("2EzixBssq4l66uoPQzfB0mme9tVl9LTJ76bi9+vMvzk=").SetEditable(false)
	ui.Sharelist.AddInputField("7:", 2, 80).SetText("8YItD4P7+dX9T0ye587AmZpQXCOTjcDQKm3pWQOHWFc=").SetEditable(false)
	ui.Secretform.Cursor = true
	ui.Sharelist.SelBgColor = gocui.ColorWhite

	g.SetManagerFunc(ui.layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (ui *UI) secretInput(g *gocui.Gui, v *gocui.View) error {
	ui.Secretform.Draw()
	ui.secretActive = true
	return nil
}

func (ui *UI) shareInput(g *gocui.Gui, v *gocui.View) error {
	ui.Sharelist.Draw()
	return nil
}

func (ui *UI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}
	ui.Menu.Draw()

	if ui.secretActive {
		ui.Secretform.Draw()
	}
	if ui.listActive {
		ui.Sharelist.Draw()
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
