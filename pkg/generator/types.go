package generator

type TypeMapper struct {
	refs   map[string]*Type
	names  map[string]struct{}
	prefix string
}

func NewTypeMapper(prefix string) *TypeMapper {
	return &TypeMapper{
		refs:   make(map[string]*Type),
		names:  make(map[string]struct{}),
		prefix: prefix,
	}
}

func (t *TypeMapper) Add(from string, schema *Type) *TypeMapper {
	t.refs[t.prefix+from] = schema
	t.names[from] = struct{}{}
	return t
}

func (t *TypeMapper) Exists(name string) bool {
	_, ok := t.names[name]
	return ok
}

func (t *TypeMapper) Convert(schema *Type) {
	if to, ok := t.refs[schema.Ref]; ok {
		// *schema.Type = *to
		schema.Ref = ""
		schema.Type = to.Type
		schema.Format = to.Format
	}
}

var mapper = NewTypeMapper("#/definitions/").
	Add("id", &Type{Type: "string"}).
	Add("string", &Type{Type: "string"}).
	Add("base64Binary", &Type{Type: "string", Format: "byte"}).
	Add("boolean", &Type{Type: "boolean"}).
	Add("canonical", &Type{Type: "string"}).
	Add("code", &Type{Type: "string"}).
	Add("date", &Type{Type: "string", Format: "date"}).
	Add("dateTime", &Type{Type: "string", Format: "date-time"}).
	Add("decimal", &Type{Type: "number"}).
	Add("instant", &Type{Type: "string"}).
	Add("integer", &Type{Type: "integer"}).
	Add("markdown", &Type{Type: "string"}).
	Add("oid", &Type{Type: "string"}).
	Add("positiveInt", &Type{Type: "integer"}).
	Add("time", &Type{Type: "string", Pattern: "^([01][0-9]|2[0-3]):[0-5][0-9]:([0-5][0-9]|60)(\\.[0-9]+)?$"}).
	Add("unsignedInt", &Type{Type: "integer"}).
	Add("uri", &Type{Type: "string", Format: "uri"}).
	Add("url", &Type{Type: "string", Format: "uri"}).
	Add("uuid", &Type{Type: "uuid", Format: "uuid"}).
	Add("xhtml", &Type{Type: "string"})
