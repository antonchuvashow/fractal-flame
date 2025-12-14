package config

// DTO is a structure for json config.
type DTO struct {
	Size       SizeDTO        `json:"size"`
	Iterations int            `json:"iteration_count"`
	Output     string         `json:"output_path"`
	Threads    int            `json:"threads"`
	Seed       float64        `json:"seed"`
	Transforms []TransformDTO `json:"functions"`
	Symmetry   int            `json:"symmetry_level"`
	Gamma      float64        `json:"gamma"`
	Brightness float64        `json:"brightness"`
	Palette    string         `json:"palette"`
	Affine     []AffineDTO    `json:"affine_params"`
}

type SizeDTO struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type TransformDTO struct {
	Name       string  `json:"name"`
	Weight     float64 `json:"weight"`
	ColorIndex float64 `json:"color_index"`
}

type AffineDTO struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
	D float64 `json:"d"`
	E float64 `json:"e"`
	F float64 `json:"f"`
}
