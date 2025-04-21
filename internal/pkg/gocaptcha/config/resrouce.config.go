/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package config

// ResourceChar .
type ResourceChar struct {
	Type      string              `json:"type"`
	Languages map[string][]string `json:"languages"`
}

// ResourceFileConfig .
type ResourceFileConfig struct {
	Type     string            `json:"type"`
	FileDir  string            `json:"file_dir"`
	FileMaps map[string]string `json:"file_maps"`
}

// ResourceMultiFileConfig .
type ResourceMultiFileConfig struct {
	Type       string            `json:"type"`
	FileDir    string            `json:"file_dir"`
	FileMaps   map[string]string `json:"file_maps"`
	FileMaps02 map[string]string `json:"file_maps_02"`
	FileMaps03 map[string]string `json:"file_maps_03"`
}

// ResourceConfig defines the configuration structure for the gocaptcha resource
type ResourceConfig struct {
	Version     string                  `json:"version"`
	Char        ResourceChar            `json:"char"`
	Font        ResourceFileConfig      `json:"font"`
	ShapeImage  ResourceFileConfig      `json:"shapes_image"`
	MasterImage ResourceFileConfig      `json:"master_image"`
	ThumbImage  ResourceFileConfig      `json:"thumb_image"`
	TileImage   ResourceMultiFileConfig `json:"tile_image"`
}
