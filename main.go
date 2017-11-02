package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	data      map[string]interface{}
	resource  string
	namespace string
)

func main() {
	resource = "pods"
	namespace = "default"
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.InputEsc = true
	g.SelBgColor = gocui.ColorWhite
	g.SelFgColor = gocui.ColorBlack
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", 'j', gocui.ModNone, scrollDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", 'k', gocui.ModNone, scrollUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", 'n', gocui.ModNone, openNamespaceSelector); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", 'j', gocui.ModNone, scrollDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", 'k', gocui.ModNone, scrollUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", gocui.KeyEsc, gocui.ModNone, closeView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", gocui.KeyEnter, gocui.ModNone, setNamespace); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("status", 0, 0, maxX/2, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "%s in %s", resource, namespace)
	}
	if v, err := g.SetView("main", 0, 2, maxX/2, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		fmt.Fprintln(v, "One")
		fmt.Fprintln(v, "Two")
		fmt.Fprintln(v, "Three")
		g.SetCurrentView("main")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
	_, maxY := v.Size()
	_, curY := v.Cursor()
	if curY < maxY {
		v.MoveCursor(0, 1, false)
	}
	return nil
}

func scrollUp(g *gocui.Gui, v *gocui.View) error {
	_, curY := v.Cursor()
	if curY > 0 {
		v.MoveCursor(0, -1, false)
	}
	return nil
}

func openNamespaceSelector(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.SetView("namespaceSelector", 2, 2, 20, 20); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		fmt.Fprintln(v, "One")
		fmt.Fprintln(v, "Two")
		fmt.Fprintln(v, "Three")
		if _, err = g.SetCurrentView(v.Name()); err != nil {
			return err
		}
	}
	return nil
}

func closeView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	_, err := g.SetCurrentView("main")
	return err
}

func setNamespace(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	var err error
	if namespace, err = v.Line(y); err != nil {
		return err
	} else if statusView, err := g.View("status"); err != nil {
		return err
	} else {
		statusView.Clear()
		fmt.Fprintf(statusView, "%s in %s", resource, namespace)
		return closeView(g, v)
	}
}
