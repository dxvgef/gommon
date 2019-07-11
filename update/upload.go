// 用于net/http标准包的文件上传

package update

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

// Uploader 上传实例及参数
type Uploader struct {
	Request      *http.Request
	FieldName    string   // 上传控件的name值
	MaxSize      int64    // 文件大小限制（KB）
	AllowMIME    []string // 允许上传的文件MIME值
	SaveName     string   // 存储文件名（不含后缀名），留空则保存原文件名
	SaveRootPath string   // 存储根路径（绝对路径）
	SaveSubPath  string   // 存储子路径（相对SaveRootPath）
	SaveSuffix   string   // 存储文件的后缀名（如果指定了此属性值，则强制更换后缀名）
}

// Result 上传结果
type Result struct {
	FileSize   int64  // 文件大小
	FileMIME   string // 文件的MIME值
	FileName   string // 上传后的文件完整路径
	FileSuffix string // 上传后的文件后缀名
}

// sizeInterface 文件大小
type sizeInterface interface {
	Size() int64
}
type statInterface interface {
	Stat() (os.FileInfo, error)
}

// Exec 执行上传
func (obj *Uploader) Exec() (result Result, err error) {
	// 获得上传文件的数据
	file, head, err := obj.Request.FormFile(obj.FieldName)
	if err != nil {
		return
	}
	defer file.Close()

	// 获得文件大小
	if statInterface, ok := file.(statInterface); ok {
		fileInfo, _ := statInterface.Stat()
		result.FileSize = fileInfo.Size()
	}
	if result.FileSize == 0 {
		if sizeInterface, ok := file.(sizeInterface); ok {
			result.FileSize = sizeInterface.Size()
		}
	}

	// 判断文件大小
	if result.FileSize == 0 {
		err = errors.New("文件大小为0")
		return
	}
	if result.FileSize > obj.MaxSize*1024 {
		err = errors.New("文件大小超出限制")
		return
	}

	// 判断文件MIME值
	result.FileMIME = head.Header.Get("Content-Type")
	if inStr(obj.AllowMIME, result.FileMIME) == false {
		err = errors.New("不允许上传该类型的文件")
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
	err = os.MkdirAll(obj.SaveRootPath+"/"+obj.SaveSubPath, 0755)
	if err != nil {
		err = errors.New("创建目录失败")
		return
	}
	// 在指定的路径创建文件
	f, err := os.OpenFile(obj.SaveRootPath+"/"+obj.SaveSubPath+"/"+obj.SaveName+"."+result.FileSuffix, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		err = errors.New("创建文件失败")
		return
	}
	defer f.Close()

	// 复制数据到文件
	if _, err = io.Copy(f, file); err != nil {
		err = errors.New("复制数据到文件失败")
		return
	}

	// 返回文件路径
	// result.FileName = obj.SaveSubPath + "/" + obj.SaveName + fileExt
	result.FileName = obj.SaveName + "." + result.FileSuffix

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
