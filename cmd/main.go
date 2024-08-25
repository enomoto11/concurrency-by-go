// go:build ignore
package main

import (
	"fmt"
	"os"
	"text/template"
)

func main() {
	chapterNum := os.Args[1]
	if chapterNum == "" {
		fmt.Println("Please provide a chapter number.")
		os.Exit(1)
	}

	if err := MakeDir(chapterNum); err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	if err := Touch(chapterNum); err != nil {
		fmt.Println("Error creating files:", err)
		os.Exit(1)
	}
}

func MakeDir(chapterNum string) error {
	fileInfo, err := os.Lstat("./")
	if err != nil {
		return err
	}

	fileMode := fileInfo.Mode()
	unixPerms := fileMode & os.ModePerm

	if err = os.Mkdir(fmt.Sprintf("../chap%s/", chapterNum), unixPerms); err != nil {
		return err
	}

	return nil
}

func Touch(chapterNum string) error {
	if err := touch_type(chapterNum); err != nil {
		return err
	}

	if err := touch_chap_subchapter(chapterNum); err != nil {
		return err
	}

	return nil
}

func touch_type(chapterNum string) error {
	t, err := template.ParseFiles("../templates/type.tmpl")
	if err != nil {
		return err
	}

	fp, err := os.Create("../chap" + chapterNum + "/type.go")
	if err != nil {
		return err
	}
	defer fp.Close()

	if err = t.Execute(fp, map[string]string{"Chapter": chapterNum}); err != nil {
		return err
	}

	return nil
}

func touch_chap_subchapter(chapterNum string) error {
	t, err := template.ParseFiles("../templates/chap.subchap.tmpl")
	if err != nil {
		return err
	}
	fp, err := os.Create("../chap" + chapterNum + "/" + chapterNum + ".1.go")
	if err != nil {
		return err
	}
	defer fp.Close()

	if err = t.Execute(fp, map[string]string{"Chapter": chapterNum}); err != nil {
		return err
	}

	return nil
}
