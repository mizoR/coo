package usage

import (
	"fmt"
	"os"
	"reflect"
)

func Show(t reflect.Type) {
	os.Stderr.WriteString("Options:\n\n")

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag

		var o string
		if s := tag.Get("short"); s != "" {
			o = fmt.Sprintf("-%s, --%s", tag.Get("short"), tag.Get("long"))
		} else {
			o = fmt.Sprintf("--%s", tag.Get("long"))
		}

		fmt.Fprintf(
			os.Stderr,
			"  %-25s %s\n",
			o,
			tag.Get("description"),
		)
	}
}
