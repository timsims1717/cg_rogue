package player

//type PlayerAction struct {
//	Selector selectors.Selector
//	Action   func([]world.Coords, selectors.ActionValues)
//	Values   selectors.ActionValues
//	Complete bool
//}
//
//func NewPlayerAction(sel selectors.Selector, values selectors.ActionValues, act func([]world.Coords, selectors.ActionValues)) *PlayerAction {
//	if sel == nil {
//		return nil
//	}
//	switch reflect.TypeOf(sel).Kind() {
//	case reflect.Ptr:
//		if reflect.ValueOf(sel).IsNil() {
//			return nil
//		}
//	}
//	return &PlayerAction{
//		Selector: sel,
//		Action:   act,
//		Values:   values,
//	}
//}
//
//func (p *PlayerAction) Update() {
//	if p.Complete {
//		return
//	}
//	p.Selector.Update()
//	if p.Selector.IsDone() {
//		result := p.Selector.Finish()
//		p.Action(result, p.Values)
//		p.Complete = true
//	}
//}

//type PlayerAction interface {
//	DoAction([]world.Coords)
//}
//
//type AbstractPlayerAction struct {
//	PlayerAction PlayerAction
//	RawText      string
//	text         *text.Text
//	Selector     selectors.Selector
//	Values       selectors.ActionValues
//	isDone       bool
//	start        bool
//	isCard       bool
//}
//
//func NewCardAction(sel selectors.Selector, values selectors.ActionValues, rawText string) *AbstractPlayerAction {
//	return &AbstractPlayerAction{
//		PlayerAction: nil,
//		RawText:      rawText,
//		text:         text.New(pixel.ZV, text2.BasicAtlas),
//		Selector:     sel,
//		Values:       values,
//		isDone:       true,
//		isCard:       true,
//	}
//}
//
//func NewNonCardAction(sel selectors.Selector, values selectors.ActionValues) *AbstractPlayerAction {
//	return &AbstractPlayerAction{
//		PlayerAction: nil,
//		Selector:     sel,
//		Values:       values,
//		isDone:       true,
//		isCard:       false,
//	}
//}
//
//func (pa *AbstractPlayerAction) Update() {
//	if pa.isCard {
//		pa.text.Clear()
//		pa.text.Color = colornames.Black
//		pa.text.Dot.X -= pa.text.BoundsOf(pa.RawText).W() / 2.
//		fmt.Fprintln(pa.text, pa.RawText)
//	}
//	if pa.isDone || !pa.start {
//		return
//	}
//	pa.Selector.Update()
//	if pa.Selector.IsDone() {
//		result := pa.Selector.Finish()
//		pa.PlayerAction.DoAction(result)
//		pa.isDone = true
//	}
//}
//
//func (pa *AbstractPlayerAction) Draw(canvas *pixelgl.Canvas, offset int) {
//	if pa.isCard {
//		pa.text.Draw(canvas, pixel.IM.Scaled(pa.text.Orig, 1.2).Moved(pixel.V(BaseCardWidth*0.5, BaseCardHeight*0.5-float64(offset)*(text2.BasicAtlas.LineHeight()+20.))))
//	}
//}