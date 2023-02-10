package service

func GetCommonOutput() map[string]interface{} {
	output := make(map[string]interface{})
	output["menu"] = GetMenu()
	return output
}
