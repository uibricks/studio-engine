package constants

const (
	MappingTableProjectColumnName = "project_version"
)

const (
	QueryCopyRepositoriesToHistory = `insert into repositories_history (project_id, config, project_version, repos_created_at, repos_deleted_at)
	select project_id, config, project_version, created_at, deleted_at from 
	repositories o where project_id = '%s' %s`

	QueryDeleteObsoleteRepositories = `delete from repositories where project_id = '%s' %s`

	QueryDuplicateActiveMapping = `insert into repositories(project_id,project_version,config) 
	select project_id, %d, config from 
	repositories where project_id='%s' order by created_at desc limit 1;`

	QueryInsertEmptyMapping = `insert into repositories(project_id ,project_version) values('%s',%d);`

	QueryDeleteFromRepositoriesHistory = `DELETE FROM repositories_history where project_id=?`
)
