package file

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"app_server/model"
	"app_server/pkg/idgen"
	"app_server/pkg/ossc"
	"app_server/service/auth"
	"app_server/service/file"
)

func WxFileUpload(c *gin.Context) {
	// 0. 获取用户ID（从认证信息中）
	authHeader := c.GetHeader("Authorization")
	var userID uint
	var err error
	if authHeader != "" {
		userID, err = auth.ParseUserID(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "认证失败: " + err.Error(),
			})
			return
		}
	}

	// 1. 获取微信 wx.uploadFile 上传的文件
	uploadFile, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "获取上传文件失败: " + err.Error(),
		})
		return
	}
	defer uploadFile.Close()

	// 获取文件信息
	filename := header.Filename
	fileSize := header.Size
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 获取用途类型（默认为临时上传）
	usageType := c.DefaultQuery("usage_type", model.UsageTypeTempUpload)

	// 2. 计算文件哈希（如果有用户ID）
	var fileHash string
	var fileData []byte

	if userID > 0 {
		// 读取文件内容到缓冲区
		buf := &bytes.Buffer{}
		teeReader := io.TeeReader(uploadFile, buf)

		// 计算哈希
		fileHash, err = file.CalculateFileHash(teeReader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "计算文件哈希失败: " + err.Error(),
			})
			return
		}

		// 保存文件数据供后续上传
		fileData = buf.Bytes()

		// 查找是否已存在相同文件
		fileService := file.NewService()
		existingFile, err := fileService.GetFileByHash(userID, fileHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "查询已有文件失败: " + err.Error(),
			})
			return
		}

		// 如果文件已存在，直接返回
		if existingFile != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "文件已存在",
				"data": WxFileUploadResponse{
					ID:        strconv.FormatUint(uint64(existingFile.ID), 10),
					URL:       existingFile.OssKey,
					PublicURL: existingFile.PublicURL,
					Filename:  existingFile.OriginalName,
				},
			})
			return
		}
	} else {
		// 没有用户ID时，读取所有数据
		fileData, err = io.ReadAll(uploadFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "读取文件内容失败: " + err.Error(),
			})
			return
		}
	}

	// 生成唯一文件名
	uniqueFilename := idgen.Base36() + ext

	// 3. 流式上传文件到 aliyun oss
	// 指定上传的目录路径
	var objectKey string
	if userID > 0 && usageType != model.UsageTypeTempUpload {
		// 用户文件直接存储到用户目录
		objectKey = fmt.Sprintf("user/%d/%s/%s/%s", userID, usageType,
			time.Now().Format("2006/01/02"), uniqueFilename)
	} else {
		// 临时文件或无用户ID的文件存储到临时目录
		objectKey = fmt.Sprintf("uploads/%s/%s", time.Now().Format("2006/01/02"), uniqueFilename)
	}

	// 上传文件
	err = ossc.Get().UserFileBucket().PutObject(objectKey, bytes.NewReader(fileData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "文件上传至OSS失败: " + err.Error(),
		})
		return
	}

	// 4. 创建文件记录（如果有用户ID）
	var fileID uint
	var publicURL string

	if userID > 0 {
		fileService := file.NewService()
		userFile, err := fileService.CreateFileRecord(
			userID,
			filename,
			fileSize,
			contentType,
			ext,
			objectKey,
			fileHash,
			usageType,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "创建文件记录失败: " + err.Error(),
			})
			return
		}
		fileID = userFile.ID
		publicURL = userFile.PublicURL
	} else {
		// 没有用户ID时，直接生成签名URL
		publicURL, err = ossc.GetPublic().UserFileBucket().SignURL(objectKey, "GET", 10*365*24*60*60)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "生成签名URL失败: " + err.Error(),
			})
			return
		}
	}

	responseData := WxFileUploadResponse{
		URL:       objectKey,
		PublicURL: publicURL,
		Filename:  filename,
	}

	if fileID > 0 {
		responseData.ID = strconv.FormatUint(uint64(fileID), 10)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "上传成功",
		"data": responseData,
	})
}

type WxFileUploadResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	PublicURL string `json:"publicUrl"`
	Filename  string `json:"filename"`
}
