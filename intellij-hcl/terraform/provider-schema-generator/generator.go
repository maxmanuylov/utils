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
    json, err := json.MarshalIndent(provider_schema{
        providerName, "provider", getObjectSchema(provider.Schema), getResourcesSchema(provider.ResourcesMap),
    }, "", "  ")

    if err != nil {
        return err
    }

    file, err := os.Create(outputFilePath)
    if err != nil {
        return err
    }

    defer file.Close()

    _, err = file.Write(json)
    if err != nil {
        return err
    }

    if err = file.Sync(); err != nil {
        return err
    }

    return nil
}

func getResourcesSchema(resources map[string]*schema.Resource) resources_schema {
    schema := make(resources_schema)
    for name, resource := range resources {
        schema[name] = getObjectSchema(resource.Schema)
    }
    return schema
}

func getObjectSchema(fields map[string]*schema.Schema) object_schema {
    schema := make(object_schema)
    for name, field := range fields {
        schema[name] = getFieldSchema(field)
    }
    return schema
}

func getFieldSchema(field *schema.Schema) field_schema {
    schema := make(field_schema, 0)

    fieldValue := reflect.ValueOf(field).Elem()
    fieldType := fieldValue.Type()

    for i := 0; i < fieldValue.NumField(); i++ {
        option := fieldType.Field(i)
        option_value := fieldValue.Field(i).Interface()

        if !reflect.DeepEqual(option_value, reflect.Zero(option.Type).Interface()) {
            schema = append(schema, field_option_schema{
                option.Name, option.Type.String(), fmt.Sprintf("%v", option_value),
            })
        }
    }

    return schema
}

type provider_schema struct {
    Name string `json:"name"`
    Type string `json:"type"`
    Schema object_schema `json:"schema"`
    Resources resources_schema `json:"resources"`
}

type resources_schema map[string]object_schema

type object_schema map[string]field_schema

type field_schema []field_option_schema

type field_option_schema struct {
    Name string `json:"name"`
    Type string `json:"type"`
    Value string `json:"value"`
}
