package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	configContent := `
{
  "dir": "./resources",
  "resources": {
    "char": {
      "type": "chinese",
      "text": []
    },
    "font": {
      "type": "load",
      "file_dir": "./fonts/",
      "file_maps": {
        "aa": "aa.ttf",
        "bb": "bb.ttf"
      }
    },
    "shape_image": {
      "type": "load",
      "file_dir": "./shape_images/",
      "file_maps": {
        "aa": "aa.png",
        "bb":"bb.png"
      }
    },
    "master_image": {
      "type": "load",
      "file_dir": "./master_images/",
      "file_maps": {
        "aa": "aa.png",
        "bb": "bb.png"
      }
    },
    "thumb_image": {
      "type": "load",
      "file_dir": "./thumb_images/",
      "file_maps": {
        "aa": "aa.png",
        "bb": "bb.png"
      }
    },
    "tile_image": {
      "type": "load",
      "file_dir": "./tile_images/",
      "file_maps_01": {
        "overlay_image_01": "overlay_image_01.png",
        "overlay_image_02": "overlay_image_02.png"
      },
      "file_maps_02": {
        "shadow_image_01": "shadow_image_01.png",
        "shadow_image_02": "shadow_image_02.png"
      },
      "file_maps_03": {
        "mask_image_01": "mask_image_01.png"
      }
    }
  },
  "builder": {
    "click_config_maps": {
      "default_cn_click": {
        "master": {
          "image_size": {
            "width": 300,
            "height": 200
          },
          "range_length": {
            "min": 6,
            "max": 7
          },
          "range_angles": [
            {
              "min": 20,
              "max": 35
            },
            {
              "min": 35,
              "max": 45
            },
            {
              "min": 290,
              "max": 305
            },
            {
              "min": 305,
              "max": 325
            },
            {
              "min": 325,
              "max": 330
            }
          ],
          "range_size": {
            "min": 26,
            "max": 32
          },
          "range_colors": [
            "#fde98e",
            "#60c1ff",
            "#fcb08e",
            "#fb88ff",
            "#b4fed4",
            "#cbfaa9",
            "#78d6f8"
          ],
          "display_shadow": true,
          "shadow_color": "#101010",
          "shadow_point": {
            "x": -1,
            "y": -1
          },
          "image_alpha": 1,
          "use_shape_original_color": true
        },
        "thumb": {
          "image_size": {
            "width": 150,
            "height": 40
          },
          "range_verify_length": {
            "min": 2,
            "max": 4
          },
          "disabled_range_verify_length": false,
          "range_text_size": {
            "min": 22,
            "max": 28
          },
          "range_text_colors": [
            "#1f55c4",
            "#780592",
            "#2f6b00",
            "#910000",
            "#864401",
            "#675901",
            "#016e5c"
          ],
          "range_background_colors": [
            "#1f55c4",
            "#780592",
            "#2f6b00",
            "#910000",
            "#864401",
            "#675901",
            "#016e5c"
          ],
          "is_non_deform_ability": false,
          "background_distort": 4,
          "background_distort_alpha": 1,
          "background_circles_num": 24,
          "background_slim_line_num": 2
        }
      },
      "dark_cn_click": {
        "master": {
          "image_size": {
            "width": 300,
            "height": 200
          },
          "range_length": {
            "min": 6,
            "max": 7
          },
          "range_angles": [
            {
              "min": 20,
              "max": 35
            },
            {
              "min": 35,
              "max": 45
            },
            {
              "min": 290,
              "max": 305
            },
            {
              "min": 305,
              "max": 325
            },
            {
              "min": 325,
              "max": 330
            }
          ],
          "range_size": {
            "min": 26,
            "max": 32
          },
          "range_colors": [
            "#fde98e",
            "#60c1ff",
            "#fcb08e",
            "#fb88ff",
            "#b4fed4",
            "#cbfaa9",
            "#78d6f8"
          ],
          "display_shadow": true,
          "shadow_color": "#101010",
          "shadow_point": {
            "x": -1,
            "y": -1
          },
          "image_alpha": 1,
          "use_shape_original_color": true
        },
        "thumb": {
          "image_size": {
            "width": 150,
            "height": 40
          },
          "range_verify_length": {
            "min": 2,
            "max": 4
          },
          "disabled_range_verify_length": false,
          "range_text_size": {
            "min": 22,
            "max": 28
          },
          "range_text_colors": [
            "#1f55c4",
            "#780592",
            "#2f6b00",
            "#910000",
            "#864401",
            "#675901",
            "#016e5c"
          ],
          "range_background_colors": [
            "#1f55c4",
            "#780592",
            "#2f6b00",
            "#910000",
            "#864401",
            "#675901",
            "#016e5c"
          ],
          "is_non_deform_ability": false,
          "background_distort": 4,
          "background_distort_alpha": 1,
          "background_circles_num": 24,
          "background_slim_line_num": 2
        }
      }
    },
    "slide_config_maps": {
      "default_cn_slide": {
        "master": {
          "image_size": {
            "width": 300,
            "height": 200
          },
          "image_alpha": 1
        },
        "thumb": {
          "range_graph_size": {
            "min": 20,
            "max": 100
          },
          "range_graph_angles": [
            {
              "min": 20,
              "max": 100
            }
          ],
          "generate_graph_number": 1,
          "enable_graph_vertical_random": false,
          "range_dead_zone_directions": ["left", "right"]
        }
      }
    },
    "rotate_config_maps": {
      "default_cn_rotate": {
        "master": {
          "image_square_size": 200
        },
        "thumb": {
          "range_angles": [
            {
              "min": 20,
              "max": 100
            }
          ],
          "range_image_square_sizes": [140, 150, 170],
          "image_alpha": 1
        }
      }
    }
  }
}`
	tmpFile, err := os.CreateTemp("", "CaptchaConfig.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString(configContent)
	assert.NoError(t, err)
	tmpFile.Close()

	config, err := Load(tmpFile.Name())
	assert.NoError(t, err)

	err = Validate(config)
	assert.NoError(t, err)

	fmt.Println(config)
}
