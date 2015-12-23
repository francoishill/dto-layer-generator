package setup

import (
	"fmt"
)

type PluginName string

func (p PluginName) String() string { return string(p) }

type OutputFilePath string

func (o OutputFilePath) String() string { return string(o) }

type dtoSetupYAML struct {
	Name   string
	Url    string
	Output struct {
		Placeholder string
		Plugins     map[PluginName]OutputFilePath
	}

	AllFields                []*DTOField `json:"all_fields"`
	IdFieldName              string      `json:"id_field_name"`
	ListableFieldNameGroups  [][]string  `json:"listable_field_name_groups"`
	GetableFieldNameGroups   [][]string  `json:"getable_field_name_groups"`
	PatchableFieldNameGroups [][]string  `json:"patchable_field_name_groups"`
}

type DTOSetup struct {
	*dtoSetupYAML

	IdField              *DTOField
	InsertableFields     []*DTOField
	ListableFieldGroups  [][]*DTOField
	GetableFieldGroups   [][]*DTOField
	PatchableFieldGroups [][]*DTOField
}

func NewDTOSetupFromYAML(setup *dtoSetupYAML) *DTOSetup {
	d := &DTOSetup{dtoSetupYAML: setup}
	d.IdField = d.getIdField()
	d.InsertableFields = d.getInsertableFields()
	d.ListableFieldGroups = d.getListableFieldGroups()
	d.GetableFieldGroups = d.getGetableFieldGroups()
	d.PatchableFieldGroups = d.getPatchableFieldGroups()
	return d
}

func (d *DTOSetup) getFieldByName(name string) *DTOField {
	for _, f := range d.AllFields {
		if f.Name == name {
			return f
		}
	}
	panic(fmt.Sprintf("Field name '%s' is not in the field list", name))
}

func (d *DTOSetup) getGroupedFieldsByNames(groupedFieldNames [][]string) [][]*DTOField {
	groups := [][]*DTOField{}
	for _, g := range groupedFieldNames {
		fieldsInGroup := []*DTOField{}
		for _, fieldName := range g {
			fieldsInGroup = append(fieldsInGroup, d.getFieldByName(fieldName))
		}
		groups = append(groups, fieldsInGroup)
	}
	return groups
}

func (d *DTOSetup) getIdField() *DTOField {
	return d.getFieldByName(d.IdFieldName)
}

func (d *DTOSetup) getInsertableFields() (fields []*DTOField) {
	fields = []*DTOField{}
	for _, f := range d.AllFields {
		if f.Name == d.IdFieldName {
			continue //When inserting an entity we do not have an ID yet
		}
		fields = append(fields, f)
	}
	return
}

func (d *DTOSetup) getListableFieldGroups() [][]*DTOField {
	return d.getGroupedFieldsByNames(d.ListableFieldNameGroups)
}

func (d *DTOSetup) getGetableFieldGroups() [][]*DTOField {
	return d.getGroupedFieldsByNames(d.GetableFieldNameGroups)
}

func (d *DTOSetup) getPatchableFieldGroups() [][]*DTOField {
	return d.getGroupedFieldsByNames(d.PatchableFieldNameGroups)
}
