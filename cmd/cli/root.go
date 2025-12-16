package cli

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/spf13/cobra"

	"hw4-fractal-flame/internal/application"
	"hw4-fractal-flame/internal/domain"
	"hw4-fractal-flame/internal/infrastructure"
	"hw4-fractal-flame/internal/infrastructure/config"
)

const (
	defaultWidth      = 1920
	defaultHeight     = 1080
	defaultIterations = 250000
	defaultGamma      = 2.2
)

// Описания для сложных параметров
const (
	functionsDescription = `Transform functions to apply.

Format: name:weight,name:weight,...

Available transforms:
  linear      - Identity transform (weight: 0.0-1.0)
  sinusoidal  - Sine-based distortion
  spherical   - Spherical inversion
  swirl       - Swirling effect
  horseshoe   - Horseshoe-shaped distortion
  polar       - Polar coordinate transform
  handkerchief - Handkerchief pattern
  heart       - Heart-shaped transform
  disc        - Disc projection
  spiral      - Spiral effect
  hyperbolic  - Hyperbolic transform
  diamond     - Diamond pattern
  ex          - Ex transform
  julia       - Julia set variation

Example: --functions "swirl:1.0,horseshoe:0.8,spherical:0.5"`

	affineDescription = `Affine transformation coefficients.

Format: a,b,c,d,e,f;a,b,c,d,e,f;...

Each transform is defined by 6 coefficients:
  x' = a*x + b*y + c
  y' = d*x + e*y + f

Coefficients explanation:
  a,e - scaling (diagonal)
  b,d - rotation/shear (off-diagonal)  
  c,f - translation

Multiple transforms separated by semicolon (;)

Example: --affine "0.5,0.0,0.25,0.0,0.5,0.0;0.5,0.0,-0.25,0.0,0.5,0.5"`
)

func GetRootCommand() (*cobra.Command, error) {
	var (
		width        int
		height       int
		iterations   int
		output       string
		threads      int
		seed         float64
		symmetry     int
		gamma        float64
		brightness   float64
		palette      string
		functions    string
		affineParams string
		configPath   string
	)

	cmd := &cobra.Command{
		Use:   "fractalflame",
		Short: "A generator for fractal flame images",
		Long: `Fractal Flame Generator

Generates fractal flame images using the chaos game algorithm
with nonlinear transformations and color mapping.

The algorithm works by:
1. Applying random affine transformations to a point
2. Applying nonlinear variation functions
3. Accumulating color and hit counts on a histogram
4. Applying gamma correction and rendering

For complex configurations, use a JSON config file with --config.`,

		Example: `  # Basic usage with default settings
  fractalflame -o flame.png

  # Custom transforms and affine parameters
  fractalflame \
    --functions "swirl:1.0,spherical:0.5" \
    --affine "0.5,0.1,0,0,0.5,0;-0.5,0,0.25,0,0.5,0.5" \
    --width 1920 --height 1080 \
    --iterations 5000 \
    -o my_fractal.png

  # Using a config file
  fractalflame --config fractal_config.json

  # High quality render with multiple threads
  fractalflame \
    -f "disk:1.0,swirl:0.7" \
    -a "0.8,0.2,0,-0.2,0.8,0" \
    -i 10000 -t 8 --gamma 2.5 \
    -o hq_fractal.png

  # Sierpinski triangle approximation
  fractalflame \
    --affine "0.5,0,0,0,0.5,0;0.5,0,0.5,0,0.5,0;0.5,0,0.25,0,0.5,0.433" \
    --functions "linear:1.0,linear:1.0,linear:1.0" \
    --symmetry 3`,

		RunE: func(*cobra.Command, []string) error {
			intSeed := int64(seed * math.MaxInt64)
			rnd := rand.New(rand.NewSource(intSeed))
			cfg := &domain.Config{}

			var err error

			if configPath != "" {
				cfg, err = config.Load(configPath)

				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}
			} else {
				transforms, err := ParseTransforms(functions, rnd)

				if err != nil {
					return fmt.Errorf("failed to parse functions: %w", err)
				}

				affines, err := ParseAffineParams(affineParams, rnd)
				if err != nil {
					return fmt.Errorf("failed to parse affine params: %w", err)
				}

				cfg = &domain.Config{
					Size:            domain.Size{Width: width, Height: height},
					Iterations:      iterations,
					Output:          output,
					Threads:         threads,
					Seed:            seed,
					Transforms:      transforms,
					AffineTransform: affines,
					SymmetryLevel:   symmetry,
					Gamma:           gamma,
					Brightness:      brightness,
					Palette: domain.Palette{
						Name:   palette,
						Colors: infrastructure.GetPaletteByName(palette),
					},
				}
			}

			app := application.NewApplication(cfg)

			err = app.Run()
			if err != nil {
				return fmt.Errorf("application run failed: %w", err)
			}

			return nil
		},
	}

	// Флаги с подробными описаниями
	cmd.Flags().IntVarP(&width, "width", "w", defaultWidth,
		"Width of the output image in pixels")
	cmd.Flags().IntVar(&height, "height", defaultHeight,
		"Height of the output image in pixels")
	cmd.Flags().IntVarP(&iterations, "iterations", "i", defaultIterations,
		"Number of iterations per sample (higher = more detail)")
	cmd.Flags().StringVarP(&output, "output", "o", "fractal.png",
		"Output file path (supports .png)")
	cmd.Flags().IntVarP(&threads, "threads", "t", 1,
		"Number of parallel threads (0 = auto)")
	cmd.Flags().Float64Var(&seed, "seed", 0,
		"Random seed for reproducible results (0-1)")
	cmd.Flags().IntVarP(&symmetry, "symmetry", "s", 1,
		"Rotational symmetry level (1 = none, 2+ = N-fold symmetry)")
	cmd.Flags().Float64VarP(&gamma, "gamma", "g", defaultGamma,
		"Gamma correction value (typically 2.0-3.0)")
	cmd.Flags().Float64VarP(&brightness, "brightness", "b", 1.0,
		"Brightness multiplier (0.5-2.0 recommended)")
	cmd.Flags().StringVarP(&palette, "palette", "p", "rainbow",
		"Color palette: fire, ocean, rainbow, grayscale")
	cmd.Flags().StringVarP(&configPath, "config", "c", "",
		"Path to JSON config file (overrides other flags)")

	// Сложные параметры с многострочными описаниями
	cmd.Flags().StringVarP(&functions, "functions", "f",
		"linear:1.0,linear:1.0,linear:1.0",
		functionsDescription)
	cmd.Flags().StringVarP(&affineParams, "affine", "a",
		"0.5,0,0,0,0.5,0;0.5,0,0.5,0,0.5,0;0.5,0,0.25,0,0.5,0.433",
		affineDescription)

	return cmd, nil
}
