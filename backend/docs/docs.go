// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "get jwt, return in header and cookie",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.login.body"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "create new account, login after success",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.register.body"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    }
                }
            }
        },
        "/story": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Story"
                ],
                "summary": "create new story",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Story"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Story"
                ],
                "summary": "create new story",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Story"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    }
                }
            }
        },
        "/story/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Story"
                ],
                "summary": "get story detail",
                "parameters": [
                    {
                        "type": "string",
                        "example": "51eecb74-bd12-40b4-bd3d-71eaa2a7d71b",
                        "description": "story id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Story"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Story"
                ],
                "summary": "create QA",
                "parameters": [
                    {
                        "type": "string",
                        "example": "51eecb74-bd12-40b4-bd3d-71eaa2a7d71b",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "What is a dog?",
                        "description": "question",
                        "name": "question",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "file",
                        "name": "attachment",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.QA"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.APIError": {
            "type": "object",
            "properties": {
                "error_code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "model.QA": {
            "type": "object",
            "properties": {
                "answer": {
                    "type": "string"
                },
                "extension": {
                    "type": "string"
                },
                "qa_id": {
                    "type": "string"
                },
                "question": {
                    "type": "string"
                }
            }
        },
        "model.Story": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.QA"
                    }
                },
                "creation_date": {
                    "type": "string"
                },
                "story_context": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "story_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "v1.login.body": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "abc@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "refo"
                }
            }
        },
        "v1.register.body": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "API Documentation",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
