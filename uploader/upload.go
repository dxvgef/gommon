package uploader

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Uploader 上传实例及参数
type Uploader struct {
	Request        *http.Request
	FieldName      string      // 上传控件的name值
	MaxSize        int64       // 文件大小限制（KB）
	AllowMIME      []string    // 允许上传的文件MIME值
	SaveName       string      // 存储文件名（不含后缀名），留空则保存原文件名
	SaveRootPath   string      // 存储根路径（绝对路径）
	SaveSubPath    string      // 存储子路径（相对SaveRootPath）
	SaveSuffix     string      // 存储文件的后缀名（如果指定了此属性值，则强制更换后缀名）
	DirPermission  os.FileMode // 文件存放目录权限，如果目录已存在，则此参数无效
	FilePermission os.FileMode // 文件权限
}

// Result 上传结果
type Result struct {
	FileSize   int64  // 文件大小
	FileMIME   string // 文件的MIME值
	FileName   string // 上传后的文件完整路径
	FileSuffix string // 上传后的文件后缀名
}

// _sizeInterface 文件大小
type _sizeInterface interface {
	Size() int64
}
type _statInterface interface {
	Stat() (os.FileInfo, error)
}

// Exec 执行上传
func (obj *Uploader) Exec() (result Result, err error) {
	var multipartFileCloseErr, newFileCloseErr error

	// 获得上传文件的数据
	multipartFile, head, err := obj.Request.FormFile(obj.FieldName)
	if err != nil {
		return
	}

	// 获得文件大小
	if statInterface, ok := multipartFile.(_statInterface); ok {
		fileInfo, statErr := statInterface.Stat()
		if statErr != nil {
			err = statErr
			if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
				err = multipartFileCloseErr
			}
			return
		}
		result.FileSize = fileInfo.Size()
	}
	if result.FileSize == 0 {
		if sizeInterfaceTmp, okTmp := multipartFile.(_sizeInterface); okTmp {
			result.FileSize = sizeInterfaceTmp.Size()
		}
	}

	// 判断文件大小
	if result.FileSize == 0 {
		err = errors.New("文件大小为0")
		if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
			err = multipartFileCloseErr
		}
		return
	}
	if result.FileSize > obj.MaxSize*1024 {
		err = errors.New("文件大小超出限制")
		if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
			err = multipartFileCloseErr
		}
		return
	}

	// 判断文件MIME值
	result.FileMIME = head.Header.Get("Content-Type")
	if !inStr(obj.AllowMIME, result.FileMIME) {
		err = errors.New("不允许上传该类型的文件")
		if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
			err = multipartFileCloseErr
		}
		return
	}

	if obj.SaveSuffix != "" {
		// 使用指定的后缀名
		result.FileSuffix = obj.SaveSuffix
	} else {
		// 获得原始文件的后缀名
		pos := strings.LastIndex(head.Filename, ".")
		if pos != 0 {
			result.FileSuffix = head.Filename[pos+1:]
		}
	}

	// 如果文件名没有指定,则使用原始文件名
	if obj.SaveName == "" {
		obj.SaveName = head.Filename
	}

	// 递归创建目录
	err = os.MkdirAll(obj.SaveRootPath+"/"+obj.SaveSubPath, obj.DirPermission)
	if err != nil {
		err = errors.New("创建目录失败")
		if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
			err = multipartFileCloseErr
		}
		return
	}
	// nolint:gosec
	// 在指定的路径创建文件
	newFile, err := os.OpenFile(filepath.Clean(obj.SaveRootPath+"/"+obj.SaveSubPath+"/"+obj.SaveName+"."+result.FileSuffix), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, obj.FilePermission)
	if err != nil {
		err = errors.New("创建文件失败")
		if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
			err = multipartFileCloseErr
		}
		return
	}

	// 复制数据到文件
	if _, err = io.Copy(newFile, multipartFile); err != nil {
		err = errors.New("复制数据到文件失败")
		if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
			err = multipartFileCloseErr
		}
		if newFileCloseErr = newFile.Close(); newFileCloseErr != nil {
			err = newFileCloseErr
		}
		return
	}

	// 返回文件路径
	// result.FileName = obj.SaveSubPath + "/" + obj.SaveName + fileExt
	result.FileName = obj.SaveName + "." + result.FileSuffix

	if multipartFileCloseErr = multipartFile.Close(); multipartFileCloseErr != nil {
		err = multipartFileCloseErr
	}
	if newFileCloseErr = newFile.Close(); newFileCloseErr != nil {
		err = newFileCloseErr
	}
	return
}

// 不变化的slice可以先sort.Strings(sk)再range速度快一倍
// 检查string值在一个string slice中是否存在
func inStr(s []string, str string) bool {
	for k := range s {
		if str == s[k] {
			return true
		}
	}
	return false
}
