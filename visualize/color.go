package visualize

// A type to define color.
type Color struct {
	Red   float32
	Green float32
	Blue  float32
	Alpha float32
}

// Creating a color from red, green, blue, and alpha bytes.
func RGBA(red, green, blue, alpha float32) Color {
	return Color{
		Red:   red,
		Green: green,
		Blue:  blue,
		Alpha: alpha,
	}
}

// A version of the math.Abs function that operations on ints.
func iAbs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

// Creating a color from hue, saturation, value, and alpha bytes.
func HSVA(hue, saturation, value, alpha float32) Color {
	if saturation <= 0 {
		return RGBA(value, value, value, alpha)
	}

	// Getting a more useful hue number for calculating the RGB values directly.
	eh := int(hue)
	if eh%120 >= 60 {
		eh = 60 - eh%60
	} else {
		eh = eh % 60
	}

	// Creating the values to place into RGB.
	c := saturation * value
	x := c * (float32(eh) / 60)
	m := value - c

	// The RGB values themselves.
	var r, g, b float32
	if hue <= 60 {
		r = c
		g = x
		b = 0
	} else if hue <= 120 {
		r = x
		g = c
		b = 0
	} else if hue <= 180 {
		r = 0
		g = c
		b = x
	} else if hue <= 240 {
		r = 0
		g = x
		b = c
	} else if hue <= 300 {
		r = x
		g = 0
		b = c
	} else if hue <= 360 {
		r = c
		g = 0
		b = x
	}

	return Color{
		Red:   r + m,
		Green: g + m,
		Blue:  b + m,
		Alpha: alpha,
	}
}

// A set of pastel colors to use in the rest of the program.
var pastelWheel = []Color{}

const colors int = 6

func GetPastel(index int) Color {
	if len(pastelWheel) == 0 {
		pastelWheel = make([]Color, colors)

		var h float32 = 0.0
		for i := 0; i < colors; i++ {
			pastelWheel[i] = HSVA(h, 0.6, 0.9, 1.0)
			h += 360 / float32(colors)
		}
	}

	return pastelWheel[index%colors]
}
