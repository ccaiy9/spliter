## Usage

```go
package main

import "fmt"

func main() {

	sql := "INSERT   INTo  `yu`.`table_name` vaLues (1257,34),(`ww`,`po`),(`wp`,`1`);"

	Option := &Option{
		AbsPath:              "/Users/Downloads/spliter/",
		Index:                "1",
		splitCompeteSqlChunk: 10,
	}

	stms, _ := SplitBatchInsertSql(&sql, Option)
	for _, st := range stms {
		fmt.Println(st)
	}
}

```

Output:

```
INSERT INTO `yu`.`table_name` VALUES (1257,34), (`ww`,`po`);
INSERT INTO `yu`.`table_name` VALUES (`wp`,`1`);
```

