package feidee

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

//数据导出到文件（随手记WEB版格式的xls文件）
func (cli *Client) ExportToFile(filename string) error {
	downloadAddr, err := cli.GetExportLink()
	resp, err := cli.httpClient.Get(downloadAddr)
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取失败: %s", err)
	}

	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return fmt.Errorf("保存失败: %s", err)
	}

	return nil
}

//获取数据导出的链接（导出为随手记WEB版格式的xls文件）
func (cli *Client) GetExportLink() (string, error) {
	addr := BaseUrl + "/money/data/index.jsp"
	resp, err := cli.httpClient.Get(addr)
	if err != nil {
		return "", fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", fmt.Errorf("读取响应出错: %s", err)
	}

	archors := doc.Find("table.out-data a")
	for i := range archors.Nodes {
		archor := archors.Eq(i)
		linkText := archor.Text()
		if linkText != "web版" {
			continue
		}
		href, found := archor.Attr("href")
		if !found {
			continue
		}

		baseUrl, _ := url.Parse(addr)
		downloadUrl, err := url.Parse(href)
		if err != nil {
			return "", fmt.Errorf("不合法的下载链接:%s", href)
		}

		return baseUrl.ResolveReference(downloadUrl).String(), nil
	}

	return "", fmt.Errorf("未找到符合条件的链接")
}