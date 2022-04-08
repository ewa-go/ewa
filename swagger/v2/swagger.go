package v2

type Swagger struct {
	Swagger             string                `json:"swagger"`
	Info                Info                  `json:"info"`
	Host                string                `json:"host"`
	BasePath            string                `json:"basePath"`
	Tags                []Tag                 `json:"tags"`
	Schemes             []string              `json:"schemes"`
	Paths               Paths                 `json:"paths"`
	SecurityDefinitions SecurityDefinitions   `json:"securityDefinitions"`
	Definitions         map[string]Definition `json:"definitions"`
}

type Definition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Xml        Xml                 `json:"xml"`
}

type Xml struct {
	Name string `json:"name"`
}

type Property struct {
	Type        string `json:"type"`
	Format      string `json:"format,omitempty"`
	Description string `json:"description"`
}

type Info struct {
	Description    string  `json:"description"`
	Version        string  `json:"version"`
	Title          string  `json:"title"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
}

type Contact struct {
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Tag struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ExternalDocs ExternalDocs `json:"externalDocs"`
}

type ExternalDocs struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Path struct {
	Tags        []string         `json:"tags"`
	Summary     string           `json:"summary"`
	Description string           `json:"description"`
	OperationId string           `json:"operationId"`
	Consumers   []string         `json:"consumers"`
	Produces    []string         `json:"produces"`
	Parameters  []Parameter      `json:"parameters"`
	Responses   map[int]Response `json:"responses"`
	Security    Security         `json:"security"`
}

type Paths map[string]Methods

type Methods map[string]*Path

type Secure map[string][]string

type Security []Secure

type Response struct {
	Description string            `json:"description"`
	Schema      map[string]string `json:"schema"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
	Type        string      `json:"type,omitempty"`
	Format      string      `json:"format,omitempty"`
	Minimum     int         `json:"minimum,omitempty"`
	Schema      Schema      `json:"schema,omitempty"`
	Body        interface{} `json:"-"`
}

type Schema struct {
	Type  string            `json:"type,omitempty"`
	Items map[string]string `json:"items,omitempty"`
}

type SecurityDefinition struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type SecurityDefinitions map[string]SecurityDefinition

const (
	ParameterTypePath = "path"
	ParameterTypeBody = "body"
)

// New инициализация swagger
func New(host string, info Info) *Swagger {
	return &Swagger{
		Swagger: "2.0",
		Info:    info,
		Host:    host,
		Paths:   Paths{},
	}
}

// SetSchemes устанавливаем схему
func (s *Swagger) SetSchemes(scheme ...string) *Swagger {
	s.Schemes = append(s.Schemes, scheme...)
	return s
}

/*
func (p Parameter) Marshal() *Definition {
	if p.Body == nil {
		return nil
	}
	reflect.TypeOf(p.Body).
}*/
