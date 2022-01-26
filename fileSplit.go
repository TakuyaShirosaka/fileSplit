package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Usage = flagUsage
	input := flag.String("i", "./sample.txt", "inputのファイルパス、ex:D:¥sample.txt")
	splitByte := flag.Int("b", 52428800, "分割バイト数を指定")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	r, err := os.Open(*input)
	b := splitByte
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, *b) // 1. スライス確保
	defer r.Close()

	extname := filepath.Base(r.Name())                 // ファイル拡張子
	ext := filepath.Ext(r.Name())                      // ファイル名+拡張子
	resultBaseName := strings.TrimSuffix(extname, ext) // ファイル名

	// 分割結果の格納先チェック
	if f, err := os.Stat("./result"); os.IsNotExist(err) || !f.IsDir() {
		os.Mkdir("./result", 777)
	}

	for prefix := 1; true; prefix++ {
		readBytesLen, _ := r.Read(buf) // 2. Read実行
		if readBytesLen == 0 {
			fmt.Println("分割完了")
			return
		}

		resultFileName := getResultFileName(resultBaseName, prefix, extname)
		if _, err := writeResultFile(buf, readBytesLen, resultFileName); err != nil {
			return
		}
	}
}

func flagUsage() {
	usageTxt := `ファイル分割ツール
    分割した結果は.exeのディレクトリにresultというディレクトリを生成して格納します。
   -i  inputのファイルパス、ex:D:¥sample.txt", default: 必須
   -b  分割バイト数(byte単位),                 default: 52428800(50MB)
   -h  helpの表示

   使用例:fileSplit.exe -i ./sample.txt
`
	fmt.Fprintf(os.Stderr, "%s\n", usageTxt)
}

func getResultFileName(resultBaseName string, prefix int, extname string) string {
	output := fmt.Sprintf("./result/%s_%d.%s", resultBaseName, prefix, extname)
	return output
}

func writeResultFile(buf []byte, readBytesLen int, output string) (bool, error) {
	if fp, err := os.Create(output); err != nil {
		fmt.Println(err)
		return false, err
	} else {
		fp.WriteString(string(buf[:readBytesLen]))
		fp.Close()
	}
	return true, nil
}
