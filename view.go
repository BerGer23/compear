package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var initialized = false

var leftEntry *widget.Entry
var rightEntry *widget.Entry
var separatorLeft *widget.Select
var separatorRight *widget.Select
var resultListLeft *widget.List
var resultListRight *widget.List
var separatorAutodetectLeft *widget.Check
var separatorAutodetectRight *widget.Check
var ignoreTimestamps *widget.Check
var trim *widget.Check
var result Analysis

func analyze() {
	if initialized {
		result = compareTokens(leftEntry.Text, rightEntry.Text, separatorLeft.Selected, separatorRight.Selected, trim.Checked)
		//TODO: dynamic
		resultListLeft.Resize(fyne.Size{Width: 500, Height: 300})
		resultListRight.Resize(fyne.Size{Width: 500, Height: 300})
	}
}

func processResult(result Token) string {
	toPrint := result.Content
	if len(result.Content) == 0 {
		toPrint = "<empty>"
	}
	return toPrint + " (" + strconv.Itoa(result.Index) + ")"
}

func setupView() {
	a := app.New()
	w := a.NewWindow("Compear")

	leftEntry = widget.NewMultiLineEntry()
	leftEntry.SetMinRowsVisible(20)
	leftEntry.SetPlaceHolder("Apples")

	rightEntry = widget.NewMultiLineEntry()
	rightEntry.SetMinRowsVisible(20)
	rightEntry.SetPlaceHolder("Oranges")

	separatorLeft = widget.NewSelect([]string{"Newline", "Comma", "Space"}, func(value string) {
		log.Println("SeparatorType set to", value)
	})
	separatorLeft.SetSelectedIndex(0)

	separatorRight = widget.NewSelect([]string{"Newline", "Comma", "Space"}, func(value string) {
		log.Println("SeparatorType set to", value)
	})
	separatorRight.SetSelectedIndex(0)

	resultListLeft = widget.NewList(
		func() int {
			return len(result.FindingsLeft)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(processResult(result.FindingsLeft[i]))
		})

	resultListRight = widget.NewList(
		func() int {
			return len(result.FindingsRight)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(processResult(result.FindingsRight[i]))
		})

	separatorAutodetectLeft = widget.NewCheck("Detect separator", func(bool) {
		analyze()
	})
	separatorAutodetectLeft.SetChecked(true)
	separatorAutodetectRight = widget.NewCheck("Detect separator", func(bool) {
		analyze()
	})
	separatorAutodetectRight.SetChecked(true)

	ignoreTimestamps = widget.NewCheck("Ignore timestamps", func(bool) {
		analyze()
	})
	ignoreTimestamps.SetChecked(false)
	ignoreTimestamps.Disable()

	trim = widget.NewCheck("Trim Token values", func(v bool) {
		analyze()
	})
	trim.SetChecked(true)

	leftEntry.OnChanged = func(v string) {
		if separatorAutodetectLeft.Checked {
			separatorLeft.SetSelected(detectSeparator(v))
		}
		analyze()
	}
	rightEntry.OnChanged = func(v string) {
		if separatorAutodetectRight.Checked {
			separatorRight.SetSelected(detectSeparator(v))
		}
		analyze()
	}

	optionsLeft := container.NewHBox()
	optionsLeft.Add(separatorAutodetectLeft)
	optionsLeft.Add(separatorLeft)

	leftContainer := container.NewVBox()
	leftContainer.Add(leftEntry)
	leftContainer.Add(optionsLeft)

	optionsRight := container.NewHBox()
	optionsRight.Add(separatorAutodetectRight)
	optionsRight.Add(separatorRight)

	rightContainer := container.NewVBox()
	rightContainer.Add(rightEntry)
	rightContainer.Add(optionsRight)

	newGrid := container.New(layout.NewGridLayout(2), leftContainer, rightContainer)

	generalOptionsContainer := container.NewHBox()
	generalOptionsContainer.Add(trim)
	generalOptionsContainer.Add(ignoreTimestamps)

	findingsContainer := container.New(layout.NewGridLayout(2), resultListLeft, resultListRight)

	w.Resize(fyne.NewSize(1000, 1400))
	mainContainer := container.NewVBox()
	mainContainer.Add(newGrid)
	mainContainer.Add(generalOptionsContainer)
	mainContainer.Add(findingsContainer)
	w.SetContent(mainContainer)
	r, err := fyne.LoadResourceFromPath("./icon.png")
	if err != nil {
		log.Println("Error: " + err.Error())
	}
	w.SetIcon(r)
	a.SetIcon(r)
	//TODO: wieso geht das noch nicht?
	initialized = true
	w.ShowAndRun()

}
