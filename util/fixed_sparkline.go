package util

import (
	. "github.com/gizak/termui/v3"
	. "github.com/gizak/termui/v3/widgets"

	"image"
	"math"
)

type FixedSparklineGroup struct {
	SparklineGroup
}

func (self *FixedSparklineGroup) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	sparklineHeight := self.Inner.Dy() / len(self.Sparklines)

	for i, sl := range self.Sparklines {
		heightOffset := (sparklineHeight * (i + 1))
		barHeight := sparklineHeight
		if i == len(self.Sparklines)-1 {
			heightOffset = self.Inner.Dy()
			barHeight = self.Inner.Dy() - (sparklineHeight * i)
		}
		if sl.Title != "" {
			barHeight--
		}

		maxVal := sl.MaxVal
		if maxVal == 0 {
			maxVal, _ = GetMaxFloat64FromSlice(sl.Data)
		}

		// draw line
		for j := 0; j < len(sl.Data) && j < self.Inner.Dx(); j++ {
			data := sl.Data[j]
			height := (data / maxVal) * float64(barHeight)
			sparkChar := BARS[len(BARS)-1]
			k := 0
			for ; k < int(height); k++ {
				buf.SetCell(
					NewCell(sparkChar, NewStyle(sl.LineColor)),
					image.Pt(j+self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset-k),
				)
			}
			if height == 0 {
				sparkChar = BARS[1]
				buf.SetCell(
					NewCell(sparkChar, NewStyle(sl.LineColor)),
					image.Pt(j+self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset),
				)
			} else {
				kHeight := height - math.Floor(height)
				sparkChar = BARS[int(kHeight*float64(len(BARS)-1))]
				buf.SetCell(
					NewCell(sparkChar, NewStyle(sl.LineColor)),
					image.Pt(j+self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset-k),
				)
			}
		}

		if sl.Title != "" {
			// draw title
			buf.SetString(
				TrimString(sl.Title, self.Inner.Dx()),
				sl.TitleStyle,
				image.Pt(self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset-barHeight),
			)
		}
	}
}
