package uploader

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type (
	// Uploader 上传实例及参数
	Uploader struct {
		DirPermission  os.FileMode // 文件存放目录权限，如果目录已存在，则此参数无效
		FilePermission os.FileMode // 文件权限
		MaxSize        int64       // 文件大小限制（KB）
		FieldName      string      // 上传控件的name值
		SaveName       string      // 存储文件名（不含后缀名），留空则保存原文件名
		SaveRootPath   string      // 存储根路径（绝对路径）
		SaveSubPath    string      // 存储子路径（相对SaveRootPath）
		SaveSuffix     string      // 存储文件的后缀名（如果指定了此属性值，则强制更换后缀名）
		AllowMIME      []string    // 允许上传的文件MIME值
		Request        *http.Request
	}
	// 错误类型
	Error struct {
		Status        int    // HTTP状态码
		OriginalError error  // 原始的错误
		FriendlyText  string // 错误文本
	}
	// Result 上传结果
	Result struct {
		FileSize   int64  // 文件大小
		FileMIME   string // 文件的MIME值
		FileName   string // 上传后的文件完整路径
		FileSuffix string // 上传后的文件后缀名
	}
	// _sizeInterface 文件大小
	_sizeInterface interface {
		Size() int64
	}
	_statInterface interface {
		Stat() (os.FileInfo, error)
	}
)

// Exec 执行上传
func (obj *Uploader) Exec() (result Result, execErr Error) {
	var (
		statInterface _statInterface
		sizeInterface _sizeInterface
		fileInfo      os.FileInfo
		newFile       *os.File
		ok            bool
	)

	// 获得上传文件的数据
	multipartFile, head, err := obj.Request.FormFile(obj.FieldName)
	if err != nil {
		execErr.FriendlyText = "无法获得要上传的文件数据"
		execErr.OriginalError = err
		execErr.Status = 400
		return
	}
	defer func() {
		if err = multipartFile.Close(); err != nil {

		}
	}()

	// 获得文件大小
	if statInterface, ok = multipartFile.(_statInterface); ok {
		fileInfo, err = statInterface.Stat()
		if err != nil {
			execErr.FriendlyText = "无法获得要上传的文件大小"
			execErr.OriginalError = err
			execErr.Status = 400
			return
		}
		result.FileSize = fileInfo.Size()
	}
	if result.FileSize == 0 {
		if sizeInterface, ok = multipartFile.(_sizeInterface); ok {
			result.FileSize = sizeInterface.Size()
		}
	}

	// 判断文件大小
	if result.FileSize == 0 {
		err = errors.New("文件大小为0")
		execErr.FriendlyText = err.Error()
		execErr.OriginalError = err
		execErr.Status = 400
		return
	}

	if result.FileSize > obj.MaxSize*1024 {
		execErr.FriendlyText = "文件大小超出限制"
		execErr.OriginalError = errors.New("文件大小(" + strconv.FormatInt(result.FileSize, 10) + ")超出限制(" + strconv.FormatInt(obj.MaxSize, 10) + ")")
		execErr.Status = 400
		return
	}

	// 判断文件MIME值
	result.FileMIME = head.Header.Get("Content-Type")
	if !inStr(obj.AllowMIME, result.FileMIME) {
		execErr.FriendlyText = "不允许上传该类型的文件"
		execErr.OriginalError = errors.New("不允许上传" + result.FileMIME + "类型的文件")
		execErr.Status = 400
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
	savePath := filepath.Clean(obj.SaveRootPath + "/" + obj.SaveSubPath)
	err = os.MkdirAll(savePath, obj.DirPermission)
	if err != nil {
		execErr.FriendlyText = "创建目录失败"
		execErr.OriginalError = errors.New("创建目录失败 " + savePath)
		execErr.Status = 500
		return
	}
	// 在指定的路径创建文件
	filePath := filepath.Clean(obj.SaveRootPath + "/" + obj.SaveSubPath + "/" + obj.SaveName + "." + result.FileSuffix)
	newFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, obj.FilePermission)
	if err != nil {
		execErr.FriendlyText = "上传文件失败"
		execErr.OriginalError = errors.New("创建文件失败 " + filePath)
		execErr.Status = 500
		return
	}
	defer func() {
		if err = newFile.Close(); err != nil {
		}
	}()

	// 复制数据到文件
	if _, err = io.Copy(newFile, multipartFile); err != nil {
		execErr.FriendlyText = "上传文件失败"
		execErr.OriginalError = errors.New("复制上传数据到文件失败 " + filePath)
		execErr.Status = 500
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
