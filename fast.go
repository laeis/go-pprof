package main

import (
	"bufio"
	"bytes"

	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)
//main use here just for print log in console
func main() {
	out := new(bytes.Buffer)
	FastSearch(out)
}

//User record about user what wiil be prase from file
type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"-"`
	Country  string   `json:"-"`
	Email    string   `json:"email"`
	Job      string   `json:"-"`
	Name     string   `json:"name"`
	Phone    string   `json:"-"`
}
// for don't allocate memory in foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
func write(foundUsers *bytes.Buffer, pos int, name string, email string) {
	foundUsers.WriteByte('[')
	foundUsers.Write(strconv.AppendInt([]byte(""), int64(pos), 10))
	foundUsers.WriteByte(']')
	foundUsers.WriteByte(' ')
	foundUsers.WriteString(name)
	foundUsers.WriteByte(' ')
	foundUsers.WriteByte('<')
	foundUsers.WriteString(email)
	foundUsers.WriteByte('>')
	foundUsers.WriteRune('\n')
}

//FastSearch вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	} else {
		defer file.Close()
	}
	scanner := bufio.NewScanner(file)
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	// foundUsers := ""
	foundUsers := &bytes.Buffer{}
	//hack for counter
	//TODO: bad code, need find best solution
	i := -1

	var email string

	user := &User{}

	for scanner.Scan() {
		line := scanner.Bytes()

		err := user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}
		//hack for counter
		//TODO: bad code, need find best solution
		i++
		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {

			notSeenBefore := false
			if ok := strings.Contains(browserRaw, "Android"); ok {
				isAndroid = true
				notSeenBefore = true
			}
			if ok := strings.Contains(browserRaw, "MSIE"); ok {
				isMSIE = true
				notSeenBefore = true
			}
			for _, item := range seenBrowsers {
				if item == browserRaw {
					notSeenBefore = false
				}
			}
			if notSeenBefore {
				seenBrowsers = append(seenBrowsers, browserRaw)
				uniqueBrowsers++
			}
		}
		if !(isAndroid && isMSIE) {
			continue
		}
		email = strings.Replace(user.Email, "@", " [at] ", -1) // r.ReplaceAllString(user["email"].(string), " [at] ")
		//TODO:: change on buff output
		write(foundUsers, i, user.Name, email)
	}
	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	// fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

/*easyjson code*/
func easyjson9e1087fdDecodeHomeArtGoSrc(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9e1087fdEncodeHomeArtGoSrc(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeHomeArtGoSrc(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeHomeArtGoSrc(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeHomeArtGoSrc(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeHomeArtGoSrc(l, v)
}
