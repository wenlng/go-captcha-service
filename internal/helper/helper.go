package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/wenlng/go-captcha/v2/base/codec"
)

// GenUniqueId .
func GenUniqueId() (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

// Marshal .
func Marshal(data interface{}) interface{} {
	typeof := reflect.TypeOf(data)
	valueof := reflect.ValueOf(data)

	for i := 0; i < typeof.Elem().NumField(); i++ {
		if valueof.Elem().Field(i).IsZero() {
			def := typeof.Elem().Field(i).Tag.Get("default")
			if def != "" {
				switch typeof.Elem().Field(i).Type.String() {
				case "int":
					result, _ := strconv.Atoi(def)
					valueof.Elem().Field(i).SetInt(int64(result))
				case "uint":
					result, _ := strconv.ParseUint(def, 10, 64)
					valueof.Elem().Field(i).SetUint(result)
				case "string":
					valueof.Elem().Field(i).SetString(def)
				case "interface {}":
					valueof.Elem().Field(i).SetZero()
				}
			}
		}
	}
	return data
}

// MarshalJson .
func MarshalJson(data interface{}) ([]byte, error) {
	typeof := reflect.TypeOf(data)
	valueof := reflect.ValueOf(data)

	for i := 0; i < typeof.Elem().NumField(); i++ {
		if valueof.Elem().Field(i).IsZero() {
			def := typeof.Elem().Field(i).Tag.Get("default")
			if def != "" {
				switch typeof.Elem().Field(i).Type.String() {
				case "int":
					result, _ := strconv.Atoi(def)
					valueof.Elem().Field(i).SetInt(int64(result))
				case "uint":
					result, _ := strconv.ParseUint(def, 10, 64)
					valueof.Elem().Field(i).SetUint(result)
				case "string":
					valueof.Elem().Field(i).SetString(def)
				}
			}
		}
	}
	return json.Marshal(data)
}

// GetPWD .
func GetPWD() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

// ReadFileStream reads file contents using streaming, suitable for large files
func ReadFileStream(filePath string) ([]byte, error) {
	cleanPath := filepath.Clean(filePath)

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file %s does not exist: %w", cleanPath, err)
		}
		return nil, fmt.Errorf("cannot access file %s: %w", cleanPath, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a file", cleanPath)
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", cleanPath, err)
	}
	defer file.Close()

	var buffer bytes.Buffer
	const chunkSize = 8192 // 8KB buffer size
	chunk := make([]byte, chunkSize)
	for {
		n, err := file.Read(chunk)
		if n > 0 {
			buffer.Write(chunk[:n])
		}
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			return nil, fmt.Errorf("failed to read file %s: %w", cleanPath, err)
		}
	}
	return buffer.Bytes(), nil
}

// LoadImageData .
func LoadImageData(filepath string) (image.Image, error) {
	stream, err := ReadFileStream(filepath)
	if err != nil {
		return nil, err
	}
	if path.Ext(filepath) == ".png" {
		png, err := codec.DecodeByteToPng(stream)
		if err != nil {
			return nil, err
		}
		return png, nil
	}

	jpeg, err := codec.DecodeByteToJpeg(stream)
	if err != nil {
		return nil, err
	}
	return jpeg, nil
}

// DeepEqual depth compares whether two values are exactly the same
func DeepEqual(a, b interface{}) bool {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	if v1.Type() != v2.Type() {
		return false
	}

	return deepEqualValue(v1, v2)
}

// deepEqualValue recursively compares two reflection values
func deepEqualValue(v1, v2 reflect.Value) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}

	switch v1.Kind() {
	case reflect.Ptr:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		return deepEqualValue(v1.Elem(), v2.Elem())

	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			if !deepEqualValue(v1.Field(i), v2.Field(i)) {
				return false
			}
		}
		return true

	case reflect.Slice, reflect.Array:
		if v1.Len() != v2.Len() {
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepEqualValue(v1.Index(i), v2.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Map:
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		for _, k := range v1.MapKeys() {
			if !deepEqualValue(v1.MapIndex(k), v2.MapIndex(k)) {
				return false
			}
		}
		return true

	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return deepEqualValue(v1.Elem(), v2.Elem())

	default:
		return reflect.DeepEqual(v1.Interface(), v2.Interface())
	}
}

// FileExists .
func FileExists(filepathStr string) bool {
	p := filepath.Clean(filepathStr)
	_, err := os.Stat(p)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// IsFile .
func IsFile(filepathStr string) bool {
	p := filepath.Clean(filepathStr)
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DeleteFile .
func DeleteFile(path string) error {
	cleanPath := filepath.Clean(path)

	if cleanPath == "" || cleanPath == "." || cleanPath == "/" || cleanPath == string(os.PathSeparator) {
		return os.ErrInvalid
	}

	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}
	if fileInfo.IsDir() {
		return os.ErrInvalid
	}

	err = os.Remove(cleanPath)
	if err != nil {
		return err
	}

	return nil
}

// EnsureDir .
func EnsureDir(path string) error {
	return EnsureDirWithPerm(path, 0755)
}

// EnsureDirWithPerm make sure all directories in the specified path exist. If they do not exist, create them
func EnsureDirWithPerm(path string, perm os.FileMode) error {
	cleanPath := filepath.Clean(path)
	if cleanPath == "" || cleanPath == "." || cleanPath == "/" || cleanPath == string(os.PathSeparator) {
		return os.ErrInvalid
	}
	fileInfo, err := os.Stat(cleanPath)
	if err == nil {
		if fileInfo.IsDir() {
			return nil
		}
		return os.ErrExist
	}
	if !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(cleanPath, perm)
	if err != nil {
		return err
	}

	return nil
}

// IsValidDirName verify that the directory name is valid
// (only letters, numbers, hyphens, underscores allowed)
func IsValidDirName(dir string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(dir)
}

// IsSubPath check whether the target path is in the root directory
func IsSubPath(path, root string) bool {
	path = filepath.Clean(path)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absRoot, absPath)
}

// TraverseDir recursively traverse all file paths in the directory
// and return the relative paths relative to the specified parent directory
func TraverseDir(root, baseDir string) ([]string, error) {
	cleanRoot := filepath.Clean(root)
	cleanBaseDir := filepath.Clean(baseDir)

	if cleanRoot == "" || cleanRoot == "/" || cleanRoot == string(os.PathSeparator) {
		return nil, os.ErrInvalid
	}
	if cleanBaseDir == "" || cleanBaseDir == "/" || cleanBaseDir == string(os.PathSeparator) {
		return nil, os.ErrInvalid
	}

	absRoot, err := filepath.Abs(cleanRoot)
	if err != nil {
		return nil, err
	}
	absBaseDir, err := filepath.Abs(cleanBaseDir)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(absRoot)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, os.ErrInvalid
	}

	if _, err := os.Stat(absBaseDir); err != nil {
		return nil, err
	}

	var filePaths []string

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(absBaseDir, path)
			if err != nil {
				return err
			}
			filePaths = append(filePaths, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filePaths, nil
}
