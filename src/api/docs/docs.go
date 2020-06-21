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
        "/auth": {
            "post": {
                "description": "Authenticate a user with a signed tezos message",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "Signed message",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AuthInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT token",
                        "schema": {
                            "$ref": "#/definitions/handlers.AuthResult"
                        }
                    }
                }
            }
        },
        "/game/propose": {
            "post": {
                "description": "Submit a new rule proposal",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Proposed rule",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RuleProposal"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.ProposalResult"
                        }
                    }
                }
            }
        },
        "/game/vote": {
            "post": {
                "description": "Receive and tabulate votes, stage vote outcome to be processed when game window is settled",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ]
            }
        },
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
        "/players": {
            "get": {
                "description": "Players",
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "List of players in the current gameplay session",
                        "schema": {
                            "$ref": "#/definitions/handlers.PlayerList"
                        }
                    }
                }
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
        }
    },
    "definitions": {
        "handlers.AuthInput": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "msg": {
                    "type": "string"
                },
                "pubKey": {
                    "type": "string"
                },
                "sig": {
                    "type": "string"
                }
            }
        },
        "handlers.AuthResult": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "handlers.PlayerList": {
            "type": "object",
            "properties": {
                "currentTurn": {
                    "type": "string"
                },
                "nextTurn": {
                    "type": "string"
                },
                "nextTurnAt": {
                    "type": "string"
                },
                "players": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handlers.ProposalResult": {
            "type": "object",
            "properties": {
                "message": {
                    "description": "\"OK!\" or error message",
                    "type": "string"
                },
                "round": {
                    "description": "Updated round value",
                    "type": "integer"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "handlers.RuleProposal": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Nomsu code",
                    "type": "string"
                },
                "index": {
                    "description": "rule index of the existing rule",
                    "type": "integer"
                },
                "kind": {
                    "description": "Mutable / Immutable",
                    "type": "string"
                },
                "type": {
                    "description": "Update, Create, Delete, Transmute",
                    "type": "string"
                }
            }
        },
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
