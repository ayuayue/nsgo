package ui

import (
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/icons/lucide"
)

// TestIconButtonShowsItsLabelAsATooltip shows the label of a ghost button
// and of a circle as the pointer rests on each.
func TestIconButtonShowsItsLabelAsATooltip(t *testing.T) {
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		return core.View{Style: core.Style{AlignItems: core.AlignStart, Padding: core.All(40), Gap: 12}, Children: []core.Node{
			IconButton{Icon: lucide.X, Label: "Close", TestID: "close"},
			IconButton{Icon: lucide.Plus, Label: "Add", Variant: IconButtonPrimary, TestID: "add"},
		}}
	})
	app.Hover("id:close")
	app.Advance(600 * time.Millisecond)
	app.Find("text:Close")

	app.Hover("id:add")
	app.Advance(600 * time.Millisecond)
	app.Find("text:Add")
}

// TestIconButtonTooltipOverridesItsLabel shows a Tooltip in place of the
// label.
func TestIconButtonTooltipOverridesItsLabel(t *testing.T) {
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		return core.View{Style: core.Style{AlignItems: core.AlignStart, Padding: core.All(40)}, Children: []core.Node{
			IconButton{Icon: lucide.X, Label: "Close", Tooltip: func() string { return "Close this panel" }, TestID: "close"},
		}}
	})
	app.Hover("id:close")
	app.Advance(600 * time.Millisecond)
	app.Find("text:Close this panel")
}

// TestIconButtonShowsNoTooltipForAnEmptyHint keeps a Tooltip that says
// nothing from showing a bubble.
func TestIconButtonShowsNoTooltipForAnEmptyHint(t *testing.T) {
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		return core.View{Style: core.Style{AlignItems: core.AlignStart, Padding: core.All(40)}, Children: []core.Node{
			IconButton{Icon: lucide.X, Label: "Close", Tooltip: func() string { return "" }, TestID: "close"},
		}}
	})
	app.Hover("id:close")
	app.Advance(600 * time.Millisecond)
	if app.Exists("text:Close") {
		t.Fatal("a tooltip shows for an empty hint")
	}
}
