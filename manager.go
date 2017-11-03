package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type manager struct {
	data      map[string]interface{}
	resource  string
	namespace string
	kube      Kube
}

func NewManager(gui *gocui.Gui, kube Kube) gocui.Manager {
	m := &manager{
		kube:      kube,
		resource:  "pods",
		namespace: "default",
	}

	m.loadResources()

	return m
}

func (m *manager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("status", 0, 0, maxX/2, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "%s in %s", m.resource, m.namespace)
	}
	if v, err := g.SetView("main", 0, 2, maxX/2, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		g.SetCurrentView("main")
	}
	return nil
}

func (m *manager) loadResources() error {
	data, err := m.kube.Get(m.resource, m.namespace)
	if err != nil {
		return err
	}
	m.data = data
	return nil
}

func bindKeys(m *manager, g *gocui.Gui) {
	if err := g.SetKeybinding("main", 'j', gocui.ModNone, m.scrollDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", 'k', gocui.ModNone, m.scrollUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", 'n', gocui.ModNone, m.openNamespaceSelector); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", gocui.KeyTab, gocui.ModNone, m.openResourceSelector); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", 'j', gocui.ModNone, m.scrollDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", 'k', gocui.ModNone, m.scrollUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("resourceSelector", 'j', gocui.ModNone, m.scrollDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("resourceSelector", 'k', gocui.ModNone, m.scrollUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", gocui.KeyEsc, gocui.ModNone, m.closeView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("namespaceSelector", gocui.KeyEnter, gocui.ModNone, m.setNamespace); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("resourceSelector", gocui.KeyEsc, gocui.ModNone, m.closeView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("resourceSelector", gocui.KeyEnter, gocui.ModNone, m.setResource); err != nil {
		log.Panicln(err)
	}
}

func (m *manager) scrollDown(g *gocui.Gui, v *gocui.View) error {
	_, maxY := v.Size()
	_, curY := v.Cursor()
	if curY < maxY {
		v.MoveCursor(0, 1, false)
	}
	return nil
}

func (m *manager) scrollUp(g *gocui.Gui, v *gocui.View) error {
	_, curY := v.Cursor()
	if curY > 0 {
		v.MoveCursor(0, -1, false)
	}
	return nil
}

func (m *manager) openNamespaceSelector(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.SetView("namespaceSelector", 2, 2, 40, 40); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true

		namespaces, err := m.kube.Get("namespaces", "")
		if err != nil {
			return err
		}

		for k := range namespaces {
			fmt.Fprintln(v, k)
		}

		if _, err = g.SetCurrentView(v.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) openResourceSelector(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.SetView("resourceSelector", 2, 2, 40, 40); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true

		for _, name := range m.kube.Resources() {
			fmt.Fprintln(v, name)
		}

		if _, err = g.SetCurrentView(v.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) closeView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}
	_, err := g.SetCurrentView("main")
	return err
}

func (m *manager) setNamespace(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	namespace, err := v.Line(y)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		m.namespace = namespace
		statusView, err := g.View("status")
		if err != nil {
			return err
		}
		statusView.Clear()
		fmt.Fprintf(statusView, "%s in %s", m.resource, m.namespace)

		mainView, err := g.View("main")
		if err != nil {
			return err
		}
		mainView.Clear()
		for k := range m.data {
			fmt.Fprintln(mainView, k)
		}

		return m.loadResources()
	})
	return m.closeView(g, v)
}

func (m *manager) setResource(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	resource, err := v.Line(y)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		m.resource = resource
		statusView, err := g.View("status")
		if err != nil {
			return err
		}
		statusView.Clear()
		fmt.Fprintf(statusView, "%s in %s", m.resource, m.namespace)

		mainView, err := g.View("main")
		if err != nil {
			return err
		}
		mainView.Clear()
		for k := range m.data {
			fmt.Fprintln(mainView, k)
		}

		return m.loadResources()
	})
	return m.closeView(g, v)
}
