package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/nsgo-dev/nsgo/core"
	"github.com/nsgo-dev/nsgo/icons/lucide"
)

// TestIconButtonShowsItsLabelAsATooltip shows the label of a ghost button and of a circle once the pointer has rested on each, and not before.
func TestIconButtonShowsItsLabelAsATooltip(t *testing.T) {
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		return core.View{Style: core.Style{AlignItems: core.AlignStart, Padding: core.All(40), Gap: 12}, Children: []core.Node{
			IconButton{Icon: lucide.X, Label: "Close", TestID: "close"},
			IconButton{Icon: lucide.Plus, Label: "Add", Variant: IconButtonPrimary, TestID: "add"},
			IconButton{Icon: lucide.Minus, Label: "Remove", Variant: IconButtonRaised, TestID: "remove"},
		}}
	})
	for _, c := range []struct{ id, label string }{{"close", "Close"}, {"add", "Add"}, {"remove", "Remove"}} {
		if app.Exists("text:" + c.label) {
			t.Fatalf("%s: a tooltip shows before the pointer rests on it", c.id)
		}
		app.Hover("id:" + c.id)
		app.Advance(600 * time.Millisecond)
		if !app.Exists("text:" + c.label) {
			t.Fatalf("%s: no tooltip says %q", c.id, c.label)
		}
		app.HoverOut()
		app.Settle()
	}
}

// TestIconButtonTooltipOverridesItsLabel shows a Tooltip in place of the label.
func TestIconButtonTooltipOverridesItsLabel(t *testing.T) {
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		return core.View{Style: core.Style{AlignItems: core.AlignStart, Padding: core.All(40)}, Children: []core.Node{
			IconButton{Icon: lucide.X, Label: "Close", Tooltip: func() string { return "Close this panel" }, TestID: "close"},
		}}
	})
	app.Hover("id:close")
	app.Advance(600 * time.Millisecond)
	app.Find("text:Close this panel")
	if app.Exists(`text:"Close"`) {
		t.Fatal("the label shows beside the Tooltip")
	}
}

// TestIconButtonTooltipFollowsItsText follows a LabelFunc and a Tooltip as they change, showing nothing while the text is empty.
func TestIconButtonTooltipFollowsItsText(t *testing.T) {
	var label, tip *core.Signal[string]
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		label, tip = core.NewSignal(cx, "Play"), core.NewSignal(cx, "Mute")
		return core.View{Style: core.Style{AlignItems: core.AlignStart, Padding: core.All(40), Gap: 12}, Children: []core.Node{
			IconButton{Icon: lucide.Play, LabelFunc: label.Get, TestID: "play"},
			IconButton{Icon: lucide.Volume2, Label: "Volume", Tooltip: tip.Get, Variant: IconButtonRaised, TestID: "volume"},
		}}
	})
	app.Hover("id:play")
	app.Advance(600 * time.Millisecond)
	app.Find("text:Play")
	app.Post(func() { label.Set("Pause") })
	app.Settle()
	app.Find("text:Pause")
	if app.Exists(`text:"Play"`) {
		t.Fatal("the tooltip keeps the old label")
	}
	app.HoverOut()
	app.Settle()

	before := drawn(app.Dump())
	app.Hover("id:volume")
	app.Advance(600 * time.Millisecond)
	app.Find("text:Mute")
	app.Post(func() { tip.Set("") })
	app.Settle()
	if app.Exists("text:Mute") || drawn(app.Dump()) != before {
		t.Fatalf("a bubble shows for an empty Tooltip:\n%s", app.Dump())
	}
	app.Post(func() { tip.Set("Unmute") })
	app.Settle()
	app.Find("text:Unmute")
}

// TestIconButtonShowsNoTooltipForAnEmptyHint keeps a Tooltip that says nothing, or a button with no label, from showing a bubble.
func TestIconButtonShowsNoTooltipForAnEmptyHint(t *testing.T) {
	app, _ := start(t, desktopEnv, func(cx *core.Scope) core.Node {
		return core.View{Style: core.Style{AlignItems: core.AlignStart, Padding: core.All(40), Gap: 12}, Children: []core.Node{
			IconButton{Icon: lucide.X, Label: "Close", Tooltip: func() string { return "" }, TestID: "close"},
			IconButton{Icon: lucide.Plus, TestID: "bare"},
		}}
	})
	before := drawn(app.Dump())
	for _, id := range []string{"close", "bare"} {
		app.Hover("id:" + id)
		app.Advance(600 * time.Millisecond)
		if after := drawn(app.Dump()); after != before {
			t.Fatalf("%s: a bubble shows as the pointer rests on it:\n%s", id, after)
		}
	}
}

// drawn is a Dump without its empty views, such as an overlay that holds nothing.
func drawn(dump string) string {
	var b strings.Builder
	for line := range strings.Lines(dump) {
		if !strings.HasSuffix(strings.TrimSpace(line), " 0x0") {
			b.WriteString(line)
		}
	}
	return b.String()
}
