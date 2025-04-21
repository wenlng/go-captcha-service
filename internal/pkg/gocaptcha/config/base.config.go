/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package config

import (
	"github.com/wenlng/go-captcha/v2/base/option"
)

// RangeVal .
type RangeVal struct {
	Min, Max int
}

// Size .
type Size struct {
	Width, Height int
}

// Point .
type Point struct {
	X, Y int
}

// ClickMasterOption .
type ClickMasterOption struct {
	ImageSize             option.Size       `json:"image_size"`
	RangeLength           option.RangeVal   `json:"range_length"`
	RangeAngles           []option.RangeVal `json:"range_angles"`
	RangeSize             option.RangeVal   `json:"range_size"`
	RangeColors           []string          `json:"range_colors"`
	DisplayShadow         bool              `json:"display_shadow"`
	ShadowColor           string            `json:"shadow_color"`
	ShadowPoint           option.Point      `json:"shadow_point"`
	ImageAlpha            float32           `json:"image_alpha"`
	UseShapeOriginalColor bool              `json:"use_shape_original_color"`
}

// ClickThumbOption .
type ClickThumbOption struct {
	ImageSize                 option.Size     `json:"image_size"`
	RangeVerifyLength         option.RangeVal `json:"range_verify_length"`
	DisabledRangeVerifyLength bool            `json:"disabled_range_verify_length"`
	RangeTextSize             option.RangeVal `json:"range_text_size"`
	RangeTextColors           []string        `json:"range_text_colors"`
	RangeBackgroundColors     []string        `json:"range_background_colors"`
	BackgroundDistort         int             `json:"background_distort"`
	BackgroundDistortAlpha    float32         `json:"background_distort_alpha"`
	BackgroundCirclesNum      int             `json:"background_circles_num"`
	BackgroundSlimLineNum     int             `json:"background_slim_line_num"`
	IsThumbNonDeformAbility   bool            `json:"is_thumb_non_deform_ability"`
}

// ClickConfig .
type ClickConfig struct {
	Version  string            `json:"version"`
	Language string            `json:"language"`
	Master   ClickMasterOption `json:"master"`
	Thumb    ClickThumbOption  `json:"thumb"`
}

// SlideMasterOption .
type SlideMasterOption struct {
	ImageSize  option.Size `json:"image_size"`
	ImageAlpha float32     `json:"image_alpha"`
}

// SlideThumbOption .
type SlideThumbOption struct {
	RangeGraphSizes           option.RangeVal   `json:"range_graph_size"`
	RangeGraphAngles          []option.RangeVal `json:"range_graph_angles"`
	GenerateGraphNumber       int               `json:"generate_graph_number"`
	EnableGraphVerticalRandom bool              `json:"enable_graph_vertical_random"`
	RangeDeadZoneDirections   []string          `json:"range_dead_zone_directions"`
}

// SlideConfig .
type SlideConfig struct {
	Version string            `json:"version"`
	Master  SlideMasterOption `json:"master"`
	Thumb   SlideThumbOption  `json:"thumb"`
}

// RotateMasterOption .
type RotateMasterOption struct {
	ImageSquareSize int `json:"image_square_size"`
}

// RotateThumbOption .
type RotateThumbOption struct {
	RangeAngles           []option.RangeVal `json:"range_angles"`
	RangeImageSquareSizes []int             `json:"range_image_square_sizes"`
	ImageAlpha            float32           `json:"image_alpha"`
}

// RotateConfig .
type RotateConfig struct {
	Version string             `json:"version"`
	Master  RotateMasterOption `json:"master"`
	Thumb   RotateThumbOption  `json:"thumb"`
}
