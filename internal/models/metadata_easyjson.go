// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	time "time"

	uuid "github.com/google/uuid"
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

func easyjsonBa0ee0e3DecodeGithubComFlutterDizasterFileServerInternalModels(in *jlexer.Lexer, out *Metadata) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			if in.IsNull() {
				in.Skip()
				out.ID = nil
			} else {
				if out.ID == nil {
					out.ID = new(uuid.UUID)
				}
				if data := in.UnsafeBytes(); in.Ok() {
					in.AddError((*out.ID).UnmarshalText(data))
				}
			}
		case "name":
			out.Name = string(in.String())
		case "file":
			out.File = bool(in.Bool())
		case "public":
			out.Public = bool(in.Bool())
		case "mime":
			out.Mime = string(in.String())
		case "created":
			if in.IsNull() {
				in.Skip()
				out.Created = nil
			} else {
				if out.Created == nil {
					out.Created = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Created).UnmarshalJSON(data))
				}
			}
		case "owner_id":
			if in.IsNull() {
				in.Skip()
				out.OwnerID = nil
			} else {
				if out.OwnerID == nil {
					out.OwnerID = new(uuid.UUID)
				}
				if data := in.UnsafeBytes(); in.Ok() {
					in.AddError((*out.OwnerID).UnmarshalText(data))
				}
			}
		case "grant":
			if in.IsNull() {
				in.Skip()
				out.Grant = nil
			} else {
				in.Delim('[')
				if out.Grant == nil {
					if !in.IsDelim(']') {
						out.Grant = make([]string, 0, 4)
					} else {
						out.Grant = []string{}
					}
				} else {
					out.Grant = (out.Grant)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Grant = append(out.Grant, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "json":
			(out.JSON).UnmarshalEasyJSON(in)
		case "file-size":
			out.FileSize = int64(in.Int64())
		case "url":
			out.URL = string(in.String())
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
func easyjsonBa0ee0e3EncodeGithubComFlutterDizasterFileServerInternalModels(out *jwriter.Writer, in Metadata) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != nil {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((*in.ID).MarshalText())
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.File {
		const prefix string = ",\"file\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.File))
	}
	if in.Public {
		const prefix string = ",\"public\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Public))
	}
	if in.Mime != "" {
		const prefix string = ",\"mime\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Mime))
	}
	if in.Created != nil {
		const prefix string = ",\"created\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((*in.Created).MarshalJSON())
	}
	if in.OwnerID != nil {
		const prefix string = ",\"owner_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((*in.OwnerID).MarshalText())
	}
	if len(in.Grant) != 0 {
		const prefix string = ",\"grant\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v2, v3 := range in.Grant {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	if in.JSON != "" {
		const prefix string = ",\"json\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.JSON).MarshalEasyJSON(out)
	}
	if in.FileSize != 0 {
		const prefix string = ",\"file-size\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.FileSize))
	}
	if in.URL != "" {
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Metadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComFlutterDizasterFileServerInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Metadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComFlutterDizasterFileServerInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComFlutterDizasterFileServerInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Metadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComFlutterDizasterFileServerInternalModels(l, v)
}
