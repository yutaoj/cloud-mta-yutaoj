package validate

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"cloud-mta/mta"
)

// GetValidationMode converts validation mode flags to validation process flags.
func GetValidationMode(validationFlag string) (bool, bool, error) {
	switch validationFlag {
	case "":
		return true, true, nil
	case "schema":
		return true, false, nil
	case "project":
		return false, true, nil
	}
	return false, false, errors.New("wrong argument of validation mode; expected argument of [all, schema, project]")
}

// MtaYaml validates an MTA.yaml file.
func MtaYaml(projectPath, mtaFilename string, validateSchema bool, validateProject bool) error {
	if validateProject || validateSchema {

		mtaPath := filepath.Join(projectPath, mtaFilename)
		// ParseFile contains MTA yaml content.
		yamlContent, err := ioutil.ReadFile(mtaPath)

		if err != nil {
			return errors.Wrapf(err, "Validation of %v failed on reading MTA content.", mtaPath)
		}
		// Validates MTA content.
		issues, err := validate(yamlContent, projectPath, validateSchema, validateProject)
		if len(issues) > 0 {
			return errors.Errorf("Validation of %v failed. Issues: \n%v %s", mtaPath, issues.String(), err)
		}
	}

	return nil
}

// validate validates an MTA schema.
func validate(yamlContent []byte, projectPath string, validateSchema bool, validateProject bool) (YamlValidationIssues, error) {
	//noinspection GoPreferNilSlice
	issues := []YamlValidationIssue{}
	if validateSchema {
		validations, schemaValidationLog := BuildValidationsFromSchemaText(schemaDef)
		if len(schemaValidationLog) > 0 {
			return schemaValidationLog, nil
		}
		yamlValidationLog, err := Yaml(yamlContent, validations...)
		if err != nil && len(yamlValidationLog) == 0 {
			yamlValidationLog = append(yamlValidationLog, []YamlValidationIssue{{Msg: "Validation failed." + err.Error()}}...)
		}
		issues = append(issues, yamlValidationLog...)

	}
	if validateProject {
		mtaStr := mta.MTA{}
		Unmarshal := yaml.Unmarshal
		err := Unmarshal(yamlContent, &mtaStr)
		if err != nil {
			return nil, errors.Wrap(err, "Read failed getting the MTA .yaml path while reading the MTA file.")
		}
		projectIssues := validateYamlProject(&mtaStr, projectPath)
		issues = append(issues, projectIssues...)
	}
	return issues, nil
}
