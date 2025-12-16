package application

import (
	"log/slog"
	"math"
	"math/rand"

	"hw4-fractal-flame/internal/domain"
	"hw4-fractal-flame/internal/infrastructure"
)

type Application struct {
	Config *domain.Config
}

func NewApplication(config *domain.Config) *Application {
	return &Application{
		Config: config,
	}
}

func (a *Application) Run() error {
	transforms := a.prepareTransforms()

	wr := infrastructure.NewPngWriter()
	gen := domain.NewGenerator(transforms, a.Config.Palette.Colors, a.Config.SymmetryLevel)
	renderer := infrastructure.NewRenderer(a.Config.Palette.Colors, a.Config.Gamma, a.Config.Brightness)

	hist := a.generateParallel(gen, a.Config.Threads)

	slog.Info("Rendering image")

	img := renderer.Render(hist)

	slog.Info("Writing image to file", "path", a.Config.Output)

	err := wr.Write(img, a.Config.Output)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) prepareTransforms() []domain.Transform {
	transformations := make([]domain.Transform, len(a.Config.Transforms))

	for i, tr := range a.Config.Transforms {
		affineIdx := i
		if i >= len(a.Config.AffineTransform) {
			affineIdx = len(a.Config.AffineTransform) - 1
		}

		affine := a.Config.AffineTransform[affineIdx].Transform
		transformations[i] = domain.Composite{Transforms: []domain.Transform{affine, tr.Transform}}
	}

	return transformations
}

func (a *Application) generateParallel(generator *domain.Generator, threads int) *domain.Histogram {
	slog.Info(
		"Starting parallel generation",
		"threads", threads,
		"iterations_per_thread", a.Config.Iterations/threads,
	)

	results := make(chan *domain.Histogram, threads)
	iterationsPerThread := a.Config.Iterations / threads
	intSeed := int64(a.Config.Seed * math.MaxInt64)

	for i := range threads {
		go func(idx int64) {
			slog.Info("Generation thread started", "thread_id", idx)
			rnd := rand.New(rand.NewSource(intSeed + idx))

			hist := generator.Generate(
				a.Config.Size.Width,
				a.Config.Size.Height,
				iterationsPerThread,
				rnd,
			)
			results <- hist
		}(int64(i))
	}

	// Collect
	histograms := make([]*domain.Histogram, threads)
	for i := range threads {
		histograms[i] = <-results
	}

	slog.Info("All generation threads finished, merging histograms")
	merged := MergeHistograms(histograms)
	slog.Info("Histograms merged")

	return merged
}

func MergeHistograms(histograms []*domain.Histogram) *domain.Histogram {
	if len(histograms) == 0 {
		return nil
	}

	width := histograms[0].Width
	height := histograms[0].Height
	result := domain.NewHistogram(width, height)

	for _, h := range histograms {
		for y := range height {
			for x := range width {
				result.Count[y][x] += h.Count[y][x]
				result.ColorSum[y][x] += h.ColorSum[y][x]
			}
		}
	}

	return result
}
