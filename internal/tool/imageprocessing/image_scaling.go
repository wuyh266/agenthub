package imageprocessing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ImageScaling 是一个图像缩放工具。它实现了 Tool 接口，提供图像缩放功能。
type ImageScaling struct {
	endpoint string
	client   *http.Client
}

// NewImageScaling 创建一个新的 ImageScaling 实例，使用endpoint。设置了30秒的超时时间。
func NewImageScaling(endpoint string) *ImageScaling {
	return &ImageScaling{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (i *ImageScaling) Name() string {
	return "image_scaling"
}

func (i *ImageScaling) Description() string {
	return "- 功能：把本地 jpeg/png 图片等比缩放到指定宽度\n- input 格式：一行，空格分隔两个字段——图片路径 目标宽度，比如(./picture/window.png 256)\n- 约束：宽度 10-4096;仅支持 jpeg/png\n- 输出：缩放后图片的保存路径"
}

func (i *ImageScaling) Execute(ctx context.Context, input string) (string, error) {
	// Execute 解析 input 得到图片路径与目标宽度，判定 MIME 类型，打开文件，构建 multipart 请求体存入 buf，等 B-2b 发送。

	parts := strings.Fields(input) // 这个Fields函数会把输入的字符串按空格分割成若干个字段，返回一个字符串切片。比如输入"./picture/window.png 256"，就会得到一个长度为2的切片，分别是"./picture/window.png"和"256"。跟TrimSpace不同，TrimSpace只是去掉首尾空格，不会分割字符串。这里用Fields是为了方便后续获取图片路径和宽度两个参数。
	if len(parts) != 2 {
		return "", fmt.Errorf("输入格式错误，正确格式为：图片路径 目标宽度，例如：./picture/window.png 256")
	}
	imgPath := parts[0]
	widthStr := parts[1]
	width, err := strconv.Atoi(widthStr) // Atoi 是字符串转整数的函数，如果转换失败，err 就不为 nil
	if err != nil {
		return "", fmt.Errorf("宽度必须是数字")
	}
	if width < 10 || width > 4096 {
		return "", fmt.Errorf("宽度必须在 10-4096 范围内")
	}
	var mimeType string                           // 这里这个mime类型是为了后续的multipart/form-data请求中使用的。multipart/form-data是一种常用的HTTP请求格式，正常的比如我们计网里面的数据传输是单独传输字节流，但是我们这里要同时上传图像和width两个参数，他就需要用到这个mime/multipart/form-data格式，这个格式用一串随机生成的分隔符（boundary）来分隔不同的字段，mimeType就是告诉服务端这个文件的类型是什么，服务端就知道怎么处理这个文件了。
	ext := strings.ToLower(filepath.Ext(imgPath)) // filepath.Ext(imgPath) 会返回文件的扩展名，比如".png"、".jpg"。然后用 strings.ToLower 转成小写，方便后续判断。小写的原因是有些时候后缀是大写的.PNG
	if ext == ".jpg" || ext == ".jpeg" {
		mimeType = "image/jpeg"
	} else if ext == ".png" {
		mimeType = "image/png"
	} else {
		return "", fmt.Errorf("仅支持 jpeg/png")
	}
	file, err := os.Open(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	// head 是 map[string][]string 类型，Set 方法由 textproto 包提供；这里设置的头会成为文件 part 的头部
	head := textproto.MIMEHeader{}
	head.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, filepath.Base(imgPath))) //filepath.Base 把 "./picture/window.png" 截成 "window.png"，只留文件名，不把本地目录结构漏给服务端
	head.Set("Content-Type", mimeType)
	pw, err := mw.CreatePart(head)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(pw, file) //把图片字节从 file 灌进 part。io.Copy(dst, src) 循环搬运直到 src 读完，返回（写入字节数, err）
	if err != nil {
		return "", fmt.Errorf("复制图片数据到请求体失败: %w", err)
	}
	err = mw.WriteField("width", strconv.Itoa(width))
	if err != nil {
		return "", fmt.Errorf("写入 width 字段失败: %w", err)
	}
	err = mw.Close()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.endpoint, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := i.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("请求失败，状态码: %d, 响应体: %s", resp.StatusCode, string(bodyBytes))
	}
	// 读取响应体：缩放后的图片二进制字节具体长什么样详见 imageprocess项目
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}
	outPath := strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + "_resized" + ext
	err = os.WriteFile(outPath, imageData, 0644)
	if err != nil {
		return "", fmt.Errorf("保存缩放后图片失败: %w", err)
	}
	return outPath, nil
}
