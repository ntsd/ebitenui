package widget

import (
	img "image"

	"github.com/ebitenui/ebitenui/image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ProgressBar struct {
	Min     int
	Max     int
	current int

	widgetOpts []WidgetOpt
	trackImage *ProgressBarImage
	fillImage  *ProgressBarImage

	trackPadding Insets

	init        *MultiOnce
	widget      *Widget
	lastCurrent int
}

type ProgressBarImage struct {
	Idle     *image.NineSlice
	Hover    *image.NineSlice
	Disabled *image.NineSlice
}

type ProgressBarOpt func(s *ProgressBar)

type ProgressBarOptions struct {
}

var ProgressBarOpts ProgressBarOptions

func NewProgressBar(opts ...ProgressBarOpt) *ProgressBar {
	pb := &ProgressBar{
		Min:     1,
		Max:     100,
		current: 1,

		trackImage: &ProgressBarImage{},
		fillImage:  &ProgressBarImage{},

		lastCurrent: 1,

		init: &MultiOnce{},
	}

	pb.init.Append(pb.createWidget)

	for _, o := range opts {
		o(pb)
	}

	return pb
}

func (o ProgressBarOptions) WidgetOpts(opts ...WidgetOpt) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.widgetOpts = append(s.widgetOpts, opts...)
	}
}

func (o ProgressBarOptions) Images(track *ProgressBarImage, fill *ProgressBarImage) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.trackImage = track
		s.fillImage = fill
	}
}

func (o ProgressBarOptions) TrackPadding(i Insets) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.trackPadding = i
	}
}

func (o ProgressBarOptions) Values(min int, max int, current int) ProgressBarOpt {
	return func(s *ProgressBar) {
		s.Min = min
		s.Max = max
		s.current = current
	}
}

func (s *ProgressBar) Configure(opts ...ProgressBarOpt) {
	for _, o := range opts {
		o(s)
	}
}

func (s *ProgressBar) GetWidget() *Widget {
	s.init.Do()
	return s.widget
}

func (s *ProgressBar) PreferredSize() (int, int) {
	w, h := s.trackImage.Idle.MinSize()
	if s.widget != nil {
		if w < s.widget.MinWidth {
			w = s.widget.MinWidth
		}
		if h < s.widget.MinHeight {
			h = s.widget.MinHeight
		}
	}
	return w, h
}

func (s *ProgressBar) SetLocation(rect img.Rectangle) {
	s.init.Do()
	s.widget.Rect = rect
}

func (s *ProgressBar) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	s.init.Do()
	s.widget.Render(screen, def)
	s.draw(screen)
}

func (s *ProgressBar) draw(screen *ebiten.Image) {
	i := s.trackImage.Idle
	fill := s.fillImage.Idle
	if s.widget.Disabled {
		if s.trackImage.Disabled != nil {
			i = s.trackImage.Disabled
		}
		if s.fillImage.Disabled != nil {
			fill = s.fillImage.Disabled
		}
	}

	if i != nil {
		i.Draw(screen, s.widget.Rect.Dx(), s.widget.Rect.Dy(), s.widget.drawImageOptions)
	}
	if fill != nil && s.currentPercentage() > 0 {
		fillX := s.widget.Rect.Dx() - s.trackPadding.Left - s.trackPadding.Right
		fillX = int(float64(fillX) * s.currentPercentage())

		fillY := s.widget.Rect.Dy() - s.trackPadding.Top - s.trackPadding.Bottom
		fill.Draw(screen, fillX, fillY, s.drawFillOptions)
	}
}

func (s *ProgressBar) drawFillOptions(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(float64(s.GetWidget().Rect.Min.X+s.trackPadding.Left), float64(s.GetWidget().Rect.Min.Y+s.trackPadding.Top))
}

func (s *ProgressBar) currentPercentage() float64 {
	if s.Max <= s.Min {
		return 0
	}

	return float64(s.current-s.Min) / float64(s.Max-s.Min)
}

func (s *ProgressBar) SetCurrent(value int) bool {
	oldValue := s.current
	if value < s.Min {
		s.current = s.Min
	} else if value > s.Max {
		s.current = s.Max
	} else {
		s.current = value
	}
	return oldValue != s.current
}

func (s *ProgressBar) GetCurrent() int {
	return s.current
}

func (s *ProgressBar) createWidget() {
	s.widget = NewWidget(s.widgetOpts...)
}
