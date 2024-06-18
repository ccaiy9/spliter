package spliter

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Option struct {
	AbsPath              string
	Index                string
	splitCompeteSqlChunk int // byte
}

func doWrite2Local(sql *string, path string) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(*sql)
	if err != nil {
		return err
	}

	return nil
}

func SplitBatchInsertSql(sql *string, op ...*Option) ([]string, error) {
	var result []string

	pattern := "(?is)(INSERT|REPLACE)\\s+INTO\\s+(`?(\\w+)`?\\.)?(`?(\\w+)`?)\\s*(?:\\((.*?)\\))?\\s*VALUES\\s*(.*);"
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(*sql)

	if len(matches) < 7 {

		if len(op) > 0 && op[0] != nil {
			errFileName := fmt.Sprintf("err_split_sql_%s.log", op[0].Index)
			err := doWrite2Local(sql, op[0].AbsPath+errFileName)
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil, errors.New("match compressed sql error")
	}

	chunk := 1 * 1024 * 1024
	if len(op) > 0 && op[0].splitCompeteSqlChunk > 0 {
		chunk = op[0].splitCompeteSqlChunk
	}

	operation := strings.ToUpper(strings.TrimSpace(matches[1]))
	dbName := strings.TrimSpace(matches[3])
	tableName := strings.TrimSpace(matches[5])

	if dbName == "" {
		tableName = "`" + tableName + "`"
	} else {
		tableName = "`" + dbName + "`.`" + tableName + "`"
	}

	var columns []string
	if matches[6] != "" {
		columns = strings.Split(matches[6], ",")
	}

	values := strings.TrimSpace(matches[7])
	valueStatements := strings.Split(values, "),")

	var batchValues []string
	var batchValueSize int

	for _, value := range valueStatements {
		value = strings.TrimSpace(value)
		value = strings.Trim(value, "()")
		value = "(" + value + ")"

		// Secondary compression to prevent OOM caused by large SQL
		if batchValueSize+len(value) >= chunk {
			batchValues = append(batchValues, value)
			valuesString := strings.Join(batchValues, ", ")
			statement := buildStatement(columns, operation, tableName, valuesString)
			result = append(result, statement)

			batchValues = batchValues[:0]
			batchValueSize = 0
		} else {
			batchValues = append(batchValues, value)
			batchValueSize += len(value)
		}
	}

	if len(batchValues) > 0 {
		remainString := strings.Join(batchValues, ", ")
		statement := buildStatement(columns, operation, tableName, remainString)
		result = append(result, statement)
	}

	return result, nil
}

func buildStatement(columns []string, operation string, tableName string, value string) string {
	var statement string
	if len(columns) > 0 {
		statement = fmt.Sprintf("%s INTO %s (%s) VALUES %s;", operation, tableName, strings.Join(columns, ","), value)
	} else {
		statement = fmt.Sprintf("%s INTO %s VALUES %s;", operation, tableName, value)
	}

	return statement
}
