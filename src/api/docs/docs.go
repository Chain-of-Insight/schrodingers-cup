// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ping": {
            "get": {
                "description": "Ping/Health Check",
                "produces": [
                    "text/plain"
                ],
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/settle": {
            "post": {
                "description": "Settle game window"
            }
        },
        "/test": {
            "post": {
                "description": "Run some arbitrary nomsu code",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Nomsu code to run",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.TestInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Nomsu execution result",
                        "schema": {
                            "$ref": "#/definitions/handlers.TestResult"
                        }
                    }
                }
            }
        },
        "/vote": {
            "post": {
                "description": "Receive and tabulate votes, stage vote outcome to be processed when game window is settled"
            }
        }
    },
    "definitions": {
        "handlers.TestInput": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "handlers.TestResult": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "string"
                },
                "resultHtml": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Schrodinger's Cup API",
	Description: "A game of Peter Suber's Nomic running on the Tezos network.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
