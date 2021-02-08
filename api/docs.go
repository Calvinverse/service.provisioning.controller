// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package api

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
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/doc": {
            "get": {
                "description": "Returns the OpenAPI document for the current service",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "doc"
                ],
                "summary": "Returns the OpenAPI document for the current service",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/environment.Environment"
                        }
                    }
                }
            }
        },
        "/v1/environment": {
            "put": {
                "description": "Creates a new environment based on the provided information.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "environment"
                ],
                "summary": "Creates a new environment.",
                "parameters": [
                    {
                        "description": "Environment ID",
                        "name": "id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/environment.Environment"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/environment.Environment"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "/v1/environment/": {
            "get": {
                "description": "Returns a list of known environment IDs.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "environment"
                ],
                "summary": "Provide the list of known environment IDs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "/v1/environment/{id}": {
            "get": {
                "description": "Returns information about the environment with the given id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "environment"
                ],
                "summary": "Provide information about an environment.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Environment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/environment.Environment"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes the environment with the given id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "environment"
                ],
                "summary": "Deletes an environment.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Environment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/environment.Environment"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "/v1/self/info": {
            "get": {
                "description": "Respond to an info request with information about the application.",
                "consumes": [
                    "application/json",
                    "text/xml"
                ],
                "produces": [
                    "application/json",
                    "text/xml"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Respond to an info request",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/health.InfoResponse"
                        }
                    },
                    "415": {
                        "description": "Unsupported media type",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/self/liveliness": {
            "get": {
                "description": "Respond to an liveliness request with information about the status of the latest health checks.",
                "consumes": [
                    "application/json",
                    "text/xml"
                ],
                "produces": [
                    "application/json",
                    "text/xml"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Respond to an liveliness request",
                "parameters": [
                    {
                        "enum": [
                            "summary",
                            "detailed"
                        ],
                        "type": "string",
                        "description": "options are summary or detailed",
                        "name": "type",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/health.LivelinessDetailedResponse"
                        }
                    },
                    "415": {
                        "description": "Unsupported media type",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/self/ping": {
            "get": {
                "description": "Respond to a ping request with a pong response.",
                "consumes": [
                    "application/json",
                    "text/xml"
                ],
                "produces": [
                    "application/json",
                    "text/xml"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Respond to a ping request",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/health.PingResponse"
                        }
                    },
                    "415": {
                        "description": "Unsupported media type",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/self/readiness": {
            "get": {
                "description": "Respond to an readiness request with information about ability of the application to start serving requests.",
                "consumes": [
                    "application/json",
                    "text/xml"
                ],
                "produces": [
                    "application/json",
                    "text/xml"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Respond to an readiness request",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/health.ReadinessResponse"
                        }
                    }
                }
            }
        },
        "/v1/self/started": {
            "get": {
                "description": "Respond to an started request with information indicating if the application has started successfully.",
                "consumes": [
                    "application/json",
                    "text/xml"
                ],
                "produces": [
                    "application/json",
                    "text/xml"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Respond to an started request",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/health.StartedResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "environment.Environment": {
            "type": "object"
        },
        "health.DetailedCheckInformation": {
            "type": "object",
            "properties": {
                "name": {
                    "description": "Name returns the name of the health check.",
                    "type": "string"
                },
                "status": {
                    "description": "Status returns the status of the health check, either success or failure.",
                    "type": "string"
                },
                "timestamp": {
                    "description": "Timestamp returns the time the healtcheck was executed.",
                    "type": "string"
                }
            }
        },
        "health.InfoResponse": {
            "type": "object",
            "properties": {
                "buildtime": {
                    "description": "BuildTime stores the date and time the application was built.",
                    "type": "string"
                },
                "revision": {
                    "description": "Revision stores the GIT SHA of the commit on which the application build was based.",
                    "type": "string"
                },
                "version": {
                    "description": "Version stores the version number of the application.",
                    "type": "string"
                }
            }
        },
        "health.LivelinessDetailedResponse": {
            "type": "object",
            "properties": {
                "checks": {
                    "description": "Status of all the health checks",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/health.DetailedCheckInformation"
                    }
                },
                "status": {
                    "description": "Global status",
                    "type": "string"
                },
                "time": {
                    "description": "Time the liveliness response was created at",
                    "type": "string"
                }
            }
        },
        "health.PingResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "string"
                }
            }
        },
        "health.ReadinessResponse": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                }
            }
        },
        "health.StartedResponse": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
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
	Host:        "localhost:8080",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Provisioning.Controller server API",
	Description: "Provides information about deployed environments and the templates used to created these environments.",
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
