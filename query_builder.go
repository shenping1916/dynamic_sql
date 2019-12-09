package dynamic

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

type DynamicQueryBuilder string

type QueryParams map[string]string

func (qp QueryParams) GetInt(key string) interface{} {
	mapVal := qp[key]

	if mapVal == "" {
		return ""
	} else {
		val, err := strconv.Atoi(mapVal)
		if err != nil {
			return ""
		}
		return val
	}
}

func (qp QueryParams) GetString(key string) interface{} {
	return qp[key]
}

type Expression struct {
	Key   string
	Exp   string
	Value interface{}
}

func (dqb DynamicQueryBuilder) NewExp(key string, assignment string, value interface{}) Expression {
	return Expression{Key: key, Exp: assignment, Value: value}
}

func componentToString(c interface{}) DynamicQueryBuilder {
	switch v := c.(type) {
	case Expression:
		return DynamicQueryBuilder(c.(Expression).ToString())
	case string, *string:
		return DynamicQueryBuilder(c.(string))
	case DynamicQueryBuilder:
		return v
	default:
		return ""
	}
}

func (dqb DynamicQueryBuilder) And(component ...interface{}) DynamicQueryBuilder {
	return dqb.getOperationExpression("AND", component...)
}

func (dqb DynamicQueryBuilder) OR(component ...interface{}) DynamicQueryBuilder {
	return dqb.getOperationExpression("OR", component...)
}

func (dqb DynamicQueryBuilder) GroupBy(component ...string) DynamicQueryBuilder {
	query := dqb
	if len(component) == 1 {
		return query + DynamicQueryBuilder("GROUP BY "+component[0]+" ")
	}
	if len(component) > 1 {
		return query + DynamicQueryBuilder("GROUP BY "+strings.Join(component, ",")+" ")
	}

	return ""
}

func (dqb DynamicQueryBuilder) OrderBy(component ...string) DynamicQueryBuilder {
	query := dqb
	if len(component) == 1 {
		return query + DynamicQueryBuilder("ORDER BY "+component[0]+" ")
	}

	if len(component) > 1 {
		s := strings.Join(component[0:len(component)-1], ",") + " " + component[len(component)-1]
		return query + DynamicQueryBuilder("ORDER BY "+s+" ")
	}

	return ""
}

func (dqb DynamicQueryBuilder) getOperationExpression(operation string, component ...interface{}) DynamicQueryBuilder {
	if len(component) == 0 {
		return ""
	}
	if len(component) == 1 {
		return componentToString(component[0])
	} else {
		clauses := make([]string, 0, len(component))
		for _, v := range component {
			value := componentToString(v)
			if value != "" {
				clauses = append(clauses, " "+string(value)+" ")
			}
		}

		if len(clauses) > 0 {
			return DynamicQueryBuilder("(" + strings.Join(clauses, ""+operation+"") + ") ")
		}

		return ""
	}
}

func (dqb DynamicQueryBuilder) Limit(offset int, length int) DynamicQueryBuilder {
	query := string(dqb)
	query += " LIMIT " + strconv.Itoa(length) + " OFFSET " + strconv.Itoa(offset)
	return DynamicQueryBuilder(query)
}

func (dqb DynamicQueryBuilder) CopyQuery(dest *string) DynamicQueryBuilder {
	*dest = dqb.ToString()
	return dqb
}

func (dqb DynamicQueryBuilder) BindSql(sql string) string {
	if dqb != "" && dqb != "( )" {
		index := strings.Index(dqb.ToString(), "LIMIT")
		if index == 1 {
			return sql + dqb.ToString()
		}

		return sql + " WHERE " + string(dqb)
	}
	return sql
}

func (dqb DynamicQueryBuilder) ToString() string {
	return string(dqb)
}

func (e Expression) ToString() string {
	switch e.Value.(type) {
	case int, int16, int32, int64:
		val := strconv.Itoa(e.Value.(int))
		clause := e.Key + " " + e.Exp + " " + e.getReplaceExp()
		return fmt.Sprintf(clause, val)

	default:
		if trim := strings.TrimSpace(e.Value.(string)); (trim == "") || (trim == "''") {
			return ""
		} else {
			e.Value = template.HTMLEscapeString(e.Value.(string))
			e.Value = strings.Replace(e.Value.(string), "&#39;", "'", -1)
			clause := e.Key + " " + e.Exp + " " + e.getReplaceExp()
			val := fmt.Sprintf(clause, e.Value)
			return val
		}
	}
}

func (e Expression) getReplaceExp() string {
	switch e.Value.(type) {
	case int, int64, int32, int16:
		return "%d"
	case string:
		if strings.Contains(e.Value.(string), "&#39;") {
			return "'%s'"
		}
		return "%s"
	default:
		return "'%s'"
	}
}
