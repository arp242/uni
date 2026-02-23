//go:build go_run_only

package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"zgo.at/zli"
)

func main() {
	os.Setenv("README", "yes")
	tpl := template.New("").Funcs(template.FuncMap{
		"trim": func(n int, lines string) string {
			s := strings.SplitN(lines, "\n", n+2)
			return strings.Join(s[:len(s)-1], "\n") + "\n    …"
		},
		"example": func(args ...string) (string, error) {
			o, err := exec.Command("uni", args...).CombinedOutput()
			if err != nil {
				return "", err
			}

			for i := range args {
				if strings.Index(args[i], " ") > -1 {
					args[i] = "'" + args[i] + "'"
				}
			}

			out := "    % uni " + strings.Join(args, " ") + "\n"
			for _, line := range bytes.Split(bytes.TrimRight(o, "\n"), []byte{'\n'}) {
				out += "    " + string(line) + "\n"
			}
			return out[:len(out)-1], nil
		},
	})

	tpl, err := tpl.ParseFiles(".readme.gotxt")
	zli.F(err)

	fp, err := os.Create("README.md")
	zli.F(err)

	zli.F(tpl.ExecuteTemplate(fp, ".readme.gotxt", nil))
}
