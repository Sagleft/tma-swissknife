package router

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sagleft/tma-swissknife/helpers"
	"github.com/Sagleft/tma-swissknife/rest"
	"github.com/gin-gonic/gin"
)

func onSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, rest.Success(data))
}

func handleError(ctx *gin.Context, err error) {
	ctx.JSON(
		http.StatusInternalServerError,
		rest.ErrorMessage(err),
	)
}

// path -> hash
func HashAssets(assetsPath string) (map[string]string, error) {
	manifest := make(map[string]string)

	// Рекурсивно обходим директорию
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// пропускаем ошибки доступа
			return nil
		}
		// Игнорируем директории
		if info.IsDir() {
			return nil
		}

		// Открываем файл
		f, ferr := os.Open(path)
		if ferr != nil {
			return nil
		}
		defer f.Close()

		// Читаем весь файл
		b, err := io.ReadAll(f)
		if err != nil {
			return nil
		}

		// Вычисляем MD5 через вашу функцию
		hash := helpers.XxHash64Base32(b)

		// Формируем путь относительно assetsPath
		rel, rerr := filepath.Rel(assetsPath, path)
		if rerr != nil {
			rel = path
		}

		manifest[rel] = hash
		return nil
	}

	// Начинаем обход
	if err := filepath.Walk(assetsPath, walkFn); err != nil {
		return nil, fmt.Errorf("scan %q: %w", assetsPath, err)
	}

	return manifest, nil
}

func ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}
