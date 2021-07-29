package gui

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed Sarasa-Mono-SC-Bold.ttf
var bold []byte

//go:embed Sarasa-Mono-SC-Regular.ttf
var regular []byte

var boldfont = &fyne.StaticResource{
	StaticName:    "NotoSansCJK-Bold",
	StaticContent: bold,
}

var regularfont = &fyne.StaticResource{
	StaticName:    "NotoSansCJK-Regular",
	StaticContent: regular,
}

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

// return bundled font resource
func (*myTheme) Font(s fyne.TextStyle) fyne.Resource {
	// if s.Monospace {
	// 	return theme.DefaultTheme().Font(s)
	// }
	if s.Bold {
		if s.Italic {
			return regularfont
		}
		return boldfont
	}
	// if s.Italic {
	// 	return theme.DefaultTheme().Font(s)
	// }
	return regularfont
}

func (*myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
