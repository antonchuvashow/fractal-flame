package domain

type Point struct {
	X, Y, ColorIndex float64
}

type Histogram struct {
	Width    int
	Height   int
	Count    [][]uint32
	ColorSum [][]float64
}

func NewHistogram(width, height int) *Histogram {
	count := make([][]uint32, height)
	colorSum := make([][]float64, height)

	for y := range height {
		count[y] = make([]uint32, width)
		colorSum[y] = make([]float64, width)
	}

	return &Histogram{
		Width:    width,
		Height:   height,
		Count:    count,
		ColorSum: colorSum,
	}
}

func (h *Histogram) Contains(x, y int) bool {
	return x >= 0 && x < h.Width && y >= 0 && y < h.Height
}

func (h *Histogram) Increment(x, y int, colorIndex float64) {
	h.Count[y][x]++
	h.ColorSum[y][x] += colorIndex
}

func (h *Histogram) GetCount(x, y int) uint32 {
	return h.Count[y][x]
}

func (h *Histogram) GetAverageColorIndex(x, y int) float64 {
	if h.Count[y][x] == 0 {
		return 0
	}

	return h.ColorSum[y][x] / float64(h.Count[y][x])
}

// MaxCount get maximum amount of hits.
func (h *Histogram) MaxCount() uint32 {
	var maxCount uint32

	for y := range h.Height {
		for x := range h.Width {
			if h.Count[y][x] > maxCount {
				maxCount = h.Count[y][x]
			}
		}
	}

	return maxCount
}
