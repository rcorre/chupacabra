package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

var ()

func init() {
}

func main() {
	conf := NewConfig()
	kube, err := NewKube(conf.KubeConfigPath)
	if err != nil {
		log.Panicln(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	if err != nil {
		panic(err)
	}

	g.InputEsc = true
	g.SelBgColor = gocui.ColorWhite
	g.SelFgColor = gocui.ColorBlack
	m := NewManager(g, kube)
	g.SetManager(m)
	bindKeys(m.(*manager), g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
