package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/karrick/godirwalk"
	"github.com/russross/blackfriday/v2"
)

func main() {
	// Define flags for the directory to scan and max character count.
	var dir string
	var maxCharCount int
	flag.StringVar(&dir, "dir", ".", "指定扫描的根目录")
	flag.IntVar(&maxCharCount, "maxCharCount", 800000, "合并成pdf的最大字符数")
	flag.Parse()

	// Check if wkhtmltopdf is installed.
	_, err := exec.LookPath("wkhtmltopdf")
	if err != nil {
		fmt.Println("未安装wkhtmltopdf，请安装后执行，您可以在这里下载安装：https://wkhtmltopdf.org/downloads.html")
		return
	}

	// Check if the .temp directory exists.
	if _, err := os.Stat(".temp"); !os.IsNotExist(err) {
		// If it exists, delete it.
		err = os.RemoveAll(".temp")
		if err != nil {
			panic(err)
		}
	}

	// Create the .temp directory.
	err = os.Mkdir(".temp", 0755)
	if err != nil {
		panic(err)
	}

	// Step 1: Traverse the specified directory and its subdirectories to find all .md and .mdx files.
	var files []string
	err = godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if strings.HasSuffix(osPathname, ".md") || strings.HasSuffix(osPathname, ".mdx") {
				files = append(files, osPathname)
			}
			return nil
		},
		Unsorted: true, // set true for faster walk, but files are unsorted
	})
	if err != nil {
		panic(err)
	}

	// Step 2: Merge the found files into one .md file and convert to .pdf file if characters exceed maxCharCount.
	var content []byte
	var pdfCount int
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		// Check if appending the current file's content will exceed the max character count.
		if len(content)+len(data) > maxCharCount {
			// Write the content to a file and convert it to a PDF file.
			err = generatePDF(pdfCount, content)
			if err != nil {
				panic(err)
			}

			// Reset the content and increase the PDF count.
			content = []byte{}
			pdfCount++
		}

		// Append the current file's content to the content.
		content = append(content, data...)
		content = append(content, '\n')
	}

	// Convert remaining content to a PDF file.
	if len(content) > 0 {
		err = generatePDF(pdfCount, content)
		if err != nil {
			panic(err)
		}
	}

	// Delete the .temp directory.
	err = os.RemoveAll(".temp")
	if err != nil {
		panic(err)
	}
}

// This function writes content to a .md file, converts it to a .html file, and then converts it to a .pdf file.
func generatePDF(pdfCount int, content []byte) error {
	// Write the content to a .md file.
	err := ioutil.WriteFile(filepath.Join(".temp", "merged.md"), content, 0644)
	if err != nil {
		return err
	}

	// Convert the .md file to a .html file.
	mdContent, err := ioutil.ReadFile(filepath.Join(".temp", "merged.md"))
	if err != nil {
		return err
	}
	htmlContent := blackfriday.Run(mdContent)

	// Parse the HTML content with goquery.
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return err
	}

	// Modify the HTML content.
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// Remove all img tags.
		s.Remove()
	})

	// Get the modified HTML content.
	modifiedHtmlContent, err := doc.Html()
	if err != nil {
		return err
	}

	// Add a style tag to specify a font that supports Chinese.
	modifiedHtmlContent = fmt.Sprintf(`
  <html>
  <head>
   <meta charset="utf-8">
   <style>
    body {
     font-family: "Microsoft YaHei","SimSun";
    }
   </style>
  </head>
  <body>
   %s
  </body>
  </html>
 `, modifiedHtmlContent)

	// Write the modified HTML content to a file.
	err = ioutil.WriteFile(filepath.Join(".temp", "merged.html"), []byte(modifiedHtmlContent), 0644)
	if err != nil {
		return err
	}

	// Convert the HTML file to a PDF file with wkhtmltopdf.
	pdfFilename := fmt.Sprintf("merged%d.pdf", pdfCount)
	cmd := exec.Command("wkhtmltopdf", filepath.Join(".temp", "merged.html"), pdfFilename)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
