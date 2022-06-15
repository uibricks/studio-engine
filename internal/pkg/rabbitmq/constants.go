package rabbitmq

const (
	Mapping_Queue_Name    string = "mapping"
	Project_Queue_Name           = "project"
	Expression_Queue_Name        = "expression"
)

const (
	Action_Save_Mapping       string = "save_mapping"
	Action_Delete_Mapping     string = "delete_mapping"
	Action_Resolve_Expression string = "eval_expression"
	Action_Restore_Mapping    string = "restore_mapping"
)

const (
	Status_Success             string = "success"
	Status_Error                      = "error"
	reconnect_Delay_In_Seconds        = 3
)
