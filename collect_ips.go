package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// 采集源配置
var urls = []struct {
	URL      string
	Tag      string // 提取标签：tr 或 li
}{
	{"https://cf.090227.xyz/ct?ips=6", "tr"},
	{"https://ip.164746.xyz", "tr"},
	{"https://cf.090227.xyz/cu", "tr"},
	{"https://cf.090227.xyz/cmcc?ips=8", "tr"},
	{"https://ipdb.api.030101.xyz/?type=cfv4;proxy", "tr"},
	{"https://ipdb.api.030101.xyz/?type=bestproxy&country=true", "tr"},
	{"https://www.wetest.vip/page/cloudflare/address_v4.html", "tr"},
    {"https://stock.hostmonit.com/CloudFlareYes", "tr"},
    {"https://api.uouin.com/cloudflare.html", "tr"},
    {"https://vps789.com/cfip", "tr"},
    {"https://www.byoip.top/", "tr"},
    {"https://mrxn.net/Vercel.html", "tr"},
    {"https://addressesapi.090227.xyz/CloudFlareYes", "tr"},
    {"https://www.cloudflare.com/ips-v4", "tr"},
}

// IPv4 正则表达式
var ipPattern = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

func main() {
	// 删除旧的 ip.txt
	if _, err := os.Stat("ip.txt"); err == nil {
		os.Remove("ip.txt")
	}

	// 创建 HTTP 客户端（支持 TLS 1.2+）
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	// 用于去重
	seen := make(map[string]bool)

	// 打开输出文件
	file, err := os.OpenFile("ip.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 无法创建 ip.txt: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, src := range urls {
		fmt.Printf("[采集] 正在请求: %s\n", src.URL)

		req, err := http.NewRequest("GET", src.URL, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[错误] 创建请求失败 %s: %v\n", src.URL, err)
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[错误] 请求失败 %s: %v\n", src.URL, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[错误] 读取响应失败 %s: %v\n", src.URL, err)
			continue
		}

		// 解析 HTML
		ips := extractIPs(string(body), src.Tag)
		for _, ip := range ips {
			if !seen[ip] {
				seen[ip] = true
				file.WriteString(ip + "\n")
			}
		}
		fmt.Printf("[完成] %s 提取到 %d 个新IP\n", src.URL, len(ips))
	}

	fmt.Println("✅ 所有 IP 地址已保存到 ip.txt 文件中。")
}

// extractIPs 从 HTML 中提取 IP 地址
func extractIPs(htmlContent, tag string) []string {
	var ips []string

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return ips
	}

	// 递归遍历节点
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			text := getTextContent(n)
			matches := ipPattern.FindAllString(text, -1)
			ips = append(ips, matches...)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	return ips
}

// getTextContent 获取节点内所有文本内容
func getTextContent(n *html.Node) string {
	var sb strings.Builder
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	return sb.String()
}
