package provider_schema_generator

import (
    "encoding/json"
    "fmt"
    "github.com/hashicorp/terraform/helper/schema"
    "os"
    "path/filepath"
    "reflect"
)

func Generate(provider *schema.Provider) {
    if len(os.Args) < 3 {
        fmt.Fprintln(os.Stderr, "Usage: <program> <provider name> <output directory>")
        os.Exit(3)
    }

    providerName := os.Args[1]

    outputFilePath := filepath.Join(os.Args[2], fmt.Sprintf("%s.json", providerName))

    if err := DoGenerate(provider, providerName, outputFilePath); err != nil {
        fmt.Fprintln(os.Stderr, "Error: ", err.Error())
        os.Exit(255)
    }
}

func DoGenerate(provider *schema.Provider, providerName string, outputFilePath string) error {
    providerJson, err := json.MarshalIndent(&provider_schema{
        Name:        providerName,
        Type:        "provider",
        Provider:    getObjectSchema(provider.Schema),
        Resources:   getResourcesSchema(provider.ResourcesMap),
        DataSources: getResourcesSchema(provider.DataSourcesMap),
    }, "", "  ")

    if err != nil {
        return err
    }

    file, err := os.Create(outputFilePath)
    if err != nil {
        return err
    }

    defer file.Close()

    _, err = file.Write(providerJson)
    if err != nil {
        return err
    }

    return file.Sync()
}

func getResourcesSchema(resources map[string]*schema.Resource) resources_schema {
    if resources == nil {
        return nil
    }

    resourcesSchema := make(resources_schema)
    for name, resource := range resources {
        resourcesSchema[name] = getObjectSchema(resource.Schema)
    }

    return resourcesSchema
}

func getObjectSchema(fields map[string]*schema.Schema) object_schema {
    objectSchema := make(object_schema)
    for name, field := range fields {
        objectSchema[name] = getFieldSchema(field)
    }
    return objectSchema
}

func getFieldSchema(field *schema.Schema) field_schema {
    fieldSchema := make(field_schema)

    fieldValue := reflect.ValueOf(field).Elem()
    fieldType := fieldValue.Type()

    for i := 0; i < fieldValue.NumField(); i++ {
        option := fieldType.Field(i)
        if option.Type.Kind() == reflect.Func {
            continue
        }

        option_value := fieldValue.Field(i).Interface()

        if !reflect.DeepEqual(option_value, reflect.Zero(option.Type).Interface()) {
            fieldSchema[option.Name] = option_value
        }
    }

    fieldSchema["Type"] = field.Type.String()

    if field.Default != nil {
        defaultValue := reflect.ValueOf(field.Default)
        fieldSchema["Default"] = &default_value_schema{
            Type:  defaultValue.Type().String(),
            Value: fmt.Sprintf("%v", defaultValue.Interface()),
        }
    }

    if field.Elem != nil {
        if schemaElem, ok := field.Elem.(*schema.Schema); ok {
            fieldSchema["Elem"] = &elem_schema{
                Type:     "SchemaElements",
                ElemType: schemaElem.Type.String(),
            }
        } else if resourceElem, ok := field.Elem.(*schema.Resource); ok {
            fieldSchema["Elem"] = &elem_schema{
                Type: "SchemaInfo",
                Info: getObjectSchema(resourceElem.Schema),
            }
        }
    }

    return fieldSchema
}

type provider_schema struct {
    Name        string           `json:"name"`
    Type        string           `json:"type"`
    Provider    object_schema    `json:"provider"`
    Resources   resources_schema `json:"resources"`
    DataSources resources_schema `json:"data-sources"`
}

type resources_schema map[string]object_schema

type object_schema map[string]field_schema

type field_schema map[string]interface{}

type default_value_schema struct {
    Type  string `json:"type"`
    Value string `json:"value"`
}

type elem_schema struct {
    Type     string        `json:"type"`
    ElemType string        `json:"elements-type,omitempty"`
    Info     object_schema `json:"info,omitempty"`
}
