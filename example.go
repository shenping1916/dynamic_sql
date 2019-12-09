package dynamic

import "fmt"

func StringApostrophe(s string) string {
	return fmt.Sprintf("'%s'", s)
}

func GenerateDynamicSqlToCodeCoveragePeriod(queryParams QueryParams) string {
	var dqb DynamicQueryBuilder
	query := dqb.And(
		dqb.NewExp("app.app_name", "=", StringApostrophe(queryParams["app_name"])),
		dqb.NewExp("app.guid", "=", "app_tag.app_guid"),
		dqb.NewExp("app_tag.tag", "=", StringApostrophe(queryParams["tag"])),
		dqb.NewExp("app_tag.guid", "=", "app_tag_batch.tag_guid"),
		dqb.NewExp("app_tag_batch.env", "=", StringApostrophe(queryParams["env"])),
		dqb.NewExp("app_tag_batch.guid", "=", "code_execute.batch_guid"),
		dqb.NewExp("code_execute.update_time", ">=", StringApostrophe(queryParams["start_time"])),
		dqb.NewExp("code_execute.update_time", "<=", StringApostrophe(queryParams["end_time"])),
	).GroupBy("code_execute.create_time", "code_execute.update_time").
		OrderBy("code_execute.update_time", "ASC").
		BindSql(
			`
               SELECT
                   TO_CHAR(code_execute.create_time,'YYYY-MM-DD HH24:MI:SS.MS') create_time,
                   TO_CHAR(code_execute.update_time,'YYYY-MM-DD HH24:MI:SS.MS') update_time
               FROM
                   t_jcc_app app,
                   t_jcc_app_tag app_tag,
                   t_jcc_app_tag_batch app_tag_batch,
                   t_jcc_code_execute code_execute
    `)
	return query
}

