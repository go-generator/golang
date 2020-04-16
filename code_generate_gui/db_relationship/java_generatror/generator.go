package java_generatror

import "golang/code_generate_gui/db_relationship/constants"

func WritePackage(pkg string) string {
	return "package " + pkg
}

func WriteImportOneToMany() string {
	return constants.ImportOneToMany
}

func WriteImportPK() string {
	return constants.ImportComPK
}
