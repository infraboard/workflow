package feishu

type Color string

func (c Color) String() string {
	return string(c)
}

const (
	COLOR_BLUE      Color = "blue"
	COLOR_WATHET    Color = "wathet"
	COLOR_TURQUOISE Color = "turquoise"
	COLOR_GREEN     Color = "green"
	COLOR_YELLOW    Color = "yellow"
	COLOR_ORANGE    Color = "orange"
	COLOR_RED       Color = "red"
	COLOR_CARMINE   Color = "carmine"
	COLOR_VIOLET    Color = "violet"
	COLOR_PURPLE    Color = "purple"
	COLOR_INDIGO    Color = "indigo"
	COLOR_GREY      Color = "grey"
)
