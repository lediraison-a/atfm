package app

import (
	"atfm/app/config"
	"atfm/app/icons"
	"atfm/app/models"
	"fmt"
	"regexp"
	"strings"
)

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func RenderFileInfo(item models.FileInfo, fileInfoFormat []string, conf config.DisplayConfig) string {
	t := ""
	tmode := item.Mode.String()
	tdate := item.ModTime.Format(conf.DateFormat)
	tsize := ByteCountSI(item.Size)
	if item.IsDir {
		tsize = ""
	}
	tname := item.Name
	ic, _ := icons.GetFileIcon(item.Name, conf.Theme.Text_default)
	for i, v := range fileInfoFormat {
		vt := v
		if strings.Contains(vt, `{mod}`) && tmode != "" {
			re := regexp.MustCompile(`{mod}`)
			vt = string(re.ReplaceAll([]byte(vt), []byte(tmode)))
		} else if strings.Contains(vt, `{date}`) && tdate != "" {
			re := regexp.MustCompile(`{date}`)
			vt = string(re.ReplaceAll([]byte(vt), []byte(tdate)))
		} else if strings.Contains(vt, `{size}`) && tsize != "" {
			re := regexp.MustCompile(`{size}`)
			vt = string(re.ReplaceAll([]byte(vt), []byte(tsize)))
		} else if strings.Contains(vt, `{name}`) && tname != "" {
			re := regexp.MustCompile(`{name}`)
			vt = string(re.ReplaceAll([]byte(vt), []byte(tname)))
		} else if strings.Contains(vt, `{icon}`) && tdate != "" {
			re := regexp.MustCompile(`{icon}`)
			vt = string(re.ReplaceAll([]byte(vt), []byte(ic)))
		} else if strings.Contains(vt, `{icon}`) && ic != "" {
			re := regexp.MustCompile(`{icon}`)
			vt = string(re.ReplaceAll([]byte(vt), []byte(ic)))
		} else if strings.Contains(vt, `{symlink}`) && item.Symlink != "" {
			re := regexp.MustCompile(`{symlink}`)
			vt = string(re.ReplaceAll([]byte(vt), []byte(item.Symlink)))
		} else {
			continue
		}
		t += vt
		if i != len(conf.FileInfoFormat)-1 {
			t += conf.InfoSeparator
		}
	}
	return t
}
