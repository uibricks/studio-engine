package constants

import (
	"github.com/uibricks/studio-engine/internal/pkg/constants"
)

const (
	QueryCopyProjectsToHistory = `insert into project_history (name, client_id, luid, version, user_version, config, state, project_created_at,project_updated_at, project_deleted_at, created_by, updated_by)
	select name, client_id, luid, version, user_version, config, state, created_at,updated_at, deleted_at, created_by, updated_by from 
	(select row_number() over (partition by state order by updated_at desc) rnum, o.*
	from objects o where luid = '%s') p
	where (state != 'active') or (state = '` + constants.STATE_ACTIVE + `' and p.rnum > 2);`

	QueryDeleteObsoleteProjects = `delete from objects where id in (
	select id from (select row_number() over (partition by state order by updated_at desc) rnum, o.*
	from objects o where luid = '%s') p
	where (state != 'active') or (state = '` + constants.STATE_ACTIVE + `' and p.rnum > 2));`

	QueryFetchActiveProjects = `SELECT * FROM ( 
		SELECT *, ROW_NUMBER() OVER (PARTITION BY luid ORDER BY updated_at DESC) _col 
		FROM objects where (state='` + constants.STATE_ACTIVE + `' or state='` + constants.STATE_CACHE + `') and deleted_at is null and %s) x  
		WHERE x._col = 1;`

	QueryFetchRecentProjects = `SELECT * FROM ( 
		SELECT *, ROW_NUMBER() OVER (PARTITION BY luid ORDER BY updated_at DESC) _col 
		FROM objects where (state='` + constants.STATE_ACTIVE + `' or state='` + constants.STATE_CACHE + `') and updated_at > '%s' and deleted_at is null) x  
		WHERE x._col = 1;`

	QueryFetchTrashProjects = `SELECT * FROM ( 
		SELECT *, ROW_NUMBER() OVER (PARTITION BY luid ORDER BY updated_at DESC) _col 
		FROM objects where (state='` + constants.STATE_ACTIVE + `' or state='` + constants.STATE_CACHE + `') and deleted_at is not null) x  
		WHERE x._col = 1;`

	QueryDeleteFromProjectHistory = `DELETE FROM project_history where luid::text=?`

	QueryFetchProjectVersionsByLuid = `select * from (
	select version, user_version, created_at, updated_at, state, created_by, updated_by from objects where luid = '%[1]s'
	union
	select version, user_version, project_created_at as created_at, project_updated_at as updated_at, state, created_by, updated_by from project_history where luid = '%[1]s') versions 
	order by updated_at desc;`
)
